package servers

import (
	"context"
	"grpc-storage/pkg/api/storage"
)

// =============== server ====================== //

type key struct {
	key string
}

type keyValue struct {
	key   string
	value string
}

type Server struct {
	Tries     uint
	storage   storage.StorageClient
	Available uint64
	Total     uint64
	Number    int32

	sc storage.Storage_SetClient
	gc storage.Storage_GetClient

	inSet  chan keyValue
	outSet chan *storage.SetResponse

	inGet  chan key
	outGet chan *storage.GetResponse
}

func (s *Server) Set(ctx context.Context, Key, Value string) (*storage.SetResponse, error) {
	sc, err := s.storage.Set(context.Background())
	if err != nil {
		return nil, err
	}
	s.sc = sc

	s.inSet <- keyValue{key: Key, value: Value}
	return <-s.outSet, nil
}

func (s *Server) Get(ctx context.Context, Key string) (*storage.GetResponse, error) {
	gc, err := s.storage.Get(context.Background())
	if err != nil {
		return nil, err
	}

	s.gc = gc

	s.inGet <- key{key: Key}
	return <-s.outGet, nil
}

func (s *Server) set() error {
	for s.sc == nil {

	}

	for {
		select {
		case <-s.sc.Context().Done():
			return s.sc.Context().Err()
		case kv := <-s.inSet:
			err := s.sc.Send(&storage.SetRequest{Key: kv.key, Value: kv.value})
			if err != nil {
				return err
			}

			rc, err := s.sc.Recv()
			if err != nil {
				return err
			}

			s.outSet <- rc
		}
	}
}

func (s *Server) get() error {
	for s.gc == nil {

	}

	for {
		select {
		case <-s.gc.Context().Done():
			return s.gc.Context().Err()
		case key := <-s.inGet:
			err := s.gc.Send(&storage.GetRequest{Key: key.key})
			if err != nil {
				return err
			}

			rc, err := s.gc.Recv()
			if err != nil {
				return err
			}

			s.outGet <- rc
		}
	}
}
