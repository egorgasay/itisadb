package domains

import (
	"context"
	"itisadb/internal/models"
)

type Core interface {
	Get(ctx context.Context, server *int32, key string) (string, models.RAM, error)
	Set(ctx context.Context, server *int32, key string, val string, uniques bool) (int32, models.RAM, error)
	Delete(ctx context.Context, server *int32, key string) (models.RAM, error)

	Object(ctx context.Context, server *int32, name string) (int32, models.RAM, error)
	ObjectToJSON(ctx context.Context, server *int32, name string) (string, models.RAM, error)
	DeleteObject(ctx context.Context, server *int32, name string) (models.RAM, error)
	IsObject(ctx context.Context, server *int32, name string) (bool, models.RAM, error)
	Size(ctx context.Context, server *int32, name string) (uint64, models.RAM, error)
	AttachToObject(ctx context.Context, server *int32, dst string, src string) (models.RAM, error)

	GetFromObject(ctx context.Context, server *int32, object string, key string) (string, models.RAM, error)
	SetToObject(ctx context.Context, server *int32, object string, key string, val string, uniques bool) (int32, models.RAM, error)
	DeleteAttr(ctx context.Context, server *int32, attr string, object string) (models.RAM, error)

	Connect(address string, available uint64, total uint64) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string
	Authenticate(ctx context.Context, login string, password string) (string, error)
}
