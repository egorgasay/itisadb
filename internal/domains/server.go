package domains

import (
	"context"

	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

type Server interface {
	RAM() models.RAM
	RefreshRAM(ctx context.Context) (res gost.Result[gost.Nothing])
	Number() int32
	IsOffline() bool
	Reconnect(ctx context.Context) (res gost.ResultN)

	appLogic
}

type appLogic interface {
	GetOne(ctx context.Context, claims gost.Option[models.UserClaims], key string, opt models.GetOptions) (res gost.Result[models.Value])
	DelOne(ctx context.Context, claims gost.Option[models.UserClaims], key string, opt models.DeleteOptions) gost.Result[gost.Nothing]
	SetOne(ctx context.Context, claims gost.Option[models.UserClaims], key string, val string, opt models.SetOptions) (res gost.Result[int32])

	NewObject(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing])
	SetToObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key, value string, opts models.SetToObjectOptions) (res gost.Result[gost.Nothing])
	GetFromObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key string, opts models.GetFromObjectOptions) (res gost.Result[string])

	ObjectToJSON(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.ObjectToJSONOptions) (res gost.Result[string])
	ObjectSize(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.SizeOptions) (res gost.Result[uint64])
	DeleteObject(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.DeleteObjectOptions) gost.ResultN
	AttachToObject(ctx context.Context, claims gost.Option[models.UserClaims], dst, src string, opts models.AttachToObjectOptions) gost.ResultN
	ObjectDeleteKey(ctx context.Context, claims gost.Option[models.UserClaims], object, key string, opts models.DeleteAttrOptions) gost.ResultN
}
