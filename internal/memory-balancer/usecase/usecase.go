package usecase

import (
	"context"
	"errors"
	"fmt"
	//"github.com/tomakado/containers/queue"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-storage/internal/memory-balancer/servers"
	repo "grpc-storage/internal/memory-balancer/storage"
)

var ErrNoData = errors.New("the value is not found")
var ErrUnknownServer = errors.New("unknown server")

const (
	_ = iota * -1
	dbOnly
	setToAll
	AllAndDB
)

type UseCase struct {
	servers *servers.Servers
	logger  *zap.Logger
	storage *repo.Storage
	//queue   *queue.Queue[int32]
}

func New(repository *repo.Storage, logger *zap.Logger) (*UseCase, error) {
	s, err := servers.New()
	if err != nil {
		return nil, err
	}
	return &UseCase{
		servers: s,
		storage: repository,
		logger:  logger,
	}, nil
}

func (uc *UseCase) Set(ctx context.Context, key, val string, serverNumber int32, uniques bool) (int32, error) {
	setDB := uc.storage.Set
	if uniques {
		setDB = uc.storage.SetUnique
	}

	if uc.servers.Len() == 0 && serverNumber != -1 {
		err := setDB(key, val)
		if err != nil {
			uc.logger.Warn(err.Error())
			return 0, fmt.Errorf("error while setting new pair to dbstorage with no active grpc-storages: %w", err)
		}
		return -1, nil
	}

	switch serverNumber {
	case dbOnly:
		uc.logger.Info("setting k:val to db")
		return dbOnly, setDB(key, val)
	case setToAll:
		failedServers := uc.servers.SetToAll(ctx, key, val, uniques)
		if len(failedServers) != 0 {
			return setToAll, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}
		return setToAll, nil
	case AllAndDB:
		uc.logger.Info("setting key:val to all instance")
		failedServers := uc.servers.SetToAll(ctx, key, val, uniques)
		if len(failedServers) != 0 {
			return AllAndDB, fmt.Errorf("some servers wouldn't get values: %v", failedServers)
		}
		uc.logger.Info("setting key:val to db")
		return AllAndDB, setDB(key, val)
	}

	var cl *servers.Server
	var ok bool

	if serverNumber > 0 {
		cl, ok = uc.servers.GetClientByID(serverNumber)
		if !ok || cl == nil {
			return 0, ErrUnknownServer
		}
	} else {
		cl, ok = uc.servers.GetClient()
		if !ok || cl == nil {
			err := setDB(key, val)
			if err != nil {
				uc.logger.Warn(err.Error())
				return 0, fmt.Errorf("error while adding new pair to dbstorage with offline grpc-storage: %w", err)
			}
			return -1, nil
		}
	}

	resp, err := cl.Set(context.Background(), key, val, uniques)
	if err != nil {
		return 0, err
	}

	cl.Total = resp.Total
	cl.Available = resp.Available

	return cl.Number, nil
}

func (uc *UseCase) FindInDB(key string) (string, error) {
	value, err := uc.storage.Get(key)
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
		return uc.FindInDB(key)
	}

	if serverNumber == 0 {
		value, err := uc.servers.DeepSearch(ctx, key)
		if errors.Is(err, servers.ErrNotFound) {
			return uc.FindInDB(key)
		}
		return value, err
	} else if serverNumber == -1 {
		return uc.FindInDB(key)
	} else if !uc.servers.Exists(serverNumber) {
		return "", ErrUnknownServer
	}

	cl, ok := uc.servers.GetClientByID(serverNumber)
	if !ok || cl == nil {
		return uc.FindInDB(key)
	}

	res, err := cl.Get(context.Background(), key)
	if err == nil {
		cl.Total = res.Total
		cl.Available = res.Available
		return res.Value, nil
	}

	uc.logger.Warn(err.Error())
	st, ok := status.FromError(err)
	if !ok {
		return "", err
	}
	if st.Code().String() == codes.NotFound.String() || st.Code().String() != codes.Unavailable.String() {
		return uc.FindInDB(key)
	}

	if cl.Tries > 2 {
		uc.Disconnect(cl.Number)
	}
	cl.Tries = 0

	return uc.FindInDB(key)
}

func (uc *UseCase) Connect(address string, available, total uint64, server int32) (int32, error) {
	uc.logger.Info("New request for connect from " + address)
	number, err := uc.servers.AddClient(address, available, total, server)
	if err != nil {
		uc.logger.Warn(err.Error())
		return 0, err
	}

	return number, nil
}

func (uc *UseCase) Disconnect(number int32) {
	uc.servers.Disconnect(number)
}

func (uc *UseCase) Servers() []string {
	return uc.servers.GetServers()
}
