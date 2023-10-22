package core

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	"runtime"
	"sync"
)

const (
	searchEverywhere = iota * -1
	setToAll
)

const mainStorage = 1

const (
	deleteFromAll = -1
)

type Core struct {
	logger *zap.Logger

	servers domains.Servers
	storage domains.Storage
	tlogger domains.TransactionLogger
	session domains.Session

	objects map[string]int32
	mu      sync.RWMutex

	cfg *config.Config

	pool chan struct{} // TODO: ADD TO CONFIG
}

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.Logger,
	storage domains.Storage,
	tlogger domains.TransactionLogger,
	servers domains.Servers,
	session domains.Session,
) (*Core, error) {
	var err error

	objects := make(map[string]int32)
	if tlogger != nil {
		objects, err = tlogger.RestoreObjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to restore objects: %w", err)
		}
	}

	_, err = storage.CreateUser(
		models.User{
			Login:    "itisadb",
			Password: "itisadb",
			Level:    constants.SecretLevel, // TODO:
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
		objects: objects,
		mu:      sync.RWMutex{},
		cfg:     cfg,
		pool:    make(chan struct{}, 10_000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
	}, nil
}

func toServerNumber(server *int32) int32 {
	if server == nil {
		return mainStorage
	}

	return *server
}

func (c *Core) Set(ctx context.Context, userID int, key, val string, opts models.SetOptions) (int32, error) {
	if c.useMainStorage(opts.Server) {
		err := c.storage.Set(key, val, opts)
		if err != nil {
			return mainStorage, err
		}

		if c.cfg.TransactionLoggerConfig.On {
			c.tlogger.WriteSet(key, val)
		}

		return mainStorage, nil
	}

	serverNumber := toServerNumber(opts.Server)

	if serverNumber == setToAll {
		failedServers := c.servers.SetToAll(ctx, key, val, opts)
		if len(failedServers) != 0 {
			return setToAll, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}

		return setToAll, nil
	}
	cl, ok := c.servers.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return 0, constants.ErrUnknownServer
	}

	err := cl.Set(context.Background(), key, val, opts)
	if err != nil {
		return 0, err
	}

	return cl.GetNumber(), nil
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
		(server != nil && *server == mainStorage)
}

func (c *Core) get(ctx context.Context, userID int, key string, opts models.GetOptions) (string, error) {
	if c.useMainStorage(opts.Server) {
		v, err := c.storage.Get(key)
		if err != nil {
			return v, err
		}
		return v, nil
	}

	serverNumber := toServerNumber(opts.Server)

	if serverNumber == searchEverywhere {
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

	res, err := cl.Get(context.Background(), key, opts)
	if err == nil {
		cl.ResetTries()
		return res.Value, nil
	}

	c.logger.Warn(err.Error())

	cl.IncTries()

	if cl.GetTries() > 2 {
		err = c.Disconnect(ctx, cl.GetNumber())
		if err != nil {
			c.logger.Warn(err.Error())
		}
	}

	return "", constants.ErrNotFound
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
	if c.useMainStorage(opts.Server) {
		if err := c.storage.Delete(key); err != nil {
			c.logger.Warn("failed to delete", zap.Error(err))
			return err
		}

		if c.cfg.TransactionLoggerConfig.On {
			c.tlogger.WriteDelete(key)
		}

		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	serverNumber := toServerNumber(opts.Server)

	if serverNumber == deleteFromAll {
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

	err := cl.Delete(ctx, key, opts)
	if err != nil {
		return err
	}
	return nil
}
