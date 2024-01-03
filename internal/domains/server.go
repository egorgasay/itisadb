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
	Tries() uint32
	IncTries()
	ResetTries()

	appLogic
}

type appLogic interface {
	GetOne(ctx context.Context, userID int, key string, opt models.GetOptions) (res gost.Result[models.Value])
	DelOne(ctx context.Context, userID int, key string, opt models.DeleteOptions) gost.Result[gost.Nothing]
	SetOne(ctx context.Context, userID int, key string, val string, opt models.SetOptions) (res gost.Result[int32])

	NewObject(ctx context.Context, userID int, name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing])
	SetToObject(ctx context.Context, userID int, object string, key string, value string, opts models.SetToObjectOptions) (res gost.Result[gost.Nothing])
	GetFromObject(ctx context.Context, userID int, object string, key string, opts models.GetFromObjectOptions) (res gost.Result[string])

	ObjectToJSON(ctx context.Context, name string, opts models.ObjectToJSONOptions) (res gost.Result[string])
}

//Set(ctx context.Context, key string, value string, opts models.SetOptions) error
//Get(ctx context.Context, key string, opts models.GetOptions) (*api.GetResponse, error)
//GetFromObject(ctx context.Context, name string, key string, opts models.GetFromObjectOptions) (*api.GetFromObjectResponse, error)
//SetToObject(ctx context.Context, name string, key string, value string, opts models.SetToObjectOptions) error
//NewObject(ctx context.Context, name string, opts models.ObjectOptions) error
//Size(ctx context.Context, name string, opts models.SizeOptions) (*api.ObjectSizeResponse, error)
//DeleteObject(ctx context.Context, name string, opts models.DeleteObjectOptions) error
//Delete(ctx context.Context, Key string, opts models.DeleteOptions) error
//AttachToObject(ctx context.Context, dst string, src string, opts models.AttachToObjectOptions) error
//DeleteAttr(ctx context.Context, attr string, object string, opts models.DeleteAttrOptions) error
