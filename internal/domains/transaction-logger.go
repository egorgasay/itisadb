package domains

import (
	"context"
	"itisadb/internal/models"
)

type TransactionLogger interface {
	Run()
	Err() <-chan error
	Stop() error
	Restore(r Restorer) error
	WriteSet(key string, value string)
	WriteDelete(key string)
	WriteSetToObject(name string, key string, val string)
	WriteCreateObject(name string)
	WriteDeleteObject(name string)
	WriteAttach(dst string, src string)
	WriteDeleteAttr(name string, key string)
	WriteCreateUser(user models.User)
	RestoreObjects(ctx context.Context) (map[string]int32, error)
	SaveObjectLoc(ctx context.Context, object string, server int32) error
}

type Restorer interface {
	Set(key, value string, uniques bool) error
	Delete(key string) error
	SetToObject(name, key, value string, uniques bool) error
	DeleteObject(name string) error
	CreateObject(name string) error
	AttachToObject(dst, src string) error
	CreateUser(user models.User) (int, error)
}
