package usecase

import (
	"github.com/pbnjay/memory"
	"itisadb/internal/grpc-storage/storage"
	"itisadb/pkg/logger"
)

type UseCase struct {
	storage *storage.Storage
	logger  logger.ILogger
}

func New(storage *storage.Storage, logger logger.ILogger) *UseCase {
	return &UseCase{storage: storage, logger: logger}
}

func (uc *UseCase) Set(key, val string, uniques bool) (RAM, error) {
	err := uc.storage.Set(key, val, uniques)
	uc.storage.WriteSet(key, val)
	return RAMUsage(), err
}

func (uc *UseCase) SetToIndex(name, key, val string) (RAM, error) {
	err := uc.storage.SetToIndex(name, key, val)
	uc.storage.WriteSet(key, val)
	return RAMUsage(), err
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

func (uc *UseCase) Get(key string) (RAM, string, error) {
	s, err := uc.storage.Get(key)
	return RAMUsage(), s, err
}

func (uc *UseCase) GetFromIndex(name, key string) (RAM, string, error) {
	s, err := uc.storage.GetFromIndex(name, key)
	return RAMUsage(), s, err
}

func (uc *UseCase) GetIndex(name string) (RAM, map[string]string, error) {
	index, err := uc.storage.GetIndex(name)
	return RAMUsage(), index, err
}

func (uc *UseCase) Save() {
	uc.logger.Info("Saving values ...")
	err := uc.storage.Save()
	if err != nil {
		uc.logger.Warn(err.Error())
	}
}

func (uc *UseCase) NewIndex(name string) (RAM, error) {
	return RAMUsage(), uc.storage.CreateIndex(name)
}

func (uc *UseCase) Size(name string) (RAM, uint64, error) {
	size, err := uc.storage.Size(name)
	return RAMUsage(), size, err
}

func (uc *UseCase) DeleteIndex(name string) (RAM, error) {
	return RAMUsage(), uc.storage.DeleteIndex(name)
}

func (uc *UseCase) AttachToIndex(dst, src string) (RAM, error) {
	return RAMUsage(), uc.storage.AttachToIndex(dst, src)
}

func (uc *UseCase) Delete(key string) RAM {
	uc.storage.Delete(key)
	return RAMUsage()
}
