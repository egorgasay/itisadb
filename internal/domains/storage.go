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
	Set(key string, val string, opts models.SetOptions) error
	Get(key string) (models.Value, error)
	DeleteIfExists(key string)
	Delete(key string) error
}

type ObjectsStorage interface {
	/*
	   Common operations with objects
	*/

	CreateObject(name string, opts models.ObjectOptions) (err error)
	DeleteObject(name string) error
	SetToObject(name string, key string, value string, opts models.SetToObjectOptions) error
	GetFromObject(name string, key string) (string, error) // TODO: impl -> models.Value

	/*
	   PRO operations with objects
	*/

	ObjectToJSON(name string) (string, error)
	Size(name string) (uint64, error)
	IsObject(name string) bool
	DeleteAttr(name string, key string) error

	/*
		ObjectInfo operations
	*/

	AddObjectInfo(name string, info models.ObjectInfo)
	DeleteObjectInfo(name string)
	GetObjectInfo(name string) (models.ObjectInfo, error)

	/*
	   Complicated and not fully implemented or tested
	*/

	AttachToObject(dst string, src string) error
}

type UserStorage interface {
	NewUser(user models.User) (r gost.Result[int])
	GetUserByID(id int) (r gost.Result[models.User])
	GetUserByName(username string) (r gost.Result[models.User])
	GetUserIDByName(username string) (r gost.Result[int])
	DeleteUser(id int) (r gost.Result[bool])
	SaveUser(user models.User) (r gost.ResultN)
	GetUserLevel(id int) (r gost.Result[models.Level])
	GetUsersFromChangeID(id uint64) gost.Result[[]models.User]
	GetUserChangeID() uint64
	SetUserChangeID(id uint64)
}
