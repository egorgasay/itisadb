package domains

import "context"

//go:generate mockgen -destination=mocks/storage/mock_storage.go -package=mocks . Storage
type Storage interface {
	Set(key string, val string, unique bool) error
	Get(key string) (string, error)
	DeleteIfExists(key string)
	Delete(key string) error

	SetToObject(name string, key string, value string, uniques bool) error
	GetFromObject(name string, key string) (string, error)
	AttachToObject(dst string, src string) error
	DeleteObject(name string) error
	CreateObject(name string) (err error)
	ObjectToJSON(name string) (string, error)
	Size(name string) (uint64, error)
	IsObject(name string) bool
	DeleteAttr(name string, key string) error

	RestoreObjects(ctx context.Context) (map[string]int32, error)
	SaveObjectLoc(ctx context.Context, object string, server int32) error
}
