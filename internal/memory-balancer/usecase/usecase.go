package usecase

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"sync"

	repo "grpc-storage/internal/memory-balancer/storage"
	"grpc-storage/pkg/api/storage"

	"github.com/tomakado/containers/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var ErrNoData = errors.New("the value is not found")

type UseCase struct {
	clients map[uint64]*client
	sync.RWMutex
	logger  *zap.Logger
	storage *repo.Storage
	queue   *queue.Queue[uint64]
}

type client struct {
	tries     uint
	storage   storage.StorageClient
	available uint64
	total     uint64
}

type RAM struct {
}

func New(repository *repo.Storage, logger *zap.Logger) *UseCase {
	Queue := &queue.Queue[uint64]{}

	clients := make(map[uint64]*client, 10)
	return &UseCase{
		clients: clients,
		storage: repository,
		logger:  logger,
		queue:   Queue,
	}
}

func (uc *UseCase) Set(key string, val string) (uint64, error) {
	uc.RLock()
	defer uc.RUnlock()
	if len(uc.clients) == 0 {
		err := uc.storage.Set(key, val)
		if err != nil {
			uc.logger.Warn(err.Error())
			return 0, fmt.Errorf("error while setting new pair to dbstorage with no active grpc-storages: %w", err)
		}
		return 0, nil
	}
	serverNumber := uint64(len(key)%len(uc.clients) + 1)
	cl, ok := uc.clients[serverNumber]
	if !ok || cl == nil {
		err := uc.storage.Set(key, val)
		if err != nil {
			uc.logger.Warn(err.Error())
			return 0, fmt.Errorf("error while adding new pair to dbstorage with offline grpc-storage: %w", err)
		}
		return 0, nil
	}

	resp, err := cl.storage.Set(context.Background(), &storage.SetRequest{Key: key, Value: val})
	if err != nil {
		return 0, nil
	}

	cl.total = resp.Total
	cl.available = resp.Available

	return serverNumber, nil
}

func (uc *UseCase) Get(key string) (string, error) {
	uc.RLock()
	defer uc.RUnlock()

	if len(uc.clients) == 0 {
		value, err := uc.storage.Get(key)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return "", fmt.Errorf("error while getting new pair to dbstorage with no active grpc-storages: %w", ErrNoData)
			}
			uc.logger.Warn(err.Error())
			return value, fmt.Errorf("error while getting new pair to dbstorage with no active grpc-storages: %w", err)
		}
		return value, nil
	}

	serverNumber := uint64(len(key)%(len(uc.clients)) + 1)
	cl, ok := uc.clients[serverNumber]
	if !ok || cl == nil {
		value, err := uc.storage.Get(key)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return "", fmt.Errorf("error while getting new pair to dbstorage with offline grpc-storages: %w", ErrNoData)
			}
			uc.logger.Warn(err.Error())
			return value, fmt.Errorf("error while getting new pair to dbstorage with offline grpc-storage: %w", err)
		}
		return value, nil
	}

	res, err := cl.storage.Get(context.Background(), &storage.GetRequest{Key: key})
	if err == nil {
		cl.total = res.Total
		cl.available = res.Available
		return res.Value, nil
	}

	uc.logger.Warn(err.Error())
	st, ok := status.FromError(err)
	if !ok {
		return "", err
	}
	if st.Code().String() == codes.NotFound.String() {
		return "", ErrNoData
	}

	if st.Code().String() != codes.Unavailable.String() { // connection error
		return "", fmt.Errorf("can't get the value from server: %w", err)
	}

	if cl.tries > 2 {
		uc.Disconnect(serverNumber)
		cl.tries = 0
	}

	get, err := uc.storage.Get(key)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", ErrNoData
	}
	return get, err
}

func (uc *UseCase) Connect(address string) (uint64, error) {
	uc.Lock()
	defer uc.Unlock()

	uc.logger.Info("New request for connect from " + address)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return 0, err
	}

	cl := storage.NewStorageClient(conn)

	var stClient = &client{
		tries:   0,
		storage: cl,
	}

	numForReuse, ok := uc.queue.Dequeue()
	number := uint64(len(uc.clients) + 1)
	if ok {
		number = numForReuse
	}

	uc.clients[number] = stClient
	return number, nil
}

func (uc *UseCase) Disconnect(number uint64) {
	uc.RLock()
	defer uc.RUnlock()
	uc.clients[number] = nil
	uc.queue.Enqueue(number)
}

func (uc *UseCase) Servers() []string {
	uc.RLock()
	defer uc.RUnlock()
	var servers = make([]string, 0, 5)
	for num, cl := range uc.clients {
		if cl != nil {
			servers = append(servers, fmt.Sprintf("s#%d Avaliable: %d MB, Total: %d MB", num, cl.available, cl.total))
		} else {
			servers = append(servers, fmt.Sprintf("s#%d empty server slot", num))
		}
	}

	return servers
}
