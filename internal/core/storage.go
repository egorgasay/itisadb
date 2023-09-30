package core

import (
	"fmt"
	"github.com/pbnjay/memory"
	"itisadb/internal/grpc-storage/storage"
	tlogger "itisadb/internal/grpc-storage/transaction-logger"
	"itisadb/pkg/logger"
)

type Keeper struct {
	storage   storage.IStorage
	logger    logger.ILogger
	isTLogger bool
	tLogger   *tlogger.TransactionLogger
}

//go:generate mockgen -destination=mocks/usecase/mock_usecase.go -package=mocks . IUseCase
type ICore interface {
	Set(key string, val string, uniques bool) (RAM, error)
	SetToObject(name string, key string, val string, uniques bool) (RAM, error)
	Get(key string) (RAM, string, error)
	GetFromObject(name string, key string) (RAM, string, error)
	ObjectToJSON(name string) (RAM, string, error)
	NewObject(name string) (RAM, error)
	Size(name string) (RAM, uint64, error)
	DeleteObject(name string) (RAM, error)
	AttachToObject(dst string, src string) (RAM, error)
	DeleteIfExists(key string) RAM
	Delete(key string) (RAM, error)
	DeleteAttr(name string, key string) (RAM, error)
}

func newKeeper(storage storage.IStorage, logger logger.ILogger, enableTLogger bool) (*Keeper, error) {
	if !enableTLogger {
		logger.Info("Transaction logger disabled")
		return &Keeper{storage: storage, logger: logger, isTLogger: false}, nil
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

	return &Keeper{storage: storage, logger: logger, isTLogger: true, tLogger: tl}, nil
}

func (uc *Keeper) Set(key, val string, uniques bool) (RAM, error) {
	err := uc.storage.Set(key, val, uniques)
	if err != nil {
		return RAMUsage(), err
	}

	if uc.isTLogger {
		uc.tLogger.WriteSet(key, val)
	}
	return RAMUsage(), err
}

func (uc *Keeper) SetToObject(name, key, val string, uniques bool) (RAM, error) {
	err := uc.storage.SetToObject(name, key, val, uniques)
	if err != nil {
		return RAMUsage(), err
	}

	if uc.isTLogger {
		uc.tLogger.WriteSetToObject(name, key, val)
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

func (uc *Keeper) Get(key string) (RAM, string, error) {
	s, err := uc.storage.Get(key)
	return RAMUsage(), s, err
}

func (uc *Keeper) GetFromObject(name, key string) (RAM, string, error) {
	s, err := uc.storage.GetFromObject(name, key)
	return RAMUsage(), s, err
}

func (uc *Keeper) ObjectToJSON(name string) (RAM, string, error) {
	object, err := uc.storage.ObjectToJSON(name)
	return RAMUsage(), object, err
}

func (uc *Keeper) NewObject(name string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.CreateObject(name)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteCreateObject(name)
	}
	return r, err
}

func (uc *Keeper) Size(name string) (RAM, uint64, error) {
	size, err := uc.storage.Size(name)
	return RAMUsage(), size, err
}

func (uc *Keeper) DeleteObject(name string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.DeleteObject(name)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteDeleteObject(name)
	}
	return r, err
}

func (uc *Keeper) AttachToObject(dst, src string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.AttachToObject(dst, src)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteAttach(dst, src)
	}
	return r, err
}

func (uc *Keeper) DeleteIfExists(key string) RAM {
	uc.storage.DeleteIfExists(key)

	if uc.isTLogger {
		uc.tLogger.WriteDelete(key)
	}
	return RAMUsage()
}

func (uc *Keeper) Delete(key string) (RAM, error) {
	err := uc.storage.Delete(key)
	if uc.isTLogger {
		uc.tLogger.WriteDelete(key)
	}
	return RAMUsage(), err
}

func (uc *Keeper) DeleteAttr(name, key string) (RAM, error) {
	r, err := RAMUsage(), uc.storage.DeleteAttr(name, key)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteDeleteAttr(name, key)
	}
	return r, err
}
