package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	"itisadb/pkg"
	"runtime"
	"sync"
)

const (
	_searchEverywhere = iota * -1
	_setToAll
)

const (
	_autoSelect  = 0
	_mainStorage = 1
)

const (
	_deleteFromAll = -1
)

type Core struct {
	logger *zap.Logger

	balancer domains.Balancer
	storage  domains.Storage
	tlogger  domains.TransactionLogger
	session  domains.Session

	cfg config.Config

	pool chan struct{} // TODO: ADD TO CONFIG

	objectsInfo gost.RwLock[map[string]models.ObjectInfo]
}

func New(
	ctx context.Context,
	cfg config.Config,
	logger *zap.Logger,
	storage domains.Storage,
	tlogger domains.TransactionLogger,
	servers domains.Balancer,
	session domains.Session,
) (*Core, error) {
	var err error

	_, err = storage.CreateUser(
		models.User{
			Login:    "itisadb",
			Password: "itisadb",
			Level:    constants.SecretLevel,
			Active:   true,
		},
	)

	if err != nil && !errors.Is(err, constants.ErrAlreadyExists) {
		return nil, err
	}

	return &Core{
		logger:   logger,
		balancer: servers,
		storage:  storage,
		tlogger:  tlogger,
		session:  session,
		cfg:      cfg,
		pool:     make(chan struct{}, 10_000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
	}, nil
}

func toServerNumber(server *int32) int32 {
	if server == nil {
		return constants.MainStorageNumber
	}

	return *server
}

func (c *Core) Set(ctx context.Context, userID int, key, value string, opts models.SetOptions) (val int32, err error) {
	if !c.hasPermission(userID, models.Level(gost.SafeDeref(opts.Level))) {
		return 0, constants.ErrForbidden
	}

	return val, c.withContext(ctx, func() error {
		val, err = c.set(ctx, key, value, opts)
		return err
	})
}

func (c *Core) set(ctx context.Context, key, val string, opts models.SetOptions) (int32, error) {
	serverNumber := toServerNumber(opts.Server)

	if serverNumber == _setToAll {
		failedServers := c.balancer.SetToAll(ctx, key, val, opts)
		if len(failedServers) != 0 {
			return _setToAll, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}

		return _setToAll, nil
	}

	cl, ok := c.balancer.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return 0, constants.ErrUnknownServer
	}

	err := cl.SetOne(ctx, key, val, opts).Error()
	if err != nil {
		return 0, err
	}

	return cl.Number(), nil
}

func (c *Core) Get(ctx context.Context, userID int, key string, opts models.GetOptions) (val models.Value, err error) {
	err := c.withContext(ctx, func() error {
		val, err = c.get(ctx, userID, key, opts)
		return err
	})

	if err != nil {
		return models.Value{}, err
	}

	return val, nil
}

func (c *Core) useMainStorage(server *int32) bool {
	return !c.cfg.Balancer.On ||
		c.balancer.Len() == 0 ||
		(server != nil && *server == _mainStorage)
}

func (c *Core) getObjectInfo(object string) (models.ObjectInfo, error) {
	info, err := c.storage.GetObjectInfo(object)
	if err != nil {
		return models.ObjectInfo{}, fmt.Errorf("can't get object info: %w", err)
	}

	return info, nil
}

func (c *Core) get(ctx context.Context, userID int, key string, opts models.GetOptions) (models.Value, error) {
	serverNumber := toServerNumber(opts.Server)

	if serverNumber == _searchEverywhere {
		v, err := c.balancer.DeepSearch(ctx, key, opts)
		if err != nil && errors.Is(err, constants.ErrNotFound) {
			return models.Value{}, constants.ErrNotFound
		}
		return v, err
	} else if !c.balancer.Exists(serverNumber) {
		return "", constants.ErrNotFound
	}

	cl, ok := c.balancer.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return "", constants.ErrUnknownServer
	}

	switch r := cl.GetOne(context.Background(), key, opts); r.Switch() {
	case gost.IsOk:
		cl.ResetTries()
		return r.Unwrap(), nil
	case gost.IsErr:
		c.logger.Warn(r.Error().Error())
		c.balancer.OnServerError(serverNumber)
	}

	return "", constants.ErrNotFound
}

func (c *Core) Connect(ctx context.Context, address string, available, total uint64) (number int32, err error) {
	c.logger.Info("New request for connect from " + address)
	return number, c.withContext(ctx, func() error {
		number, err = c.balancer.AddServer(address, available, total, c.balancer.Len())
		if err != nil {
			c.logger.Warn(err.Error())
			return err
		}

		return nil
	})
}

func (c *Core) Disconnect(ctx context.Context, server int32) error {
	return c.withContext(ctx, func() error {
		c.balancer.Disconnect(server)
		return nil
	})
}

func (c *Core) Servers() []string {
	return c.balancer.GetServers()
}

func (c *Core) withContext(ctx context.Context, fn func() error) (err error) {
	ch := make(chan struct{})

	once := sync.Once{}
	done := func() { close(ch) }

	c.pool <- struct{}{}
	go func() {
		err = fn()
		once.Do(done)
		<-c.pool
	}()

	select {
	case <-ch:
		return err
	case <-ctx.Done():
		once.Do(done)
		return ctx.Err()
	}
}

func (c *Core) Delete(ctx context.Context, userID int, key string, opts models.DeleteOptions) (err error) {
	return c.withContext(ctx, func() error {
		return c.delete(ctx, userID, key, opts)
	})
}

func (c *Core) delete(ctx context.Context, userID int, key string, opts models.DeleteOptions) error {
	serverNumber := toServerNumber(opts.Server)

	if serverNumber == _deleteFromAll {
		atLeastOnce := c.balancer.DelFromAll(ctx, key, opts)
		if !atLeastOnce {
			return constants.ErrNotFound
		}
		return nil
	}

	cl, ok := c.balancer.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return constants.ErrUnknownServer
	}

	err := cl.DelOne(ctx, key, opts).Error()
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) CalculateRAM(_ context.Context) (res gost.Result[models.RAM]) {
	res = pkg.CalcRAM()
	if res.Error() != nil {
		c.logger.Error("Failed to calculate RAM", zap.Error(res.Error()))
	}

	return res
}
