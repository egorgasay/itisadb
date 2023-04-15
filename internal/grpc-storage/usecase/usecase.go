package usecase

import (
	"github.com/pbnjay/memory"
	"github.com/sourcegraph/conc/pool"
	"grpc-storage/internal/grpc-storage/storage"
	"grpc-storage/pkg/logger"
)

type UseCase struct {
	storage *storage.Storage
	logger  logger.ILogger
	pool    *pool.Pool
}

func New(storage *storage.Storage, logger logger.ILogger) *UseCase {
	gpool := pool.New()
	gpool.WithMaxGoroutines(20000)
	return &UseCase{storage: storage, logger: logger, pool: gpool}
}

func (uc *UseCase) Set(key string, val string) RAM {
	uc.storage.Set(key, val)
	uc.pool.Go(func() { uc.storage.WriteSet(key, val) })
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
