package mocks

import "context"

//go:generate mockgen -destination=mock_usecase.go -package=mocks . IUseCase
type IUseCase interface {
	Get(ctx context.Context, key string, serverNumber int32) (string, error)
	Set(ctx context.Context, key string, val string, serverNumber int32, uniques bool) (int32, error)
	Delete(ctx context.Context, key string, num int32) error

	Index(ctx context.Context, name string) (int32, error)
	IndexToJSON(ctx context.Context, name string) (string, error)
	DeleteIndex(ctx context.Context, name string) error
	IsIndex(ctx context.Context, name string) (bool, error)
	Size(ctx context.Context, name string) (uint64, error)
	AttachToIndex(ctx context.Context, dst string, src string) error

	GetFromIndex(ctx context.Context, index string, key string, serverNumber int32) (string, error)
	SetToIndex(ctx context.Context, index string, key string, val string, uniques bool) (int32, error)
	DeleteAttr(ctx context.Context, attr string, index string) error

	Connect(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string
}
