package domains

import (
	"context"
	"itisadb/internal/models"
)

type Core interface {
	Get(ctx context.Context, userID uint, key string, opts models.GetOptions) (string, error)
	Set(ctx context.Context, userID uint, key, val string, opts models.SetOptions) (int32, error)
	Delete(ctx context.Context, userID uint, key string, opts models.DeleteOptions) error

	Object(ctx context.Context, userID uint, name string, opts models.ObjectOptions) (int32, error)
	ObjectToJSON(ctx context.Context, userID uint, name string, opts models.ObjectToJSONOptions) (string, error)
	DeleteObject(ctx context.Context, userID uint, name string, opts models.DeleteObjectOptions) error
	IsObject(ctx context.Context, userID uint, name string, opts models.IsObjectOptions) (bool, error)
	Size(ctx context.Context, userID uint, name string, opts models.SizeOptions) (uint64, error)
	AttachToObject(ctx context.Context, userID uint, dst, src string, opts models.AttachToObjectOptions) error

	GetFromObject(ctx context.Context, userID uint, object, key string, opts models.GetFromObjectOptions) (string, error)
	SetToObject(ctx context.Context, userID uint, object, key, val string, opts models.SetToObjectOptions) (int32, error)
	DeleteAttr(ctx context.Context, userID uint, attr, object string, opts models.DeleteAttrOptions) error

	Connect(address string, available uint64, total uint64) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string

	Authenticate(ctx context.Context, login, password string) (string, error)
	CreateUser(ctx context.Context, userID uint, user models.User) error
	DeleteUser(ctx context.Context, userID uint, login string) error
	ChangePassword(ctx context.Context, userID uint, login, password string) error
	ChangeLevel(ctx context.Context, userID uint, login string, level int8) error
}
