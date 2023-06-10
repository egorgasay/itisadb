package usecase

import (
	"github.com/pbnjay/memory"
	"itisadb/internal/grpc-storage/storage"
	"itisadb/pkg/logger"
)

type UseCase struct {
	storage storage.IStorage
	logger  logger.ILogger
	dirDB   bool
}

//go:generate mockgen -destination=mocks/usecase/mock_usecase.go -package=mocks . IUseCase
type IUseCase interface {
	Set(key string, val string, uniques bool) (RAM, error)
	SetToIndex(name string, key string, val string, uniques bool) (RAM, error)
	Get(key string) (RAM, string, error)
	GetFromIndex(name string, key string) (RAM, string, error)
	GetIndex(name string) (RAM, map[string]string, error)
	Save()
	NewIndex(name string) (RAM, error)
	Size(name string) (RAM, uint64, error)
	DeleteIndex(name string) (RAM, error)
	AttachToIndex(dst string, src string) (RAM, error)
	DeleteIfExists(key string) RAM
	Delete(key string) (RAM, error)
	DeleteAttr(name string, key string) (RAM, error)
}

func New(storage storage.IStorage, logger logger.ILogger) *UseCase {
	return &UseCase{storage: storage, logger: logger, dirDB: true}
}

func (uc *UseCase) Set(key, val string, uniques bool) (RAM, error) {
	err := uc.storage.Set(key, val, uniques)
	if err != nil {
		return RAMUsage(), err
	}

	if !uc.storage.NoTLogger() {
		uc.storage.WriteSet(key, val)
	}
	return RAMUsage(), err
}

func (uc *UseCase) SetToIndex(name, key, val string, uniques bool) (RAM, error) {
	err := uc.storage.SetToIndex(name, key, val, uniques)
	// uc.storage.WriteSet(name+"/"+key, val) TODO: add to index
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
	if err != nil && uc.dirDB && err == storage.ErrNotFound {
		s, err = uc.storage.GetFromDisk(key)
	}
	return RAMUsage(), s, err
}

func (uc *UseCase) GetFromIndex(name, key string) (RAM, string, error) {
	s, err := uc.storage.GetFromIndex(name, key)
	if err != nil && uc.dirDB && err == storage.ErrIndexNotFound {
		s, err = uc.storage.GetFromDiskIndex(name, key)
	}
	return RAMUsage(), s, err
}

func (uc *UseCase) GetIndex(name string) (RAM, map[string]string, error) {
	index, err := uc.storage.GetIndex(name, "")
	return RAMUsage(), index, err
}

func (uc *UseCase) Save() {
	uc.logger.Info("Saving values ...")
	err := uc.storage.Save()
	if err != nil {
		uc.logger.Warn(err.Error())
	}
}

func (uc *UseCase) NewIndex(name string) (RAM, error) {
	return RAMUsage(), uc.storage.CreateIndex(name)
}

func (uc *UseCase) Size(name string) (RAM, uint64, error) {
	size, err := uc.storage.Size(name)
	return RAMUsage(), size, err
}

func (uc *UseCase) DeleteIndex(name string) (RAM, error) {
	return RAMUsage(), uc.storage.DeleteIndex(name)
}

func (uc *UseCase) AttachToIndex(dst, src string) (RAM, error) {
	return RAMUsage(), uc.storage.AttachToIndex(dst, src)
}

func (uc *UseCase) DeleteIfExists(key string) RAM {
	uc.storage.DeleteIfExists(key)

	if uc.logger != nil {
		uc.storage.WriteDelete(key)
	}
	return RAMUsage()
}

func (uc *UseCase) Delete(key string) (RAM, error) {
	err := uc.storage.Delete(key)
	if err == nil && uc.logger != nil {
		uc.storage.WriteDelete(key)
	}
	return RAMUsage(), err
}

func (uc *UseCase) DeleteAttr(name, key string) (RAM, error) {
	return RAMUsage(), uc.storage.DeleteAttr(name, key)
}
