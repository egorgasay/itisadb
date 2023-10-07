package domains

import "itisadb/internal/models"

type Keeper interface {
	Set(key string, val string, uniques bool) error
	SetToObject(name string, key string, val string, uniques bool) error
	Get(key string) (string, error)
	GetFromObject(name string, key string) (string, error)
	ObjectToJSON(name string) (string, error)
	NewObject(name string) error
	Size(name string) (uint64, error)
	DeleteObject(name string) error
	AttachToObject(dst string, src string) error
	DeleteIfExists(key string) models.RAM
	Delete(key string) error
	DeleteAttr(name string, key string) error
}
