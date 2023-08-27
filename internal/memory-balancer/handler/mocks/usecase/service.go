package mocks

import "context"

//go:generate mockgen -destination=mock_usecase.go -package=mocks . IUseCase
type IUseCase interface {
	Get(ctx context.Context, key string, serverNumber int32) (string, error)
	Set(ctx context.Context, key string, val string, serverNumber int32, uniques bool) (int32, error)
	Delete(ctx context.Context, key string, num int32) error

	Object(ctx context.Context, name string) (int32, error)
	ObjectToJSON(ctx context.Context, name string) (string, error)
	DeleteObject(ctx context.Context, name string) error
	IsObject(ctx context.Context, name string) (bool, error)
	Size(ctx context.Context, name string) (uint64, error)
	AttachToObject(ctx context.Context, dst string, src string) error

	GetFromObject(ctx context.Context, object string, key string, serverNumber int32) (string, error)
	SetToObject(ctx context.Context, object string, key string, val string, uniques bool) (int32, error)
	DeleteAttr(ctx context.Context, attr string, object string) error

	Connect(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string
}
