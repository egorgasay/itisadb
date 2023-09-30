package core

import (
	"context"
	"errors"
	"fmt"
	"itisadb/internal/domains"
	servers2 "itisadb/internal/servers"
	"itisadb/pkg/logger"
	"runtime"
	"sync"
)

var ErrNoServers = errors.New("no servers available")
var ErrNotFound = errors.New("key not found")
var ErrWrongCredentials = errors.New("wrong credentials")

var ErrNoData = errors.New("the value is not found")
var ErrUnknownServer = errors.New("unknown server")

const (
	searchEverywhere = iota * -1
	setToAll
)

const (
	deleteFromAll = -1
)

type Core struct {
	servers domains.Servers
	logger  logger.ILogger
	storage domains.Storage

	objects map[string]int32
	mu      sync.RWMutex

	pool   chan struct{} // TODO: ADD TO CONFIG
	keeper *Keeper
}

func New(ctx context.Context, repository domains.Storage, logger logger.ILogger) (*Core, error) {
	s, err := servers2.New()
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
		servers: s,
		storage: repository,
		logger:  logger,
		objects: objects,
		keeper:  keeper,
		pool:    make(chan struct{}, 10000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
	}, nil
}

func (uc *Core) Set(ctx context.Context, key, val string, serverNumber int32, uniques bool) (int32, error) {
	if uc.servers.Len() == 0 {
		return 0, ErrNoServers
	}

	if serverNumber == setToAll {
		failedServers := uc.servers.SetToAll(ctx, key, val, uniques)
		if len(failedServers) != 0 {
			return setToAll, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}
		return setToAll, nil
	}

	var cl *servers2.Server
	var ok bool

	if serverNumber > 0 {
		cl, ok = uc.servers.GetServerByID(serverNumber)
		if !ok || cl == nil {
			return 0, ErrUnknownServer
		}
	} else {
		cl, ok = uc.servers.GetServer()
		if !ok || cl == nil {
			return 0, ErrNoServers
		}
	}

	err := cl.Set(context.Background(), key, val, uniques)
	if err != nil {
		return 0, err
	}

	return cl.GetNumber(), nil
}

func (uc *Core) Get(ctx context.Context, key string, serverNumber int32) (val string, err error) {
	return val, uc.withContext(ctx, func() error {
		val, err = uc.get(ctx, key, serverNumber)
		return err
	})
}

func (uc *Core) get(ctx context.Context, key string, serverNumber int32) (string, error) {
	if uc.servers.Len() == 0 {
		return "", ErrNoServers
	}

	if serverNumber == searchEverywhere {
		v, err := uc.servers.DeepSearch(ctx, key)
		if err != nil && errors.Is(err, servers2.ErrNotFound) {
			return "", ErrNotFound
		}
		return v, err
	} else if !uc.servers.Exists(serverNumber) {
		return "", ErrNotFound
	}

	cl, ok := uc.servers.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return "", ErrUnknownServer
	}

	res, err := cl.Get(context.Background(), key)
	if err == nil {
		cl.ResetTries()
		return res.Value, nil
	}

	uc.logger.Warn(err.Error())

	cl.IncTries()

	if cl.GetTries() > 2 {
		err = uc.Disconnect(ctx, cl.GetNumber())
		if err != nil {
			uc.logger.Warn(err.Error())
		}
	}

	return "", ErrNotFound
}

func (uc *Core) Connect(address string, available, total uint64, server int32) (int32, error) {
	uc.logger.Info("New request for connect from " + address)
	number, err := uc.servers.AddServer(address, available, total, server)
	if err != nil {
		uc.logger.Warn(err.Error())
		return 0, err
	}

	return number, nil
}

func (uc *Core) Disconnect(ctx context.Context, number int32) error {
	return uc.withContext(ctx, func() error {
		uc.servers.Disconnect(number)
		return nil
	})
}

func (uc *Core) Servers() []string {
	return uc.servers.GetServers()
}

func (uc *Core) withContext(ctx context.Context, fn func() error) (err error) {
	ch := make(chan struct{})

	once := sync.Once{}
	done := func() { close(ch) }

	uc.pool <- struct{}{}
	go func() {
		err = fn()
		once.Do(done)
		<-uc.pool
	}()

	select {
	case <-ch:
		return err
	case <-ctx.Done():
		once.Do(done)
		return ctx.Err()
	}
}

func (uc *Core) Delete(ctx context.Context, key string, num int32) (err error) {
	return uc.withContext(ctx, func() error {
		return uc.delete(ctx, key, num)
	})
}

func (uc *Core) delete(ctx context.Context, key string, num int32) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	if num == deleteFromAll {
		atLeastOnce := uc.servers.DelFromAll(ctx, key)
		if !atLeastOnce {
			return ErrNotFound
		}
		return nil
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrUnknownServer
	}

	err := cl.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Core) Authenticate(ctx context.Context, login string, password string) (string, error) {
	if password == "" {
		return "", ErrWrongCredentials
	}

	return "token_for_" + login, nil
}
