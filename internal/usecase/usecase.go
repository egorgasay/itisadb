package usecase

import "github.com/egorgasay/grpc-storage/internal/storage"

type UseCase struct {
	storage storage.Storage
}

func New(storage *storage.Storage) *UseCase {
	return &UseCase{storage: *storage}
}
