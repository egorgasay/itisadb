package core

import (
	"fmt"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	tlogger "itisadb/internal/transaction-logger"
	"itisadb/pkg/logger"
)

type Keeper struct {
	storage   domains.Storage
	logger    logger.ILogger
	isTLogger bool
	tLogger   *tlogger.TransactionLogger
}

func newKeeper(storage domains.Storage, logger logger.ILogger, enableTLogger bool) (*Keeper, error) {
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

func (uc *Keeper) Set(key, val string, uniques bool) (models.RAM, error) {
	err := uc.storage.Set(key, val, uniques)
	if err != nil {
		return ram.Update(), err
	}

	if uc.isTLogger {
		uc.tLogger.WriteSet(key, val)
	}
	return ram.Update(), nil
}

var ram = models.RAM{}

func (uc *Keeper) SetToObject(name, key, val string, uniques bool) (models.RAM, error) {
	err := uc.storage.SetToObject(name, key, val, uniques)
	if err != nil {
		return ram.Update(), err
	}

	if uc.isTLogger {
		uc.tLogger.WriteSetToObject(name, key, val)
	}
	return ram.Update(), err
}

func (uc *Keeper) Get(key string) (models.RAM, string, error) {
	s, err := uc.storage.Get(key)
	return ram.Update(), s, err
}

func (uc *Keeper) GetFromObject(name, key string) (models.RAM, string, error) {
	s, err := uc.storage.GetFromObject(name, key)
	return ram.Update(), s, err
}

func (uc *Keeper) ObjectToJSON(name string) (models.RAM, string, error) {
	object, err := uc.storage.ObjectToJSON(name)
	return ram.Update(), object, err
}

func (uc *Keeper) NewObject(name string) (models.RAM, error) {
	r, err := ram.Update(), uc.storage.CreateObject(name)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteCreateObject(name)
	}
	return r, err
}

func (uc *Keeper) Size(name string) (models.RAM, uint64, error) {
	size, err := uc.storage.Size(name)
	return ram.Update(), size, err
}

func (uc *Keeper) DeleteObject(name string) (models.RAM, error) {
	r, err := ram.Update(), uc.storage.DeleteObject(name)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteDeleteObject(name)
	}
	return r, err
}

func (uc *Keeper) AttachToObject(dst, src string) (models.RAM, error) {
	r, err := ram.Update(), uc.storage.AttachToObject(dst, src)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteAttach(dst, src)
	}
	return r, err
}

func (uc *Keeper) DeleteIfExists(key string) models.RAM {
	uc.storage.DeleteIfExists(key)

	if uc.isTLogger {
		uc.tLogger.WriteDelete(key)
	}
	return ram.Update()
}

func (uc *Keeper) Delete(key string) (models.RAM, error) {
	err := uc.storage.Delete(key)
	if uc.isTLogger {
		uc.tLogger.WriteDelete(key)
	}
	return ram.Update(), err
}

func (uc *Keeper) DeleteAttr(name, key string) (models.RAM, error) {
	r, err := ram.Update(), uc.storage.DeleteAttr(name, key)
	if err != nil {
		return r, err
	}

	if uc.isTLogger {
		uc.tLogger.WriteDeleteAttr(name, key)
	}
	return r, err
}
