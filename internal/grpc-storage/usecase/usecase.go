package usecase

import (
	"github.com/egorgasay/grpc-storage/internal/grpc-storage/storage"
	"github.com/pbnjay/memory"
	"log"
)

type UseCase struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *UseCase {
	return &UseCase{storage: storage}
}

func (uc *UseCase) Set(key string, val string) RAM {
	uc.storage.Set(key, val)
	return RAMUsage()
}

type RAM struct {
	Total     uint64
	Available uint64
}

// RAMUsage outputs the current, total and OS memory being used.
func RAMUsage() RAM {
	return RAM{
		Total:     memory.TotalMemory() / 1024 / 1024,
		Available: memory.FreeMemory() / 1024 / 1024,
	}
}

func (uc *UseCase) Get(key string) (string, error) {
	return uc.storage.Get(key)
}

func (uc *UseCase) Save() {
	log.Println("Saving values ...")
	uc.storage.Save()
}
