//package keeper
//
//import (
//	"go.uber.org/zap"
//	"itisadb/internal/domains"
//	"itisadb/internal/models"
//)
//
//type Keeper struct {
//	domains.Storage
//	logger    *zap.Logger
//	isTLogger bool
//	tLogger   domains.TransactionLogger
//}
//
//func New(Storage domains.Storage, l *zap.Logger, tlogger domains.TransactionLogger) (*Keeper, error) {
//
//	l = l.Named("KEEPER")
//	if tlogger == nil {
//		l.Info("Transaction logger disabled")
//		return &Keeper{Storage: Storage, logger: l, isTLogger: false}, nil
//	}
//
//	return &Keeper{Storage: Storage, logger: l, isTLogger: true, tLogger: tlogger}, nil
//}
//
//func (k *Keeper) Set(key, val string, uniques bool) error {
//	err := k.Storage.Set(key, val, uniques)
//	if err != nil {
//		k.logger.Warn("Failed to set", zap.Error(err))
//		return err
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteSet(key, val)
//	}
//	return nil
//}
//
//var ram = models.RAM{}
//
//func (k *Keeper) SetToObject(name, key, val string, uniques bool) error {
//	err := k.Storage.SetToObject(name, key, val, uniques)
//	if err != nil {
//		k.logger.Warn("Failed to set to object", zap.Error(err))
//		return err
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteSetToObject(name, key, val)
//	}
//	return nil
//}
//
//func (k *Keeper) ObjectToJSON(name string) (string, error) {
//	object, err := k.Storage.ObjectToJSON(name)
//	if err != nil {
//		k.logger.Warn("Failed to get object", zap.Error(err))
//		return "", err
//	}
//	return object, nil
//}
//
//func (k *Keeper) CreateObject(name string) error {
//	err := k.Storage.CreateObject(name)
//	if err != nil {
//		k.logger.Warn("Failed to create object", zap.Error(err))
//		return err
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteCreateObject(name)
//	}
//	return nil
//}
//
//func (k *Keeper) Size(name string) (uint64, error) {
//	size, err := k.Storage.Size(name)
//	if err != nil {
//		k.logger.Warn("Failed to get size", zap.Error(err))
//		return 0, err
//	}
//	return size, nil
//}
//
//func (k *Keeper) DeleteObject(name string) error {
//	err := k.Storage.DeleteObject(name)
//	if err != nil {
//		k.logger.Warn("Failed to delete object", zap.Error(err))
//		return err
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteDeleteObject(name)
//	}
//	return nil
//}
//
//func (k *Keeper) AttachToObject(dst, src string) error {
//	err := k.Storage.AttachToObject(dst, src)
//	if err != nil {
//		k.logger.Warn("Failed to attach object", zap.Error(err))
//		return nil
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteAttach(dst, src)
//	}
//	return nil
//}
//
//func (k *Keeper) DeleteIfExists(key string) {
//	k.Storage.DeleteIfExists(key)
//
//	if k.isTLogger {
//		k.tLogger.WriteDelete(key)
//	}
//}
//
//func (k *Keeper) Delete(key string) error {
//	err := k.Storage.Delete(key)
//	if err != nil {
//		k.logger.Warn("Failed to delete", zap.Error(err))
//		return err
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteDelete(key)
//	}
//	return nil
//}
//
//func (k *Keeper) DeleteAttr(name, key string) error {
//	err := k.Storage.DeleteAttr(name, key)
//	if err != nil {
//		k.logger.Warn("Failed to delete attr", zap.Error(err))
//		return err
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteDeleteAttr(name, key)
//	}
//	return err
//}
//
//func (k *Keeper) CreateUser(user models.User) (id int, err error) {
//	id, err = k.Storage.CreateUser(user)
//	if err != nil {
//		k.logger.Warn("Failed to create user", zap.Error(err))
//		return 0, err
//	}
//
//	if k.isTLogger {
//		k.tLogger.WriteCreateUser(user)
//	}
//
//	return id, nil
//}
