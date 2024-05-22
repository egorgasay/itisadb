package domains

import (
	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

type TransactionLogger interface {
	Run()
	Err() <-chan error
	Stop() error
	Restore(r Restorer) error
	WriteSet(key string, value string, opts models.SetOptions)
	WriteDelete(key string)
	WriteSetToObject(name string, key string, val string, opts models.SetToObjectOptions)
	WriteCreateObject(name string, info models.ObjectInfo)
	WriteDeleteObject(name string)
	WriteAttach(dst string, src string)
	WriteDeleteAttr(name string, key string)
	WriteNewUser(user models.User)
	WriteDeleteUser(login string)
}

type Restorer interface {
	Set(key, value string, opts models.SetOptions) gost.ResultN
	Delete(key string) gost.ResultN
	SetToObject(name, key, value string, opts models.SetToObjectOptions) gost.ResultN
	DeleteObject(name string) gost.ResultN
	CreateObject(name string, opts models.ObjectOptions) gost.ResultN
	AttachToObject(dst, src string) gost.ResultN
	NewUser(user models.User) (r gost.ResultN)
	DeleteUser(login string) (r gost.Result[bool])
	AddObjectInfo(name string, info models.ObjectInfo)
	DeleteObjectInfo(name string)
	DeleteAttr(object, key string) gost.ResultN
}
