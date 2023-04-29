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

func (uc *UseCase) Save() {
	uc.logger.Info("Saving values ...")
	err := uc.storage.Save()
	if err != nil {
		uc.logger.Warn(err.Error())
	}
}
