package servers

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/pkg/api/storage"
)

// =============== server ====================== //

var ErrAlreadyExists = errors.New("already exists")

type Server struct {
	Tries     uint
	storage   storage.StorageClient
	Available uint64
	Total     uint64
	Number    int32
}

func (s *Server) Set(ctx context.Context, Key, Value string, unique bool) (*storage.SetResponse, error) {
	res, err := s.storage.Set(ctx, &storage.SetRequest{Key: Key, Value: Value, Unique: unique})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return nil, ErrAlreadyExists
		}
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

func (s *Server) GetIndex(ctx context.Context, name string) (*storage.GetIndexResponse, error) {
	gir, err := s.storage.GetIndex(ctx, &storage.GetIndexRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return gir, err

}

func (s *Server) GetFromIndex(ctx context.Context, name, Key string) (*storage.GetResponse, error) {
	gfir, err := s.storage.GetFromIndex(ctx, &storage.GetFromIndexRequest{
		Key:  Key,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return gfir, err

}

func (s *Server) SetToIndex(ctx context.Context, name, Key, Value string, unique bool) (*storage.SetResponse, error) {
	stir, err := s.storage.SetToIndex(ctx, &storage.SetToIndexRequest{
		Key:    Key,
		Value:  Value,
		Name:   name,
		Unique: unique,
	})
	if err != nil {
		return nil, err
	}

	return stir, err

}

func (s *Server) NewIndex(ctx context.Context, name string) error {
	_, err := s.storage.NewIndex(ctx, &storage.NewIndexRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Size(ctx context.Context, name string) (*storage.IndexSizeResponse, error) {
	r, err := s.storage.Size(ctx, &storage.IndexSizeRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
