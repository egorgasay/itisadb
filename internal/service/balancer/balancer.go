package balancer

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

type Balancer struct {
	logger *zap.Logger

	servers domains.Servers
	storage domains.Storage
	tlogger domains.TransactionLogger
	session domains.Session

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
) (*Balancer, error) {
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

	return &Balancer{
		logger:        logger,
		servers:       servers,
		storage:       storage,
		tlogger:       tlogger,
		session:       session,
		cfg:           cfg,
		pool:          make(chan struct{}, 10_000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
		objectServers: gost.NewRwLock(make(map[string]int32)),
		keyServers:    gost.NewRwLock(make(map[string]int32)),
	}, nil
}

func toServerNumber(server *int32) int32 {
	if server == nil {
		return constants.MainStorageNumber
	}

	return *server
}

func (c *Balancer) Set(ctx context.Context, userID int, key, value string, opts models.SetOptions) (val int32, err error) {
	return val, c.withContext(ctx, func() error {
		val, err = c.set(ctx, userID, key, value, opts)
		return err
	})
}

func (c *Balancer) set(ctx context.Context, userID int, key, val string, opts models.SetOptions) (int32, error) {
	if opts.Server == _setToAll {
		failedServers := c.servers.SetToAll(ctx, key, val, opts)
		if len(failedServers) != 0 {
			return _setToAll, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}

		return _setToAll, nil
	}

	cl, ok := c.servers.GetServerByID(opts.Server)
	if !ok || cl == nil {
		return 0, constants.ErrUnknownServer
	}

	err := cl.SetOne(ctx, userID, key, val, opts).Error()
	if err != nil {
		return 0, err
	}

	return cl.Number(), nil
}

func (c *Balancer) Get(ctx context.Context, userID int, key string, opts models.GetOptions) (val models.Value, err error) {
	return val, c.withContext(ctx, func() error {
		val, err = c.get(ctx, userID, key, opts)
		return err
	})
}

func (c *Balancer) useMainStorage(server int32) bool {
	return !c.cfg.Balancer.On ||
		c.servers.Len() == 0 ||
		server == constants.MainStorageNumber
}

func (c *Balancer) getObjectInfo(object string) (models.ObjectInfo, error) {
	info, err := c.storage.GetObjectInfo(object)
	if err != nil {
		return models.ObjectInfo{}, fmt.Errorf("can't get object info: %w", err)
	}

	return info, nil
}

func (c *Balancer) get(ctx context.Context, userID int, key string, opts models.GetOptions) (models.Value, error) {
	if opts.Server == _searchEverywhere {
		v, err := c.servers.DeepSearch(ctx, key, opts)
		if err != nil && errors.Is(err, constants.ErrNotFound) {
			return models.Value{}, constants.ErrNotFound
		}
		return v, err
	} else if !c.servers.Exists(opts.Server) {
		return models.Value{}, constants.ErrNotFound
	}

	cl, ok := c.servers.GetServerByID(opts.Server)
	if !ok || cl == nil {
		return models.Value{}, constants.ErrUnknownServer
	}

	switch r := cl.GetOne(context.Background(), key, opts); r.Switch() {
	case gost.IsOk:
		cl.ResetTries()
		return r.Unwrap(), nil
	case gost.IsErr:
		err := r.Error()
		c.logger.Warn(err.Error())
		c.servers.OnServerError(cl, err)
	}

	return models.Value{}, constants.ErrNotFound
}

func (c *Balancer) Connect(ctx context.Context, address string, available, total uint64) (number int32, err error) {
	c.logger.Info("New request for connect from " + address)
	return number, c.withContext(ctx, func() error {
		number, err = c.servers.AddServer(address, available, total, c.servers.Len())
		if err != nil {
			c.logger.Warn(err.Error())
			return err
		}

		return nil
	})
}

func (c *Balancer) Disconnect(ctx context.Context, server int32) error {
	return c.withContext(ctx, func() error {
		c.servers.Disconnect(server)
		return nil
	})
}

func (c *Balancer) Servers() []string {
	return c.servers.GetServers()
}

func (c *Balancer) withContext(ctx context.Context, fn func() error) (err error) {
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

func (c *Balancer) Delete(ctx context.Context, userID int, key string, opts models.DeleteOptions) (err error) {
	return c.withContext(ctx, func() error {
		return c.delete(ctx, userID, key, opts)
	})
}

func (c *Balancer) delete(ctx context.Context, userID int, key string, opts models.DeleteOptions) error {
	if opts.Server == _deleteFromAll {
		atLeastOnce := c.servers.DelFromAll(ctx, key, opts)
		if !atLeastOnce {
			return constants.ErrNotFound
		}
		return nil
	}

	cl, ok := c.servers.GetServerByID(opts.Server)
	if !ok || cl == nil {
		return constants.ErrUnknownServer
	}

	err := cl.DelOne(ctx, key, opts).Error()
	if err != nil {
		return err
	}
	return nil
}

func (c *Balancer) CalculateRAM(_ context.Context) (res gost.Result[models.RAM]) {
	res = pkg.CalcRAM()
	if res.Error() != nil {
		c.logger.Error("Failed to calculate RAM", zap.Error(res.Error()))
	}

	return res
}

//func (c *Balancer) earlyObjectNotFound(name string, server int32) bool {
//	return c.earlyNotFound(c.cfg.Balancer.On, &c.objectServers, name, server)
//}
//
//func (c *Balancer) earlyKetNotFound(name string, server int32) bool {
//	return c.earlyNotFound(c.cfg.Balancer.On, &c.keyServers, name, server)
//}
//
//type RLocker interface {
//	RBorrow() gost.Arc[*map[string]int32, map[string]int32]
//	Release()
//}
//
//func (c *Balancer) earlyNotFound(isBalancer bool, locker RLocker, name string, server int32) bool {
//	if !isBalancer {
//		return false
//	}
//
//	defer locker.Release()
//
//	objServer, ok := (locker.RBorrow().Read())[name]
//	if !ok {
//		return false
//	}
//
//	return objServer != server
//}
