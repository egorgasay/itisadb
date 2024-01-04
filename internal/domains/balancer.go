package domains

import (
	"context"
	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

//go:generate mockgen -destination=mocks/balancer/mock_servers.go -package=mocks . Servers
type Servers interface {
	Len() int32
	AddServer(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(number int32)
	GetServers() []string
	GetServer(number int32) (Server, bool)
	Exists(number int32) bool

	// TODO: may be we should use Iter instead, because Servers != buisness logic
	SetToAll(ctx context.Context, userID int, key string, val string, opts models.SetOptions) []int32

	// TODO: may be we should use Iter instead, because Servers != buisness logic
	DelFromAll(ctx context.Context, userID int, key string, opts models.DeleteOptions) (atLeastOnce bool)

	// TODO: may be we should use Iter instead, because Servers != buisness logic
	DeepSearch(ctx context.Context, userID int, key string, opts models.GetOptions) (res gost.Result[gost.Pair[int32, models.Value]])
}
