package mocks

import "context"

//go:generate mockgen -destination=mocks/usecase/mock_usecase.go -package=mocks . IUseCase
type IUseCase interface {
	Index(ctx context.Context, name string) (int32, error)
	GetFromIndex(ctx context.Context, index string, key string, serverNumber int32) (string, error)
	SetToIndex(ctx context.Context, index string, key string, val string, uniques bool) (int32, error)
	GetIndex(ctx context.Context, name string) (map[string]string, error)
	IsIndex(ctx context.Context, name string) (bool, error)
	Size(ctx context.Context, name string) (uint64, error)
	DeleteIndex(ctx context.Context, name string) error
	AttachToIndex(ctx context.Context, dst string, src string) error
	DeleteAttr(ctx context.Context, attr string, index string) error
	Set(ctx context.Context, key string, val string, serverNumber int32, uniques bool) (int32, error)
	Get(ctx context.Context, key string, serverNumber int32) (string, error)
	Connect(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string
	Delete(ctx context.Context, key string, num int32) error
}
