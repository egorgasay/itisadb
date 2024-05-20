package domains

import (
	"context"

	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

//go:generate mockgen -destination=mocks/balancer/mock_servers.go -package=mocks . Servers
type Servers interface {
	Len() int32
	AddServer(ctx context.Context, address string, force bool) (int32, error)
	Disconnect(number int32)
	GetServersInfo() []string
	GetServer(number int32) (Server, bool)
	Exists(number int32) bool

	// TODO: may be we should use Iter instead, because Servers != buisness logic
	SetToAll(ctx context.Context, claims gost.Option[models.UserClaims], key string, val string, opts models.SetOptions) []int32

	// TODO: may be we should use Iter instead, because Servers != buisness logic
	DelFromAll(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.DeleteOptions) (atLeastOnce bool)

	// TODO: may be we should use Iter instead, because Servers != buisness logic
	DeepSearch(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.GetOptions) (res gost.Result[gost.Pair[int32, models.Value]])

	Iter(func(Server) error) error
}
