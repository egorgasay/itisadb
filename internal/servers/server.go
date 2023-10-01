package servers

import (
	"context"
	"itisadb/internal/models"
	"itisadb/pkg/api"
	client "itisadb/pkg/api"
	"sync"
	"sync/atomic"
)

// =============== server ====================== //

func NewServer(client api.ItisaDBClient, number int32) *Server {
	return &Server{
		client: client,
		number: 0,
		mu:     &sync.RWMutex{},
		tries:  atomic.Uint32{},
	}
}

type Server struct {
	tries  atomic.Uint32
	client client.ItisaDBClient
	ram    models.RAM
	number int32
	mu     *sync.RWMutex
}

func (s *Server) GetRAM() models.RAM {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ram
}

func (s *Server) setRAM(ram models.RAM) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ram.Total == 0 {
		return
	}
	s.ram = models.RAM{Total: ram.Total, Available: ram.Available}
}

func (s *Server) Set(ctx context.Context, Key, Value string, unique bool) error {
	_, err := s.client.Set(ctx, &client.SetRequest{
		Key:     Key,
		Value:   Value,
		Uniques: unique,
	})
	return err
}

func (s *Server) Get(ctx context.Context, Key string) (*client.GetResponse, error) {
	r, err := s.client.Get(ctx, &client.GetRequest{Key: Key})
	return r, err

}

func (s *Server) ObjectToJSON(ctx context.Context, name string) (*client.ObjectToJSONResponse, error) {
	r, err := s.client.ObjectToJSON(ctx, &client.ObjectToJSONRequest{
		Name: name,
	})
	return r, err
}

func (s *Server) GetFromObject(ctx context.Context, name, Key string) (*client.GetFromObjectResponse, error) {
	r, err := s.client.GetFromObject(ctx, &client.GetFromObjectRequest{
		Key:    Key,
		Object: name,
	})
	return r, err

}

func (s *Server) SetToObject(ctx context.Context, name, Key, Value string, unique bool) error {
	_, err := s.client.SetToObject(ctx, &client.SetToObjectRequest{
		Key:     Key,
		Value:   Value,
		Object:  name,
		Uniques: unique,
	})
	return err

}

func (s *Server) NewObject(ctx context.Context, name string) error {
	_, err := s.client.NewObject(ctx, &client.NewObjectRequest{
		Name: name,
	})
	return err
}

func (s *Server) Size(ctx context.Context, name string) (*client.ObjectSizeResponse, error) {
	r, err := s.client.Size(ctx, &client.ObjectSizeRequest{
		Name: name,
	})
	return r, err
}

func (s *Server) DeleteObject(ctx context.Context, name string) error {
	_, err := s.client.DeleteObject(ctx, &client.DeleteObjectRequest{
		Object: name,
	})
	return err
}

func (s *Server) Delete(ctx context.Context, Key string) error {
	_, err := s.client.Delete(ctx, &client.DeleteRequest{
		Key: Key,
	})
	return err
}

func (s *Server) AttachToObject(ctx context.Context, dst string, src string) error {
	_, err := s.client.AttachToObject(ctx, &client.AttachToObjectRequest{
		Dst: dst,
		Src: src,
	})
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
	_, err := s.client.DeleteAttr(ctx, &client.DeleteAttrRequest{
		Object: object,
		Key:    attr,
	})
	return err
}
