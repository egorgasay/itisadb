package servers

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/pkg/api/storage"
	"reflect"
	"sync"
)

// =============== server ====================== //

var ErrAlreadyExists = errors.New("already exists")
var ErrUnavailable = errors.New("server is unavailable")

type Server struct {
	tries   uint
	storage storage.StorageClient
	ram     RAM
	number  int32
	mu      *sync.RWMutex
}

type RAM struct {
	available uint64
	total     uint64
}

func (s *Server) GetRAM() RAM {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ram
}

type getRAM interface {
	GetRam() *storage.Ram
}

func (s *Server) setRAM(r getRAM) {
	if reflect.ValueOf(r).IsNil() {
		return
	}

	ram := r.GetRam()
	s.mu.Lock()
	defer s.mu.Unlock()
	if ram == nil || ram.Total == 0 {
		return
	}
	s.ram = RAM{total: ram.Total, available: ram.Available}
}

func (s *Server) Set(ctx context.Context, Key, Value string, unique bool) error {
	r, err := s.storage.Set(ctx, &storage.SetRequest{Key: Key, Value: Value, Unique: unique})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)

		if !ok {
			return err
		}

		if ok && st.Code() == codes.AlreadyExists {
			return ErrAlreadyExists
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}

		return err
	}

	return err
}

func (s *Server) Get(ctx context.Context, Key string) (*storage.GetResponse, error) {
	r, err := s.storage.Get(ctx, &storage.GetRequest{Key: Key})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		if st.Code() == codes.NotFound {
			return nil, ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return nil, ErrUnavailable
		}
		return nil, err
	}

	return r, err

}

func (s *Server) GetIndex(ctx context.Context, name string) (*storage.GetIndexResponse, error) {
	r, err := s.storage.GetIndex(ctx, &storage.GetIndexRequest{
		Name: name,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		if st.Code() == codes.NotFound {
			return nil, ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return nil, ErrUnavailable
		}
		return nil, err
	}

	return r, err

}

func (s *Server) GetFromIndex(ctx context.Context, name, Key string) (*storage.GetResponse, error) {
	r, err := s.storage.GetFromIndex(ctx, &storage.GetFromIndexRequest{
		Key:  Key,
		Name: name,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		if st.Code() == codes.NotFound {
			return nil, ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return nil, ErrUnavailable
		}
		return nil, err
	}

	return r, err

}

func (s *Server) SetToIndex(ctx context.Context, name, Key, Value string, unique bool) error {
	r, err := s.storage.SetToIndex(ctx, &storage.SetToIndexRequest{
		Key:    Key,
		Value:  Value,
		Name:   name,
		Unique: unique,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.AlreadyExists {
			return ErrAlreadyExists
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}

	return err

}

func (s *Server) NewIndex(ctx context.Context, name string) error {
	r, err := s.storage.NewIndex(ctx, &storage.NewIndexRequest{
		Name: name,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}
	return nil
}

func (s *Server) Size(ctx context.Context, name string) (*storage.IndexSizeResponse, error) {
	r, err := s.storage.Size(ctx, &storage.IndexSizeRequest{
		Name: name,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		if st.Code() == codes.NotFound {
			return nil, ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return nil, ErrUnavailable
		}
		return nil, err
	}
	return r, nil
}

func (s *Server) DeleteIndex(ctx context.Context, name string) error {
	r, err := s.storage.DeleteIndex(ctx, &storage.DeleteIndexRequest{
		Index: name,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.NotFound {
			return ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}
	return nil
}

func (s *Server) Delete(ctx context.Context, Key string) error {
	r, err := s.storage.Delete(ctx, &storage.DeleteRequest{
		Key: Key,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.NotFound {
			return ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}
	return nil
}

func (s *Server) AttachToIndex(ctx context.Context, dst string, src string) error {
	r, err := s.storage.AttachToIndex(ctx, &storage.AttachToIndexRequest{
		Dst: dst,
		Src: src,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.NotFound {
			return ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}

		return err
	}
	return nil
}

func (s *Server) GetNumber() int32 {
	return s.number
}

func (s *Server) GetTries() uint {
	return s.tries
}

func (s *Server) IncTries() {
	s.mu.Lock()
	s.tries++ // TODO: atomic??
	s.mu.Unlock()
}

func (s *Server) ResetTries() {
	s.mu.Lock()
	s.tries = 0 // TODO: atomic??
	s.mu.Unlock()
}

func (s *Server) DeleteAttr(ctx context.Context, attr string, index string) error {
	r, err := s.storage.DeleteAttr(ctx, &storage.DeleteAttrRequest{
		Name: index,
		Key:  attr,
	})
	s.setRAM(r)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.NotFound {
			return ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}
	return nil
}
