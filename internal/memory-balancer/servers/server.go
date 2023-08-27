package servers

import (
	"context"
	"errors"
	"fmt"
	"itisadb/pkg/api/storage"
	"reflect"
	"sync"
	"sync/atomic"
)

// =============== server ====================== //

func NewServer(storage storage.StorageClient, number int32) *Server {
	return &Server{
		storage: storage,
		number:  0,
		mu:      &sync.RWMutex{},
		tries:   atomic.Uint32{},
	}
}

var ErrAlreadyExists = errors.New("already exists")
var ErrUnavailable = errors.New("server is unavailable")

type Server struct {
	tries   atomic.Uint32
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
	return err
}

func (s *Server) Get(ctx context.Context, Key string) (*storage.GetResponse, error) {
	r, err := s.storage.Get(ctx, &storage.GetRequest{Key: Key})
	s.setRAM(r)
	return r, err

}

func (s *Server) ObjectToJSON(ctx context.Context, name string) (*storage.ObjectToJSONResponse, error) {
	r, err := s.storage.ObjectToJSON(ctx, &storage.ObjectToJSONRequest{
		Name: name,
	})
	s.setRAM(r)
	return r, err
}

func (s *Server) GetFromObject(ctx context.Context, name, Key string) (*storage.GetResponse, error) {
	r, err := s.storage.GetFromObject(ctx, &storage.GetFromObjectRequest{
		Key:  Key,
		Name: name,
	})
	s.setRAM(r)
	return r, err

}

func (s *Server) SetToObject(ctx context.Context, name, Key, Value string, unique bool) error {
	r, err := s.storage.SetToObject(ctx, &storage.SetToObjectRequest{
		Key:    Key,
		Value:  Value,
		Name:   name,
		Unique: unique,
	})
	s.setRAM(r)
	return err

}

func (s *Server) NewObject(ctx context.Context, name string) error {
	r, err := s.storage.NewObject(ctx, &storage.NewObjectRequest{
		Name: name,
	})
	s.setRAM(r)
	return err
}

func (s *Server) Size(ctx context.Context, name string) (*storage.ObjectSizeResponse, error) {
	r, err := s.storage.Size(ctx, &storage.ObjectSizeRequest{
		Name: name,
	})
	s.setRAM(r)
	return r, err
}

func (s *Server) DeleteObject(ctx context.Context, name string) error {
	r, err := s.storage.DeleteObject(ctx, &storage.DeleteObjectRequest{
		Object: name,
	})
	s.setRAM(r)
	return err
}

func (s *Server) Delete(ctx context.Context, Key string) error {
	r, err := s.storage.Delete(ctx, &storage.DeleteRequest{
		Key: Key,
	})
	s.setRAM(r)
	return err
}

var ErrCircularAttachment = fmt.Errorf("circular attachment")

func (s *Server) AttachToObject(ctx context.Context, dst string, src string) error {
	r, err := s.storage.AttachToObject(ctx, &storage.AttachToObjectRequest{
		Dst: dst,
		Src: src,
	})
	s.setRAM(r)
	return err
}

func (s *Server) GetNumber() int32 {
	return s.number
}

func (s *Server) GetTries() uint32 {
	return s.tries.Load()
}

func (s *Server) IncTries() {
	s.tries.Add(1)
}

func (s *Server) ResetTries() {
	s.tries.Store(0)
}

func (s *Server) DeleteAttr(ctx context.Context, attr string, object string) error {
	r, err := s.storage.DeleteAttr(ctx, &storage.DeleteAttrRequest{
		Name: object,
		Key:  attr,
	})
	s.setRAM(r)
	return err
}
