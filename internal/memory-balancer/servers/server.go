package servers

import (
	"context"
	"grpc-storage/pkg/api/storage"
)

// =============== server ====================== //

type Server struct {
	Tries     uint
	storage   storage.StorageClient
	Available uint64
	Total     uint64
	Number    int32
}

func (s *Server) Set(ctx context.Context, key, value string) (*storage.SetResponse, error) {
	return s.storage.Set(ctx, &storage.SetRequest{Key: key, Value: value})
}

func (s *Server) Get(ctx context.Context, key string) (*storage.GetResponse, error) {
	return s.storage.Get(ctx, &storage.GetRequest{Key: key})
}
