package domains

import "context"

type Core interface {
	Get(ctx context.Context, server *int32, key string) (string, error)
	Set(ctx context.Context, server *int32, key string, val string, uniques bool) (int32, error)
	Delete(ctx context.Context, server *int32, key string) error

	Object(ctx context.Context, server *int32, name string) (int32, error)
	ObjectToJSON(ctx context.Context, server *int32, name string) (string, error)
	DeleteObject(ctx context.Context, server *int32, name string) error
	IsObject(ctx context.Context, server *int32, name string) (bool, error)
	Size(ctx context.Context, server *int32, name string) (uint64, error)
	AttachToObject(ctx context.Context, server *int32, dst string, src string) error

	GetFromObject(ctx context.Context, server *int32, object string, key string) (string, error)
	SetToObject(ctx context.Context, server *int32, object string, key string, val string, uniques bool) (int32, error)
	DeleteAttr(ctx context.Context, server *int32, attr string, object string) error

	Connect(address string, available uint64, total uint64) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string
	Authenticate(ctx context.Context, login string, password string) (string, error)
}
