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

	appLogic
}

type appLogic interface {
	GetOne(ctx context.Context, userID int, key string, opt models.GetOptions) (res gost.Result[models.Value])
	DelOne(ctx context.Context, userID int, key string, opt models.DeleteOptions) gost.Result[gost.Nothing]
	SetOne(ctx context.Context, userID int, key string, val string, opt models.SetOptions) (res gost.Result[int32])

	NewObject(ctx context.Context, userID int, name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing])
	SetToObject(ctx context.Context, userID int, object, key, value string, opts models.SetToObjectOptions) (res gost.Result[gost.Nothing])
	GetFromObject(ctx context.Context, userID int, object, key string, opts models.GetFromObjectOptions) (res gost.Result[string])

	ObjectToJSON(ctx context.Context, userID int, name string, opts models.ObjectToJSONOptions) (res gost.Result[string])
	ObjectSize(ctx context.Context, userID int, object string, opts models.SizeOptions) (res gost.Result[uint64])
	DeleteObject(ctx context.Context, userID int, object string, opts models.DeleteObjectOptions) gost.ResultN
	AttachToObject(ctx context.Context, userID int, dst, src string, opts models.AttachToObjectOptions) gost.ResultN
	ObjectDeleteKey(ctx context.Context, userID int, object, key string, opts models.DeleteAttrOptions) gost.ResultN
}
