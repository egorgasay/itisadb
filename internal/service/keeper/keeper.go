package keeper

import (
	"go.uber.org/zap"
	"itisadb/internal/domains"
	"itisadb/internal/models"
)

type Keeper struct {
	storage   domains.Storage
	logger    *zap.Logger
	isTLogger bool
	tLogger   domains.TransactionLogger
}

func New(storage domains.Storage, l *zap.Logger, tlogger domains.TransactionLogger) (*Keeper, error) {

	l = l.Named("KEEPER")
	if tlogger == nil {
		l.Info("Transaction logger disabled")
		return &Keeper{storage: storage, logger: l, isTLogger: false}, nil
	}

	return &Keeper{storage: storage, logger: l, isTLogger: true, tLogger: tlogger}, nil
}

func (k *Keeper) Set(key, val string, uniques bool) error {
	err := k.storage.Set(key, val, uniques)
	if err != nil {
		k.logger.Warn("Failed to set", zap.Error(err))
		return err
	}

	if k.isTLogger {
		k.tLogger.WriteSet(key, val)
	}
	return nil
}

var ram = models.RAM{}

func (k *Keeper) SetToObject(name, key, val string, uniques bool) error {
	err := k.storage.SetToObject(name, key, val, uniques)
	if err != nil {
		k.logger.Warn("Failed to set to object", zap.Error(err))
		return err
	}

	if k.isTLogger {
		k.tLogger.WriteSetToObject(name, key, val)
	}
	return nil
}

func (k *Keeper) Get(key string) (string, error) {
	s, err := k.storage.Get(key)
	if err != nil {
		k.logger.Warn("Failed to get", zap.Error(err))
		return "", err
	}
	return s, nil
}

func (k *Keeper) GetFromObject(name, key string) (string, error) {
	s, err := k.storage.GetFromObject(name, key)
	if err != nil {
		k.logger.Warn("Failed to get from object", zap.Error(err))
		return "", err
	}
	return s, nil
}

func (k *Keeper) ObjectToJSON(name string) (string, error) {
	object, err := k.storage.ObjectToJSON(name)
	if err != nil {
		k.logger.Warn("Failed to get object", zap.Error(err))
		return "", err
	}
	return object, nil
}

func (k *Keeper) NewObject(name string) error {
	err := k.storage.CreateObject(name)
	if err != nil {
		k.logger.Warn("Failed to create object", zap.Error(err))
		return err
	}

	if k.isTLogger {
		k.tLogger.WriteCreateObject(name)
	}
	return nil
}

func (k *Keeper) Size(name string) (uint64, error) {
	size, err := k.storage.Size(name)
	if err != nil {
		k.logger.Warn("Failed to get size", zap.Error(err))
		return 0, err
	}
	return size, nil
}

func (k *Keeper) DeleteObject(name string) error {
	err := k.storage.DeleteObject(name)
	if err != nil {
		k.logger.Warn("Failed to delete object", zap.Error(err))
		return err
	}

	if k.isTLogger {
		k.tLogger.WriteDeleteObject(name)
	}
	return nil
}

func (k *Keeper) AttachToObject(dst, src string) error {
	err := k.storage.AttachToObject(dst, src)
	if err != nil {
		k.logger.Warn("Failed to attach object", zap.Error(err))
		return nil
	}

	if k.isTLogger {
		k.tLogger.WriteAttach(dst, src)
	}
	return nil
}

func (k *Keeper) DeleteIfExists(key string) models.RAM {
	k.storage.DeleteIfExists(key)

	if k.isTLogger {
		k.tLogger.WriteDelete(key)
	}
	return ram.Update()
}

func (k *Keeper) Delete(key string) error {
	err := k.storage.Delete(key)
	if err != nil {
		k.logger.Warn("Failed to delete", zap.Error(err))
		return err
	}

	if k.isTLogger {
		k.tLogger.WriteDelete(key)
	}
	return nil
}

func (k *Keeper) DeleteAttr(name, key string) error {
	err := k.storage.DeleteAttr(name, key)
	if err != nil {
		k.logger.Warn("Failed to delete attr", zap.Error(err))
		return err
	}

	if k.isTLogger {
		k.tLogger.WriteDeleteAttr(name, key)
	}
	return err
}

func (k *Keeper) CreateUser(user models.User) error {
	err := k.storage.CreateUser(user)
	if err != nil {
		k.logger.Warn("Failed to create user", zap.Error(err))
		return err
	}

	if k.isTLogger {
		k.tLogger.WriteCreateUser(user)
	}

	return nil
}
