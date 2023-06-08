package usecase

import "itisadb/internal/dirdb/storage"

type IUseCase interface {
}

type UseCase struct {
	storage storage.IStorage
}
