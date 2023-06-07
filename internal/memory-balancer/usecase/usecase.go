package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	//"github.com/tomakado/containers/queue"
	"go.mongodb.org/mongo-driver/mongo"
	"itisadb/internal/memory-balancer/servers"
)

func (uc *UseCase) Set(ctx context.Context, key, val string, serverNumber int32, uniques bool) (int32, error) {
	setDB := uc.storage.Set
	if uniques {
		setDB = uc.storage.SetUnique
	}

	if uc.servers.Len() == 0 && serverNumber != -1 {
		err := setDB(ctx, key, val)
		if err != nil {
			uc.logger.Warn(err.Error())
			return 0, fmt.Errorf("error while setting new pair to dbstorage with no active grpc-storages: %w", err)
		}
		return -1, nil
	}

	switch serverNumber {
	case dbOnly:
		uc.logger.Info("setting k:val to db")
		return dbOnly, setDB(ctx, key, val)
	case all:
		failedServers := uc.servers.SetToAll(ctx, key, val, uniques)
		if len(failedServers) != 0 {
			return all, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}
		return all, nil
	case allAndDB:
		uc.logger.Info("setting key:val to all instance")
		failedServers := uc.servers.SetToAll(ctx, key, val, uniques)
		if len(failedServers) != 0 {
			return allAndDB, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}
		uc.logger.Info("setting key:val to db")
		return allAndDB, setDB(ctx, key, val)
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
			err := setDB(ctx, key, val)
			if err != nil {
				uc.logger.Warn(err.Error())
				return 0, fmt.Errorf("error while adding new pair to dbstorage with offline grpc-storage: %w", err)
			}
			return -1, nil
		}
	}

	err := cl.Set(context.Background(), key, val, uniques)
	if err != nil {
		return 0, err
	}

	return cl.GetNumber(), nil
}

var timeout = 4 * time.Second

func (uc *UseCase) FindInDB(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	value, err := uc.storage.Get(ctx, key)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("error while getting new pair from dbstorage: %w", ErrNoData)
		}
		uc.logger.Warn(err.Error())
		return value, fmt.Errorf("error while getting new pair from dbstorage: %w", err)
	}
	return value, nil
}

func (uc *UseCase) Get(ctx context.Context, key string, serverNumber int32) (string, error) {
	if uc.servers.Len() == 0 {
		return uc.FindInDB(ctx, key)
	}

	if serverNumber == 0 {
		value, err := uc.servers.DeepSearch(ctx, key)
		if errors.Is(err, servers.ErrNotFound) {
			return uc.FindInDB(ctx, key)
		}
		return value, err
	} else if serverNumber == -1 {
		return uc.FindInDB(ctx, key)
	} else if !uc.servers.Exists(serverNumber) {
		return uc.FindInDB(ctx, key)
	}

	cl, ok := uc.servers.GetServerByID(serverNumber)
	if !ok || cl == nil {
		return uc.FindInDB(ctx, key)
	}

	res, err := cl.Get(context.Background(), key)
	if err == nil {
		return res.Value, nil
	}

	uc.logger.Warn(err.Error())

	if cl.GetTries() > 2 {
		err = uc.Disconnect(ctx, cl.GetNumber())
		if err != nil {
			uc.logger.Warn(err.Error())
		}
	}
	cl.ResetTries()

	return uc.FindInDB(ctx, key)
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
	ch := make(chan struct{})
	uc.pool <- struct{}{}
	go func() { // TODO: add pool
		uc.servers.Disconnect(number)
		close(ch)
		<-uc.pool
	}()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (uc *UseCase) Servers() []string {
	return uc.servers.GetServers()
}

func (uc *UseCase) Delete(ctx context.Context, key string, num int32) (err error) {
	ch := make(chan struct{})

	uc.pool <- struct{}{}
	go func() {
		err = uc.delete(ctx, key, num)
		close(ch)
		<-uc.pool
	}()

	select {
	case <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (uc *UseCase) delete(ctx context.Context, key string, num int32) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if num == 0 {
		num = all
	}

	switch num {
	case dbOnly:
		//  TODO: delete from db
	case all:
		// TODO: delete from all servers
	case allAndDB:
		// TODO: delete from all servers
		// TODO: delete from db
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrUnknownServer
	}

	err := cl.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("error while deleting value: %w", err)
	}
	return nil
}
