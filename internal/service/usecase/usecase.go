package usecase

import (
	"context"
	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	"sync"
)

type UseCase struct {
	ram     gost.Mutex[models.RAM]
	storage domains.Storage
	cfg     config.Config
	tlogger domains.TransactionLogger
}

func NewUseCase(storage domains.Storage, cfg config.Config, tlogger domains.TransactionLogger) *UseCase {
	return &UseCase{
		ram:     gost.NewMutex(models.RAM{}),
		storage: storage,
		cfg:     cfg,
		tlogger: tlogger,
	}
}

func (s *UseCase) SetRAM(ram models.RAM) {
	defer s.ram.Release()
	*s.ram.BorrowMut() = ram
}

func (s *UseCase) RAM() models.RAM {
	defer s.ram.Release()
	return s.ram.Borrow().Read()
}

func (s *UseCase) Find(ctx context.Context, key string, out chan<- string, once *sync.Once, opts models.GetOptions) {
	//TODO implement me
	panic("implement me")
}

func (s *UseCase) GetOne(ctx context.Context, key string, opts ...itisadb.GetOptions) (res gost.Result[string]) {
	v, err := s.storage.Get(key)
	if err != nil {
		return res.ErrNew(0, 0, err.Error())
	}

	return res.Ok(v.Value)
}

func (s *UseCase) DelOne(ctx context.Context, key string, opts ...itisadb.DeleteOptions) gost.Result[gost.Nothing] {
	//TODO implement me
	panic("implement me")
}

func (c *UseCase) SetOne(ctx context.Context, key string, val string, opt models.SetOptions) (res gost.Result[int32]) {
	v, err := c.storage.Get(key)
	if err == nil && (opt.Unique || v.ReadOnly) {
		return res.Err(constants.ErrAlreadyExists)
	}

	err = c.storage.Set(key, val, opt)
	if err != nil {
		return res.ErrNew(0, 0, err.Error()) // TODO:
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteSet(key, val, opt)
	}

	return res.Ok(constants.MainStorageNumber)
}

func (s *UseCase) NewObject(ctx context.Context, name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing]) {
	//TODO implement me
	panic("implement me")
}

func (s *UseCase) SetToObject(ctx context.Context, object string, key string, value string, opts models.SetToObjectOptions) (res gost.Result[gost.Nothing]) {
	//TODO implement me
	panic("implement me")
}

func (s *UseCase) GetFromObject(ctx context.Context, object string, key string, opts models.GetFromObjectOptions) (res gost.Result[string]) {
	//TODO implement me
	panic("implement me")
}
