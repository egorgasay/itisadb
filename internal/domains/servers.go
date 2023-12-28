package domains

import (
	"context"
	"itisadb/internal/models"
	"itisadb/internal/service/servers"
)

//go:generate mockgen -destination=mocks/servers/mock_servers.go -package=mocks . Servers
type Servers interface {
	GetServer() (servers.Server, bool)
	Len() int32
	AddServer(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(number int32)
	GetServers() []string
	DeepSearch(ctx context.Context, key string, opts models.GetOptions) (string, error)
	GetServerByID(number int32) (servers.Server, bool)
	Exists(number int32) bool
	SetToAll(ctx context.Context, key string, val string, opts models.SetOptions) []int32
	DelFromAll(ctx context.Context, key string, opts models.DeleteOptions) (atLeastOnce bool)
}
