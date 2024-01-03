package domains

import (
	"context"
	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

//go:generate mockgen -destination=mocks/balancer/mock_servers.go -package=mocks . Servers
type Servers interface {
	GetServer() (Server, bool)
	Len() int32
	AddServer(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(number int32)
	GetServers() []string
	DeepSearch(ctx context.Context, key string, opts models.GetOptions) (models.Value, error)
	GetServerByID(number int32) (Server, bool)
	Exists(number int32) bool
	SetToAll(ctx context.Context, key string, val string, opts models.SetOptions) []int32
	DelFromAll(ctx context.Context, key string, opts models.DeleteOptions) (atLeastOnce bool)
	OnServerError(cl Server, err *gost.Error)
}
