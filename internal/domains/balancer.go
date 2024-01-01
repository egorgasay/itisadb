package domains

import (
	"context"
	"itisadb/internal/models"
)

//go:generate mockgen -destination=mocks/balancer/mock_servers.go -package=mocks . Servers
type Balancer interface {
	GetServer() (Server, bool)
	Len() int32
	AddServer(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(number int32)
	GetServers() []string
	DeepSearch(ctx context.Context, key string, opts models.GetOptions) (string, error)
	GetServerByID(number int32) (Server, bool)
	Exists(number int32) bool
	SetToAll(ctx context.Context, key string, val string, opts models.SetOptions) []int32
	DelFromAll(ctx context.Context, key string, opts models.DeleteOptions) (atLeastOnce bool)
}
