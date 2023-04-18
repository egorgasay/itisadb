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

func (s *Server) Set(ctx context.Context, Key, Value string) (*storage.SetResponse, error) {
	res, err := s.storage.Set(ctx, &storage.SetRequest{Key: Key, Value: Value})
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *Server) Get(ctx context.Context, Key string) (*storage.GetResponse, error) {
	gr, err := s.storage.Get(ctx, &storage.GetRequest{Key: Key})
	if err != nil {
		return nil, err
	}

	return gr, err

}
