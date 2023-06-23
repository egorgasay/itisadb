package usecase

import (
	"context"
	"errors"
	"fmt"
	"itisadb/internal/memory-balancer/servers"
	"sync"
)

var ErrNoServers = errors.New("no servers available")
var ErrNotFound = errors.New("key not found")

func (uc *UseCase) Set(ctx context.Context, key, val string, serverNumber int32, uniques bool) (int32, error) {
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

	var cl *servers.Server
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

func (uc *UseCase) Get(ctx context.Context, key string, serverNumber int32) (val string, err error) {
	return val, uc.withContext(ctx, func() error {
		val, err = uc.get(ctx, key, serverNumber)
		return err
	})
}

func (uc *UseCase) get(ctx context.Context, key string, serverNumber int32) (string, error) {
	if uc.servers.Len() == 0 {
		return "", ErrNoServers
	}

	if serverNumber == searchEverywhere {
		v, err := uc.servers.DeepSearch(ctx, key)
		if err != nil && errors.Is(err, servers.ErrNotFound) {
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

func (uc *UseCase) Connect(address string, available, total uint64, server int32) (int32, error) {
	uc.logger.Info("New request for connect from " + address)
	number, err := uc.servers.AddServer(address, available, total, server)
	if err != nil {
		uc.logger.Warn(err.Error())
		return 0, err
	}

	return number, nil
}

func (uc *UseCase) Disconnect(ctx context.Context, number int32) error {
	return uc.withContext(ctx, func() error {
		uc.servers.Disconnect(number)
		return nil
	})
}

func (uc *UseCase) Servers() []string {
	return uc.servers.GetServers()
}

func (uc *UseCase) withContext(ctx context.Context, fn func() error) (err error) {
	ch := make(chan struct{})
	done := func() { close(ch) }

	once := sync.Once{}
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

func (uc *UseCase) Delete(ctx context.Context, key string, num int32) (err error) {
	return uc.withContext(ctx, func() error {
		return uc.delete(ctx, key, num)
	})
}

func (uc *UseCase) delete(ctx context.Context, key string, num int32) error {
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
