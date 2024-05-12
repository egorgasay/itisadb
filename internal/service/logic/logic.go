package logic

import (
	"context"

	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"

	"github.com/egorgasay/gost"
	"go.uber.org/zap"
)

type Logic struct {
	storage  domains.Storage
	cfg      config.Config
	tlogger  domains.TransactionLogger
	security domains.SecurityService

	logger *zap.Logger
}

func NewLogic(
	storage domains.Storage,
	cfg config.Config,
	tlogger domains.TransactionLogger,
	logger *zap.Logger,
	security domains.SecurityService,
) *Logic {

	r := storage.GetUserByName("itisadb")
	if r.IsErr() {
		logger.Info("creating default user")

		r := storage.NewUser(
			models.User{
				Login:    "itisadb",
				Password: "itisadb",
				Level:    constants.SecretLevel,
				Active:   true,
			},
		)

		if r.IsErr() {
			logger.Error("failed to create default user", zap.Error(r.Error()))
		}
	}

	r = storage.GetUserByName("demo")
	if r.IsErr() {
		logger.Info("creating demo user")

		r := storage.NewUser(
			models.User{
				Login:    "demo",
				Password: "demo",
				Level:    constants.DefaultLevel,
				Active:   true,
			},
		)

		if r.IsErr() {
			logger.Error("failed to create demo user", zap.Error(r.Error()))
		}
	}

	return &Logic{
		storage:  storage,
		cfg:      cfg,
		tlogger:  tlogger,
		logger:   logger,
		security: security,
	}
}

func (l *Logic) GetOne(_ context.Context, claims gost.Option[models.UserClaims], key string, _ models.GetOptions) (res gost.Result[models.Value]) {
	v := l.storage.Get(key)
	if v.IsNone() {
		return res.Err(constants.ErrNotFound)
	}

	value := v.Unwrap()

	if !l.security.HasPermission(claims, value.Level) {
		return res.Err(constants.ErrForbidden)
	}

	return res.Ok(value)
}

