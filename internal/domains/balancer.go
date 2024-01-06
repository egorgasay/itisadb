package domains

import (
	"context"

	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

type Balancer interface {
	Get(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.GetOptions) (models.Value, error)
	Set(ctx context.Context, claims gost.Option[models.UserClaims], key, val string, opts models.SetOptions) (int32, error)
	Delete(ctx context.Context, claims gost.Option[models.UserClaims], key string, opts models.DeleteOptions) error

	Object(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.ObjectOptions) (int32, error)
	ObjectToJSON(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.ObjectToJSONOptions) (string, error)
	DeleteObject(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.DeleteObjectOptions) error
	IsObject(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.IsObjectOptions) (bool, error)
	Size(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.SizeOptions) (uint64, error)
	AttachToObject(ctx context.Context, claims gost.Option[models.UserClaims], dst, src string, opts models.AttachToObjectOptions) error

	GetFromObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key string, opts models.GetFromObjectOptions) (string, error)
	SetToObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key, val string, opts models.SetToObjectOptions) (int32, error)
	DeleteAttr(ctx context.Context, claims gost.Option[models.UserClaims], attr, object string, opts models.DeleteAttrOptions) error

	Connect(ctx context.Context, address string, available uint64, total uint64) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string

	Authenticate(ctx context.Context, login, password string) (string, error)
	CreateUser(ctx context.Context, claims gost.Option[models.UserClaims], user models.User) error
	DeleteUser(ctx context.Context, claims gost.Option[models.UserClaims], login string) error
	ChangePassword(ctx context.Context, claims gost.Option[models.UserClaims], login, password string) error
	ChangeLevel(ctx context.Context, claims gost.Option[models.UserClaims], login string, level models.Level) error
	CalculateRAM(ctx context.Context) (res gost.Result[models.RAM])
}
