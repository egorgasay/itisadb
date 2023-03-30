package usecase

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/egorgasay/grpc-storage/pkg/api/storage"
)

type UseCase struct {
	clients map[int]storage.StorageClient
}

var ErrServerIsNotConnected = errors.New("the server is not connected")

func New() *UseCase {
	clients := make(map[int]storage.StorageClient, 10)
	return &UseCase{clients: clients}
}

func (uc *UseCase) Set(key string, val string) (uint64, error) {
	if len(uc.clients) == 0 {
		return 0, errors.New("no available storages")
	}
	serverNumber := len(key)%len(uc.clients) + 1
	cl, ok := uc.clients[serverNumber]
	if !ok {
		return 0, ErrServerIsNotConnected
	}

	_, err := cl.Set(context.Background(), &storage.SetRequest{Key: key, Value: val})
	if err != nil {
		return 0, nil
	}

	return uint64(serverNumber), nil
}

func (uc *UseCase) Get(key string) (string, error) {
	if len(uc.clients) == 0 {
		return "", errors.New("no available storages")
	}
	cl, ok := uc.clients[len(key)%(len(uc.clients))+1]
	if !ok {
		return "", ErrServerIsNotConnected
	}

	res, err := cl.Get(context.Background(), &storage.GetRequest{Key: key})
	if err != nil {
		return "", nil
	}

	return res.Value, nil
}

func (uc *UseCase) Connect(address string) error {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	cl := storage.NewStorageClient(conn)
	uc.clients[len(uc.clients)+1] = cl
	return nil
}
