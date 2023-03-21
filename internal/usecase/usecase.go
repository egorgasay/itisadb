package usecase

import (
	"errors"
	"github.com/egorgasay/grpc-storage/internal/storage"
)

type UseCase struct {
	storage storage.Storage
}

var NotFoundErr = errors.New("the value does not exist")

func New(storage *storage.Storage) *UseCase {
	return &UseCase{storage: *storage}
}

func (uc *UseCase) Set(key string, val string) {
	uc.storage.Mu.Lock()
	defer uc.storage.Mu.Unlock()
	uc.storage.RAMStorage[key] = val
}

func (uc *UseCase) Get(key string) (string, error) {
	uc.storage.Mu.RLock()
	defer uc.storage.Mu.RUnlock()
	val, ok := uc.storage.RAMStorage[key]
	if !ok {
		return "", NotFoundErr
	}

	return val, nil
}

func (uc *UseCase) Save() error {
	return nil
}

func (uc *UseCase) Load() error {
	return nil
}
