package usecase

import (
	"fmt"
	"github.com/pbnjay/memory"
	"itisadb/internal/grpc-storage/storage"
	tlogger "itisadb/internal/grpc-storage/transaction-logger"
	"itisadb/pkg/logger"
)

type UseCase struct {
	storage   storage.IStorage
	logger    logger.ILogger
	isTLogger bool
	tLogger   *tlogger.TransactionLogger
}

//go:generate mockgen -destination=mocks/usecase/mock_usecase.go -package=mocks . IUseCase
type IUseCase interface {
	Set(key string, val string, uniques bool) (RAM, error)
	SetToIndex(name string, key string, val string, uniques bool) (RAM, error)
	Get(key string) (RAM, string, error)
	GetFromIndex(name string, key string) (RAM, string, error)
	GetIndex(name string) (RAM, map[string]string, error)
	NewIndex(name string) (RAM, error)
	Size(name string) (RAM, uint64, error)
	DeleteIndex(name string) (RAM, error)
	AttachToIndex(dst string, src string) (RAM, error)
	DeleteIfExists(key string) RAM
	Delete(key string) (RAM, error)
	DeleteAttr(name string, key string) (RAM, error)
}

func New(storage storage.IStorage, logger logger.ILogger, enableTLogger bool) (*UseCase, error) {
	if !enableTLogger {
		logger.Info("Transaction logger disabled")
		return &UseCase{storage: storage, logger: logger, isTLogger: false}, nil
	}

	tl, err := tlogger.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction logger: %w", err)
	}

	logger.Info("Transaction logger enabled")

	logger.Info("Starting recovery from transaction logger")
	if err = tl.Restore(storage); err != nil {
		return nil, fmt.Errorf("failed to restore from transaction logger: %w", err)
	}
	logger.Info("Transaction logger recovery completed")

	tl.Run()
	logger.Info("Transaction logger started")

	return &UseCase{storage: storage, logger: logger, isTLogger: true, tLogger: tl}, nil
}

func (uc *UseCase) Set(key, val string, uniques bool) (RAM, error) {
	err := uc.storage.Set(key, val, uniques)
	if err != nil {
		return RAMUsage(), err
	}

	if uc.isTLogger {
		uc.tLogger.WriteSet(key, val)
	}
	return RAMUsage(), err
}

func (uc *UseCase) SetToIndex(name, key, val string, uniques bool) (RAM, error) {
	err := uc.storage.SetToIndex(name, key, val, uniques)
	if err != nil {
		return RAMUsage(), err
	}

	if uc.isTLogger {
		uc.tLogger.WriteSetToIndex(name, key, val)
	}
	return RAMUsage(), err
}

type RAM struct {
	Total     uint64
	Available uint64
}

// RAMUsage outputs the current, total and OS memory being used.
func RAMUsage() RAM {
	// TODO: do not call it every time
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
	index, err := uc.storage.GetIndex(name, "")
	return RAMUsage(), index, err
}

func (uc *UseCase) NewIndex(name string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.CreateIndex(name)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteCreateIndex(name)
	}
	return r, err
}

func (uc *UseCase) Size(name string) (RAM, uint64, error) {
	size, err := uc.storage.Size(name)
	return RAMUsage(), size, err
}

func (uc *UseCase) DeleteIndex(name string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.DeleteIndex(name)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteDeleteIndex(name)
	}
	return r, err
}

func (uc *UseCase) AttachToIndex(dst, src string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.AttachToIndex(dst, src)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteAttach(dst, src)
	}
	return r, err
}

func (uc *UseCase) DeleteIfExists(key string) RAM {
	uc.storage.DeleteIfExists(key)

	if uc.isTLogger {
		uc.tLogger.WriteDelete(key)
	}
	return RAMUsage()
}

func (uc *UseCase) Delete(key string) (RAM, error) {
	err := uc.storage.Delete(key)
	if uc.isTLogger {
		uc.tLogger.WriteDelete(key)
	}
	return RAMUsage(), err
}

func (uc *UseCase) DeleteAttr(name, key string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.DeleteAttr(name, key)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteDeleteAttr(name, key)
	}
	return r, err
}
