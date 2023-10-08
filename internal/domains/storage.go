package domains

import (
	"itisadb/internal/models"
)

//go:generate mockgen -destination=mocks/storage/mock_storage.go -package=mocks . Storage
type Storage interface {
	Set(key string, val string, opts models.SetOptions) error
	Get(key string) (string, error)
	DeleteIfExists(key string)
	Delete(key string) error

	SetToObject(name string, key string, value string, opts models.SetToObjectOptions) error
	GetFromObject(name string, key string) (string, error)
	AttachToObject(dst string, src string) error
	DeleteObject(name string) error
	CreateObject(name string, opts models.ObjectOptions) (err error)
	ObjectToJSON(name string) (string, error)
	Size(name string) (uint64, error)
	IsObject(name string) bool
	DeleteAttr(name string, key string) error

	CreateUser(user models.User) (id int, err error)
	GetUserByID(id int) (models.User, error)
	GetUserByName(username string) (int, models.User, error)
	DeleteUser(id int) error
	SaveUser(id int, user models.User) error
}