func (l *Logic) DelOne(_ context.Context, claims gost.Option[models.UserClaims], key string, _ models.DeleteOptions) (res gost.ResultN) {
	v := l.storage.Get(key)
	if v.IsNone() {
		return res.Err(constants.ErrNotFound)
	}

	value := v.Unwrap()

	if !l.security.HasPermission(claims, value.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if r := l.storage.Delete(key); r.IsErr() {
		l.logger.Warn("failed to delete", zap.Error(r.Error()))
		return res.Err(r.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDelete(key)
	}

	return res.Ok()
}

func (l *Logic) SetOne(_ context.Context, claims gost.Option[models.UserClaims], key string, val string, opt models.SetOptions) (res gost.Result[int32]) {
	if !l.security.HasPermission(claims, opt.Level) {
		return res.Err(constants.ErrForbidden)
	}

	r := l.storage.Get(key)
	if r.IsSome() {
		if !l.security.HasPermission(claims, r.Unwrap().Level) {
			return res.Err(constants.ErrForbidden)
		}

		if opt.Unique || r.Unwrap().ReadOnly {
			return res.Err(constants.ErrAlreadyExists)
		}
	}

	rSet := l.storage.Set(key, val, opt)
	if rSet.IsErr() {
		return res.Err(rSet.Error())
	}

	if l.cfg.TransactionLogger.On {
		opt.Encrypt = opt.Level == constants.SecretLevel
		l.tlogger.WriteSet(key, val, opt)
	}

	return res.Ok(constants.LocalServerNumber)
}

func (l *Logic) HasPermissionToObject(claims gost.Option[models.UserClaims], name string) (res gost.Result[bool]) {
	infoR := l.storage.GetObjectInfo(name)
	if infoR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	return res.Ok(l.security.HasPermission(claims, infoR.Unwrap().Level))
}

func (l *Logic) NewObject(_ context.Context, claims gost.Option[models.UserClaims], name string, opts models.ObjectOptions) (res gost.ResultN) {
	if !l.security.HasPermission(claims, opts.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if r := l.storage.CreateObject(name, opts); r.IsErr() {
		return res.Err(r.Error())
	}

	info := models.ObjectInfo{
		Server: constants.LocalServerNumber,
		Level:  opts.Level,
	}

	l.storage.AddObjectInfo(name, info) // TODO: maybe you should union Create + AddObjectInfo? and keep all information about object in one place?

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteCreateObject(name, info)
	}

	return res.Ok()
}

func (l *Logic) SetToObject(ctx context.Context, claims gost.Option[models.UserClaims], object string, key string, value string, opts models.SetToObjectOptions) (res gost.ResultN) {
	infoR := l.storage.GetObjectInfo(object)
	if infoR.IsNone() {
		if r := l.NewObject(ctx, claims, object, models.ObjectOptions{
			Level: constants.DefaultLevel,
		}); r.IsErr() {
			return res.Err(r.Error())
		}
	}

	info := infoR.Unwrap()

	if !l.security.HasPermission(claims, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if r := l.storage.SetToObject(object, key, value, opts); r.IsErr() {
		return res.Err(r.Error())
	}

	if l.cfg.TransactionLogger.On {
		opts.Encrypt = info.Level == constants.SecretLevel
		l.tlogger.WriteSetToObject(object, key, value, opts)
	}

	return res.Ok()
}

func (l *Logic) GetFromObject(_ context.Context, claims gost.Option[models.UserClaims], object, key string, _ models.GetFromObjectOptions) (res gost.Result[string]) {
	infoR := l.storage.GetObjectInfo(object)
	if infoR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	info := infoR.Unwrap()

	if !l.security.HasPermission(claims, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	r := l.storage.GetFromObject(object, key)
	if r.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	return res.Ok(r.Unwrap())
}

func (l *Logic) ObjectToJSON(_ context.Context, claims gost.Option[models.UserClaims], object string, _ models.ObjectToJSONOptions) (res gost.Result[string]) {
	infoR := l.storage.GetObjectInfo(object)
	if infoR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	info := infoR.Unwrap()

	if !l.security.HasPermission(claims, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	r := l.storage.ObjectToJSON(object)
	if r.IsErr() {
		return res.Err(r.Error())
	}

	return res.Ok(r.Unwrap())
}

func (l *Logic) ObjectSize(_ context.Context, claims gost.Option[models.UserClaims], object string, _ models.SizeOptions) (res gost.Result[uint64]) {
	infoR := l.storage.GetObjectInfo(object)
	if infoR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	info := infoR.Unwrap()

	if !l.security.HasPermission(claims, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	r := l.storage.Size(object)
	if r.IsErr() {
		return res.Err(r.Error())
	}

	return res.Ok(r.Unwrap())
}

func (l *Logic) DeleteObject(_ context.Context, claims gost.Option[models.UserClaims], object string, _ models.DeleteObjectOptions) (res gost.ResultN) {
	infoR := l.storage.GetObjectInfo(object)
	if infoR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	if !l.security.HasPermission(claims, infoR.Unwrap().Level) {
		return res.Err(constants.ErrForbidden)
	}

	if r := l.storage.DeleteObject(object); r.IsErr() {
		return res.Err(r.Error())
	}

	l.storage.DeleteObjectInfo(object)

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDeleteObject(object)
	}

	return res.Ok()
}

func (l *Logic) AttachToObject(_ context.Context, claims gost.Option[models.UserClaims], dst, src string, _ models.AttachToObjectOptions) (res gost.ResultN) {
	infoDstR := l.storage.GetObjectInfo(dst)
	if infoDstR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	infoSrcR := l.storage.GetObjectInfo(src)
	if infoSrcR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	if !l.security.HasPermission(claims, infoDstR.Unwrap().Level) {
		return res.Err(constants.ErrForbidden)
	}

	if !l.security.HasPermission(claims, infoSrcR.Unwrap().Level) {
		return res.Err(constants.ErrForbidden)
	}

	if r := l.storage.AttachToObject(dst, src); r.IsErr() {
		return res.Err(r.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteAttach(dst, src)
	}

	return res.Ok()
}

func (l *Logic) ObjectDeleteKey(_ context.Context, claims gost.Option[models.UserClaims], object, key string, _ models.DeleteAttrOptions) (res gost.ResultN) {
	infoR := l.storage.GetObjectInfo(object)
	if infoR.IsNone() {
		return res.Err(constants.ErrObjectNotFound)
	}

	info := infoR.Unwrap()

	if !l.security.HasPermission(claims, info.Level) {
		return res.Err(constants.ErrForbidden)
	}

	if r := l.storage.DeleteAttr(object, key); r.IsErr() {
		return res.Err(r.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDeleteAttr(object, key)
	}

	return res.Ok()
}

func (l *Logic) IsObject(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.IsObjectOptions) (res gost.Result[bool]) {
	return res.Ok(l.storage.IsObject(object))
}
