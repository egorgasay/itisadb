package usecase

import (
	"context"
	"errors"
	"itisadb/internal/memory-balancer/servers"
	"itisadb/pkg/logger"
	"sync"
)

var ErrNoData = errors.New("the value is not found")
var ErrUnknownServer = errors.New("unknown server")

const (
	searchEverywhere = iota * -1
	setToAll
)

const (
	deleteFromAll = -1
)

//go:generate mockgen -destination=mocks/storage/mock_storage.go -package=mocks . IStorage
type iStorage interface {
	RestoreIndexes(ctx context.Context) (map[string]int32, error)
	SaveIndexLoc(ctx context.Context, index string, server int32) error
}

//go:generate mockgen -destination=mocks/servers/mock_servers.go -package=mocks . iServers
type iServers interface {
	GetServer() (*servers.Server, bool)
	Len() int32
	AddServer(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(number int32)
	GetServers() []string
	DeepSearch(ctx context.Context, key string) (string, error)
	GetServerByID(number int32) (*servers.Server, bool)
	Exists(number int32) bool
	SetToAll(ctx context.Context, key string, val string, uniques bool) []int32
	DelFromAll(ctx context.Context, key string) (atLeastOnce bool)
}

type UseCase struct {
	servers iServers
	logger  logger.ILogger
	storage iStorage

	indexes map[string]int32
	mu      sync.RWMutex

	pool chan struct{} // TODO: ADD TO CONFIG
}

func New(ctx context.Context, repository iStorage, logger logger.ILogger) (*UseCase, error) {
	s, err := servers.New()
	if err != nil {
		return nil, err
	}

	indexes, err := repository.RestoreIndexes(ctx)
	if err != nil {
		return nil, err
	}

	return &UseCase{
		servers: s,
		storage: repository,
		logger:  logger,
		indexes: indexes,
		pool:    make(chan struct{}, 30000), // TODO: MOVE TO CONFIG
	}, nil
}
