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

type UseCase struct {
	servers *servers.Servers
	logger  *zap.Logger
	storage *repo.Storage
	//queue   *queue.Queue[int32]
}

func New(repository *repo.Storage, logger *zap.Logger) *UseCase {
	return &UseCase{
		servers: servers.New(),
		storage: repository,
		logger:  logger,
		//queue:   &queue.Queue[int32]{},
	}
}

func (uc *UseCase) Set(key string, val string) (int32, error) {
	if uc.servers.Len() == 0 {
		err := uc.storage.Set(key, val)
		if err != nil {
			uc.logger.Warn(err.Error())
			return 0, fmt.Errorf("error while setting new pair to dbstorage with no active grpc-storages: %w", err)
		}
		return 0, nil
	}

	cl, ok := uc.servers.GetClient()
	if !ok || cl == nil {
		err := uc.storage.Set(key, val)
		if err != nil {
			uc.logger.Warn(err.Error())
			return 0, fmt.Errorf("error while adding new pair to dbstorage with offline grpc-storage: %w", err)
		}
		return 0, nil
	}

	resp, err := cl.Set(context.Background(), key, val)
	if err != nil {
		return 0, nil
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

func (uc *UseCase) Get(key string, serverNumber int32) (string, error) {
	ctx := context.TODO()
	if uc.servers.Len() == 0 {
		return uc.FindInDB(key)
	}

	if serverNumber == -1 {
		value, err := uc.servers.DeepSearch(ctx, key)
		if errors.Is(err, servers.ErrNotFound) {
			return uc.FindInDB(key)
		}
		return value, err
	} else if serverNumber == 0 {
		return uc.FindInDB(key)
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
	if st.Code().String() == codes.NotFound.String() {
		return uc.FindInDB(key)
	}

	if st.Code().String() != codes.Unavailable.String() { // connection error
		return uc.FindInDB(key)
	}

	if cl.Tries > 2 {
		uc.Disconnect(cl.Number)
		cl.Tries = 0
	}

	return uc.FindInDB(key)
}

func (uc *UseCase) Connect(address string, available, total uint64) (int32, error) {
	uc.logger.Info("New request for connect from " + address)
	//numForReuse, _ := uc.queue.Dequeue()
	number, err := uc.servers.AddClient(address, available, total)
	if err != nil {
		uc.logger.Warn(err.Error())
		return 0, err
	}

	return number, nil
}

func (uc *UseCase) Disconnect(number int32) {
	uc.servers.Disconnect(number)
	//uc.queue.Enqueue(number)
}

func (uc *UseCase) Servers() []string {
	return uc.servers.GetServers()
}
