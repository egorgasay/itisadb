package domains

import (
	"itisadb/internal/models"
)

type TransactionLogger interface {
	Run()
	Err() <-chan error
	Stop() error
	Restore(r Restorer) error
	WriteSet(key string, value string, opts models.SetOptions)
	WriteDelete(key string)
	WriteSetToObject(name string, key string, val string)
	WriteCreateObject(name string)
	WriteDeleteObject(name string)
	WriteAttach(dst string, src string)
	WriteDeleteAttr(name string, key string)
	WriteNewUser(user models.User)
	WriteAddObjectInfo(name string, info models.ObjectInfo)
	WriteDeleteObjectInfo(name string)
	WriteDeleteUser(login string)
}

type Restorer interface {
	Set(key, value string, opts models.SetOptions) error
	Delete(key string) error
	SetToObject(name, key, value string, opts models.SetToObjectOptions) error
	DeleteObject(name string) error
	CreateObject(name string, opts models.ObjectOptions) error
	AttachToObject(dst, src string) error
	NewUser(user models.User) (int, error)
	AddObjectInfo(name string, info models.ObjectInfo)
	DeleteObjectInfo(name string)
}
