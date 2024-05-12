package domains

import (
	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

//go:generate mockgen -destination=mocks/storage/mock_storage.go -package=mocks . Storage
type Storage interface {
	CommonStorage
	ObjectsStorage
	UserStorage
}

type CommonStorage interface {
	Set(key string, val string, opts models.SetOptions) gost.ResultN
	Get(key string) (r gost.Option[models.Value])
	DeleteIfExists(key string)
	Delete(key string) gost.ResultN
}

type ObjectsStorage interface {
	/*
	   Common operations with objects
	*/

	CreateObject(name string, opts models.ObjectOptions) (r gost.ResultN)
	DeleteObject(name string) (r gost.ResultN)
	SetToObject(name string, key string, value string, opts models.SetToObjectOptions) gost.ResultN
	GetFromObject(name string, key string) (r gost.Option[string]) // TODO: impl -> models.Value

	/*
	   PRO operations with objects
	*/

	ObjectToJSON(name string) (r gost.Result[string])
	Size(name string) (r gost.Result[uint64])
	IsObject(name string) bool
	DeleteAttr(name string, key string) gost.ResultN

	/*
		ObjectInfo operations
	*/

	AddObjectInfo(name string, info models.ObjectInfo)
	DeleteObjectInfo(name string)
	GetObjectInfo(name string) (r gost.Option[models.ObjectInfo])

	/*
	   Complicated and not fully implemented or tested
	*/

	AttachToObject(dst string, src string) gost.ResultN
}

type UserStorage interface {
	NewUser(user models.User) (r gost.ResultN)
	GetUserByName(username string) (r gost.Result[models.User])
	DeleteUser(username string) (r gost.Result[bool])
	SaveUser(user models.User) (r gost.ResultN)
	GetUserLevel(username string) (r gost.Result[models.Level])
	GetUsersFromChangeID(id uint64) gost.Result[[]models.User]
	GetUserChangeID() uint64
	SetUserChangeID(id uint64)
}
