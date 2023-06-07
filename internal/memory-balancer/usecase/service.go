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
	_ = iota * -1
	dbOnly
	all
	allAndDB
)

//go:generate mockgen -destination=mocks/storage/mock_storage.go -package=mocks . iStorage
type iStorage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, val string) error
	SetUnique(ctx context.Context, key, val string) error
	Delete(ctx context.Context, key string) error
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
}

type UseCase struct {
	servers iServers
	logger  logger.ILogger
	storage iStorage

	// TODO: add copy to disk
	indexes map[string]int32
	mu      sync.RWMutex

	pool chan struct{} // TODO: ADD TO CONFIG
}

func New(repository iStorage, logger logger.ILogger) (*UseCase, error) {
	s, err := servers.New()
	if err != nil {
		return nil, err
	}
	return &UseCase{
		servers: s,
		storage: repository,
		logger:  logger,
		indexes: make(map[string]int32, 10000),
		pool:    make(chan struct{}, 30000),
	}, nil
}
