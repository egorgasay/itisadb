package storage

type IStorage interface {
	Save(key, value string) error
	Load(key string) (string, error)
	Delete(key string) error
	SavePairToIndex(index, key, value string) error
}
