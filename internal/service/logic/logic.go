package logic

import (
	"context"
	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	"sync"
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

func (l *Logic) Find(ctx context.Context, key string, out chan<- string, once *sync.Once, opts models.GetOptions) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) GetOne(ctx context.Context, key string, opt models.GetOptions) (res gost.Result[string]) {
	v, err := l.storage.Get(key)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	return res.Ok(v.Value)
}

func (l *Logic) DelOne(ctx context.Context, key string, opt models.DeleteOptions) (res gost.Result[gost.Nothing]) {
	if err := l.storage.Delete(key); err != nil {
		l.logger.Warn("failed to delete", zap.Error(err))
		return res.ErrNew(0, 0, err.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDelete(key)
	}

	return res.Ok(gost.Nothing{})
}

func (l *Logic) SetOne(ctx context.Context, key string, val string, opt models.SetOptions) (res gost.Result[int32]) {
	v, err := l.storage.Get(key)
	if err == nil && (opt.Unique || v.ReadOnly) {
		return res.Err(constants.ErrAlreadyExists)
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

func (l *Logic) NewObject(ctx context.Context, name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing]) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) SetToObject(ctx context.Context, object string, key string, value string, opts models.SetToObjectOptions) (res gost.Result[gost.Nothing]) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) GetFromObject(ctx context.Context, object string, key string, opts models.GetFromObjectOptions) (res gost.Result[string]) {
	//TODO implement me
	panic("implement me")
}
