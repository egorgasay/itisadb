package core

import (
	"context"
	"errors"
	"fmt"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/servers"
	"itisadb/pkg/logger"
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
	servers domains.Servers
	logger  logger.ILogger
	storage domains.Storage

	balancing bool
	objects   map[string]int32
	mu        sync.RWMutex

	pool   chan struct{} // TODO: ADD TO CONFIG
	keeper *Keeper
}

func New(ctx context.Context, repository domains.Storage, logger logger.ILogger) (*Core, error) {
	s, err := servers.New()
	if err != nil {
		return nil, err
	}

	objects, err := repository.RestoreObjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to restore objects: %w", err)
	}

	keeper, err := newKeeper(repository, logger, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create keeper: %w", err)
	}

	return &Core{
		servers:   s,
		storage:   repository,
		logger:    logger,
		objects:   objects,
		keeper:    keeper,
		pool:      make(chan struct{}, 10000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
		balancing: false,                                       // TODO: MOVE TO CONFIG
	}, nil
}

func toServerNumber(server *int32) int32 {
	if server == nil {
		return mainStorage
	}

	return *server
}

func (c *Core) Set(ctx context.Context, server *int32, key, val string, uniques bool) (int32, error) {
	if c.useMainStorage(server) {
		err := c.keeper.Set(key, val, uniques)
		if err != nil {
			return mainStorage, err
		}
		return mainStorage, nil
	}

	serverNumber := toServerNumber(server)

	if serverNumber == setToAll {
		failedServers := c.servers.SetToAll(ctx, key, val, uniques)
		if len(failedServers) != 0 {
			return setToAll, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}
		return setToAll, nil
	}
	cl, ok := c.servers.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return 0, constants.ErrUnknownServer
	}

	err := cl.Set(context.Background(), key, val, uniques)
	if err != nil {
		return 0, err
	}

	return cl.GetNumber(), nil
}

func (c *Core) Get(ctx context.Context, server *int32, key string) (val string, err error) {
	return val, c.withContext(ctx, func() error {
		val, err = c.get(ctx, server, key)
		return err
	})
}

func (c *Core) useMainStorage(server *int32) bool {
	return !c.balancing ||
		c.servers.Len() == 0 ||
		(server != nil && *server == mainStorage)
}

func (c *Core) get(ctx context.Context, server *int32, key string) (string, error) {
	if c.useMainStorage(server) {
		v, err := c.storage.Get(key)
		if err != nil {
			return v, err
		}
		return v, nil
	}

	serverNumber := toServerNumber(server)

	if serverNumber == searchEverywhere {
		v, err := c.servers.DeepSearch(ctx, key)
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

	res, err := cl.Get(context.Background(), key)
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

func (c *Core) Delete(ctx context.Context, server *int32, key string) (err error) {
	return c.withContext(ctx, func() error {
		return c.delete(ctx, server, key)
	})
}

func (c *Core) delete(ctx context.Context, server *int32, key string) error {
	if c.useMainStorage(server) {
		if err := c.storage.Delete(key); err != nil {
			return err
		}
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	serverNumber := toServerNumber(server)

	if serverNumber == deleteFromAll {
		atLeastOnce := c.servers.DelFromAll(ctx, key)
		if !atLeastOnce {
			return constants.ErrNotFound
		}
		return nil
	}

	cl, ok := c.servers.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return constants.ErrUnknownServer
	}

	err := cl.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) Authenticate(ctx context.Context, login string, password string) (string, error) {
	if password == "" {
		return "", constants.ErrWrongCredentials
	}

	return "token_for_" + login, nil
}
