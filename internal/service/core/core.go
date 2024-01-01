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

	servers domains.Servers
	storage domains.Storage
	tlogger domains.TransactionLogger
	session domains.Session

	cfg config.Config

	pool chan struct{} // TODO: ADD TO CONFIG
}

func New(
	ctx context.Context,
	cfg config.Config,
	logger *zap.Logger,
	storage domains.Storage,
	tlogger domains.TransactionLogger,
	servers domains.Servers,
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
		logger:  logger,
		servers: servers,
		storage: storage,
		tlogger: tlogger,
		session: session,
		cfg:     cfg,
		pool:    make(chan struct{}, 10_000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
	}, nil
}

func toServerNumber(server *int32) int32 {
	if server == nil {
		return _mainStorage
	}

	return *server
}

func (c *Core) Set(ctx context.Context, userID int, key, val string, opts models.SetOptions) (int32, error) {
	if !c.useMainStorage(opts.Server) {
		serverNumber := toServerNumber(opts.Server)

		if serverNumber == _setToAll {
			failedServers := c.servers.SetToAll(ctx, key, val, opts)
			if len(failedServers) != 0 {
				return _setToAll, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
			}

			return _setToAll, nil
		}
		cl, ok := c.servers.GetServerByID(serverNumber)
		if !ok || cl == nil {
			return 0, constants.ErrUnknownServer
		}

		err := cl.SetOne(context.Background(), key, val, opts.ToSDK()).Error()
		if err != nil {
			return 0, err
		}

		return cl.Number(), nil
	}

	v, err := c.storage.Get(key)
	if err == nil && (opts.Unique || v.ReadOnly) {
		return _mainStorage, constants.ErrAlreadyExists
	}

	err = c.storage.Set(key, val, opts)
	if err != nil {
		return _mainStorage, err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteSet(key, val, opts)
	}

	return _mainStorage, nil
}

func (c *Core) Get(ctx context.Context, userID int, key string, opts models.GetOptions) (val string, err error) {
	return val, c.withContext(ctx, func() error {
		val, err = c.get(ctx, userID, key, opts)
		return err
	})
}

func (c *Core) useMainStorage(server *int32) bool {
	return !c.cfg.Balancer.On ||
		c.servers.Len() == 0 ||
		(server != nil && *server == _mainStorage)
}

func (c *Core) getObjectInfo(object string) (models.ObjectInfo, error) {
	info, err := c.storage.GetObjectInfo(object)
	if err != nil {
		return models.ObjectInfo{}, fmt.Errorf("can't get object info: %w", err)
	}

	return info, nil
}

func (c *Core) get(ctx context.Context, userID int, key string, opts models.GetOptions) (string, error) {
	if !c.useMainStorage(opts.Server) {
		serverNumber := toServerNumber(opts.Server)

		if serverNumber == _searchEverywhere {
			v, err := c.servers.DeepSearch(ctx, key, opts)
			if err != nil && errors.Is(err, constants.ErrNotFound) {
				return "", constants.ErrNotFound
			}
			return v, err
		} else if !c.servers.Exists(serverNumber) {
			return "", constants.ErrNotFound
		}

		cl, ok := c.servers.GetServerByID(serverNumber)
		if !ok || cl == nil {
			return "", constants.ErrUnknownServer
		}

		var res string
		switch r := cl.GetOne(context.Background(), key, opts.ToSDK()); r.Switch() {
		case gost.IsOk:
			cl.ResetTries()
			res = r.Unwrap()

			return res, nil
		case gost.IsErr:
			c.logger.Warn(r.Error().Error())
		}

		cl.IncTries()

		if cl.Tries() > 2 {
			err := c.Disconnect(ctx, cl.Number())
			if err != nil {
				c.logger.Warn(err.Error())
			}
		}

		return "", constants.ErrNotFound
	}

	v, err := c.storage.Get(key)
	if err != nil {
		return "", err
	}

	return v.Value, nil
}

func (c *Core) Connect(address string, available, total uint64) (int32, error) {
	c.logger.Info("New request for connect from " + address)
	number, err := c.servers.AddServer(address, available, total, c.servers.Len())
	if err != nil {
		c.logger.Warn(err.Error())
		return 0, err
	}

	return number, nil
}

func (c *Core) Disconnect(ctx context.Context, server int32) error {
	return c.withContext(ctx, func() error {
		c.servers.Disconnect(server)
		return nil
	})
}

func (c *Core) Servers() []string {
	return c.servers.GetServers()
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
	if !c.useMainStorage(opts.Server) {
		serverNumber := toServerNumber(opts.Server)

		if serverNumber == _deleteFromAll {
			atLeastOnce := c.servers.DelFromAll(ctx, key, opts)
			if !atLeastOnce {
				return constants.ErrNotFound
			}
			return nil
		}

		cl, ok := c.servers.GetServerByID(serverNumber)
		if !ok || cl == nil {
			return constants.ErrUnknownServer
		}

		err := cl.DelOne(ctx, key, opts.ToSDK()).Error()
		if err != nil {
			return err
		}
		return nil
	}

	if err := c.storage.Delete(key); err != nil {
		c.logger.Warn("failed to delete", zap.Error(err))
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteDelete(key)
	}

	return nil
}
