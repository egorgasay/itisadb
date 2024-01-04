package logic

import (
	"context"
	"fmt"
	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
)

type Logic struct {
	ram     gost.Mutex[models.RAM]
	storage domains.Storage
	cfg     config.Config
	tlogger domains.TransactionLogger
	logger  *zap.Logger
}

func NewLogic(storage domains.Storage, cfg config.Config, tlogger domains.TransactionLogger, logger *zap.Logger) *Logic {
	return &Logic{
		ram:     gost.NewMutex(models.RAM{}),
		storage: storage,
		cfg:     cfg,
		tlogger: tlogger,
		logger:  logger,
	}
}

func (l *Logic) GetOne(_ context.Context, userID int, key string, _ models.GetOptions) (res gost.Result[models.Value]) {
	v, err := l.storage.Get(key)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, v.Level) {
		return res.Err(constants.ErrForbidden)
	}

	return res.Ok(v)
}

func (l *Logic) DelOne(_ context.Context, userID int, key string, _ models.DeleteOptions) (res gost.Result[gost.Nothing]) {
	v, err := l.storage.Get(key)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, v.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if err := l.storage.Delete(key); err != nil {
		l.logger.Warn("failed to delete", zap.Error(err))
		return res.ErrNew(0, 0, err.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDelete(key)
	}

	return res.Ok(gost.Nothing{})
}

func (l *Logic) SetOne(_ context.Context, userID int, key string, val string, opt models.SetOptions) (res gost.Result[int32]) {
	if !l.hasPermission(userID, opt.Level) {
		return res.Err(constants.ErrForbidden)
	}

	v, err := l.storage.Get(key)
	if err == nil {
		if !l.hasPermission(userID, v.Level) {
			return res.Err(constants.ErrForbidden)
		}

		if opt.Unique || v.ReadOnly {
			return res.Err(constants.ErrAlreadyExists)
		}
	}

	err = l.storage.Set(key, val, opt)
	if err != nil {
		return res.ErrNew(0, 0, err.Error()) // TODO:
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteSet(key, val, opt)
	}

	return res.Ok(constants.MainStorageNumber)
}

func (l *Logic) HasPermissionToObject(userID int, name string) (res gost.Result[bool]) {
	info, err := l.storage.GetObjectInfo(name)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	return res.Ok(l.hasPermission(userID, info.Level))
}

func (l *Logic) NewObject(_ context.Context, userID int, name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing]) {
	if !l.hasPermission(userID, opts.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if err := l.storage.CreateObject(name, opts); err != nil {
		return res.ErrNewUnknown(fmt.Sprintf("can't create object: %w", err)) // TODO: ??
	}

	info := models.ObjectInfo{
		Server: constants.MainStorageNumber,
		Level:  opts.Level,
	}

	l.storage.AddObjectInfo(name, info) // TODO: maybe you should union Create + AddObjectInfo? and keep all information about object in one place?

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteCreateObject(name)
		l.tlogger.WriteAddObjectInfo(name, info)
	}

	return res.Ok(gost.Nothing{})
}

func (l *Logic) SetToObject(_ context.Context, userID int, object string, key string, value string, opts models.SetToObjectOptions) (res gost.Result[gost.Nothing]) {
	info, err := l.storage.GetObjectInfo(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error()) // TODO: ??
	}

	if !l.hasPermission(userID, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	err = l.storage.SetToObject(object, key, value, opts)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteSetToObject(object, key, value)
	}

	return res.Ok(gost.Nothing{})
}

func (l *Logic) GetFromObject(_ context.Context, userID int, object, key string, _ models.GetFromObjectOptions) (res gost.Result[string]) {
	info, err := l.storage.GetObjectInfo(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	v, err := l.storage.GetFromObject(object, key)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	return res.Ok(v)
}

func (l *Logic) ObjectToJSON(_ context.Context, userID int, object string, _ models.ObjectToJSONOptions) (res gost.Result[string]) {
	info, err := l.storage.GetObjectInfo(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	v, err := l.storage.ObjectToJSON(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	return res.Ok(v)
}

func (l *Logic) ObjectSize(_ context.Context, userID int, object string, _ models.SizeOptions) (res gost.Result[uint64]) {
	info, err := l.storage.GetObjectInfo(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	v, err := l.storage.Size(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	return res.Ok(v)
}

func (l *Logic) DeleteObject(_ context.Context, userID int, object string, _ models.DeleteObjectOptions) (res gost.ResultN) {
	info, err := l.storage.GetObjectInfo(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	err = l.storage.DeleteObject(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	l.storage.DeleteObjectInfo(object)

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDeleteObject(object)
		l.tlogger.WriteDeleteObjectInfo(object)
	}

	return res.Ok()
}

func (l *Logic) AttachToObject(_ context.Context, userID int, dst, src string, _ models.AttachToObjectOptions) (res gost.ResultN) {
	infoDst, err := l.storage.GetObjectInfo(dst)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	infoSrc, err := l.storage.GetObjectInfo(src)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, infoDst.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if !l.hasPermission(userID, infoSrc.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if err := l.storage.AttachToObject(dst, src); err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteAttach(dst, src)
	}

	return res.Ok()
}

func (l *Logic) ObjectDeleteKey(_ context.Context, userID int, object, key string, _ models.DeleteAttrOptions) (res gost.ResultN) {
	info, err := l.storage.GetObjectInfo(object)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if !l.hasPermission(userID, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if err := l.storage.DeleteAttr(object, key); err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDeleteAttr(object, key)
	}

	return res.Ok()
}
