package balancer

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	"itisadb/pkg"
)

type Balancer struct {
	logger *zap.Logger

	servers  domains.Servers
	storage  domains.Storage
	tlogger  domains.TransactionLogger
	session  domains.Session
	security domains.SecurityService

	cfg config.Config

	pool chan struct{} // TODO: ADD TO CONFIG

	objectServers gost.RwLock[map[string]int32]
	keyServers    gost.RwLock[map[string]int32]
}

func New(
	ctx context.Context,
	cfg config.Config,
	logger *zap.Logger,
	storage domains.Storage,
	tlogger domains.TransactionLogger,
	servers domains.Servers,
	session domains.Session,
	security domains.SecurityService,
) (*Balancer, error) {
	var err error

	if err != nil && !errors.Is(err, constants.ErrAlreadyExists) {
		return nil, err
	}

	for _, server := range cfg.Balancer.Servers {
		logger.Info("Adding server", zap.String("server", server))

		func() {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, constants.ServerConnectTimeout)
			defer cancel()

			s, err := servers.AddServer(ctxWithTimeout, server, true)
			if err != nil {
				logger.Error("Failed to add server", zap.String("server", server), zap.Error(err))
			} else {
				logger.Info("Added server", zap.Int32("server", s))
			}
		}()
	}

	return &Balancer{
		logger:        logger,
		servers:       servers,
		storage:       storage,
		tlogger:       tlogger,
		session:       session,
		cfg:           cfg,
		pool:          make(chan struct{}, 20_000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
		objectServers: gost.NewRwLock(make(map[string]int32)),
		keyServers:    gost.NewRwLock(make(map[string]int32)),
		security:      security,
	}, nil
}

func (c *Balancer) Set(ctx context.Context, claims gost.Option[models.UserClaims], key, value string, opts models.SetOptions) (val int32, err error) {
	return val, gost.WithContextPool(ctx, func() error {
		val, err = c.set(ctx, claims, key, value, opts)
		return err
	}, c.pool)
}

func (c *Balancer) set(ctx context.Context, claims gost.Option[models.UserClaims], key, val string, opts models.SetOptions) (int32, error) {
	if opts.Server == constants.AutoServerNumber {
		res := c.getKeyServer(key)
		if res.IsSome() {
			opts.Server = res.Unwrap()
		}
	} else if opts.Server == constants.SetToAllServers {
		failedServers := c.servers.SetToAll(ctx, claims, key, val, opts)
		if len(failedServers) != 0 {
			return opts.Server, fmt.Errorf("some servers failed: %v", failedServers)
		}

		return opts.Server, nil
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return 0, constants.ErrUnknownServer
	}

	err := cl.SetOne(ctx, claims, key, val, opts).Error()
	if err != nil {
		return 0, err
	}

	c.addKeyServer(key, cl.Number())

	return cl.Number(), nil
}

func (c *Balancer) Get(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.GetOptions) (val models.Value, err error) {
	return val, gost.WithContextPool(ctx, func() error {
		val, err = c.get(ctx, claims, key, opts)
		return err
	}, c.pool)
}

func (c *Balancer) getObjectInfo(object string) (models.ObjectInfo, error) {
	info, err := c.storage.GetObjectInfo(object)
	if err != nil {
		return models.ObjectInfo{}, fmt.Errorf("can't get object info: %w", err)
	}

	return info, nil
}

func (c *Balancer) get(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.GetOptions) (models.Value, error) {
	if opts.Server == constants.AutoServerNumber {
		res := c.getKeyServer(key)
		if res.IsNone() {
			if r := c.servers.DeepSearch(ctx, claims, key, opts); r.IsErr() {
				return models.Value{}, fmt.Errorf("can't get key after deep search: %w", r.Error())
			} else {
				res := r.Unwrap()
				c.addKeyServer(key, res.Left)

				return res.Right, nil
			}
		}

		opts.Server = res.Unwrap()
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return models.Value{}, constants.ErrUnknownServer
	}

	switch r := cl.GetOne(ctx, claims, key, opts); r.IsOk() {
	case true:
		return r.Unwrap(), nil
	default:
		return models.Value{}, r.Error().WrapfNotNilMsg("can't get key from server: %s", cl.Number())
	}
}

func (c *Balancer) Connect(ctx context.Context, address string, available, total uint64) (number int32, err error) {
	c.logger.Info("New request for connect from " + address)
	return number, gost.WithContextPool(ctx, func() error {
		number, err = c.servers.AddServer(ctx, address, false)
		if err != nil {
			c.logger.Warn(err.Error())
			return err
		}

		return nil
	}, c.pool)
}

func (c *Balancer) Disconnect(ctx context.Context, server int32) error {
	return gost.WithContextPool(ctx, func() error {
		c.servers.Disconnect(server)
		return nil
	}, c.pool)
}

func (c *Balancer) Servers() []string {
	return c.servers.GetServersInfo()
}

func (c *Balancer) Delete(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.DeleteOptions) (err error) {
	return gost.WithContextPool(ctx, func() error {
		return c.delete(ctx, claims, key, opts)
	}, c.pool)
}

func (c *Balancer) delete(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.DeleteOptions) (err error) {
	if opts.Server == constants.DeleteFromAllServers {
		atLeastOnce := c.servers.DelFromAll(ctx, claims, key, opts)
		if !atLeastOnce {
			return constants.ErrNotFound
		}
		return nil
	} else if opts.Server == constants.AutoServerNumber {
		switch res := c.getKeyServer(key); res.IsSome() {
		case true:
			opts.Server = res.Unwrap()
			defer func() {
				if err != nil {
					c.delKeyServer(key)
				}
			}()
		default:
			if r := c.servers.DeepSearch(ctx, claims, key, models.GetOptions{}); r.IsErr() {
				return fmt.Errorf("can't delete key after deep search: %w", r.Error())
			} else {
				opts.Server = r.Unwrap().Left
			}
		}
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return constants.ErrUnknownServer
	}

	return cl.DelOne(ctx, claims, key, opts).Error().WrapNotNilMsg("can't delete key").IntoStd()
}

func (c *Balancer) CalculateRAM(_ context.Context) (res gost.Result[models.RAM]) {
	res = pkg.CalcRAM()
	if res.Error() != nil {
		c.logger.Error("Failed to calculate RAM", zap.Error(res.Error()))
	}

	return res
}
