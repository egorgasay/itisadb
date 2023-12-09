package servers

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk"
	api "github.com/egorgasay/itisadb-shared-proto/go"
	"itisadb/internal/models"
	"sync"
	"sync/atomic"
)

// =============== server ====================== //

func NewServer(cl api.ItisaDBClient, number int32) *RemoteServer {
	return &RemoteServer{
		client: cl,
		number: number,
		mu:     &sync.RWMutex{},
		tries:  atomic.Uint32{},
	}
}

type RemoteServer struct {
	tries  atomic.Uint32
	client api.ItisaDBClient
	ram    models.RAM
	number int32
	mu     *sync.RWMutex

	sdk *itisadb.Client
}

func (s *RemoteServer) RAM() models.RAM {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ram
}

func (s *RemoteServer) SetRAM(ram models.RAM) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ram.Total == 0 {
		return
	}
	s.ram = models.RAM{Total: ram.Total, Available: ram.Available}
}

func (s *RemoteServer) Set(ctx context.Context, key, value string, opts models.SetOptions) error {
	_, err := s.client.Set(ctx, &api.SetRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (s *RemoteServer) Get(ctx context.Context, key string, opts models.GetOptions) (*api.GetResponse, error) {
	r, err := s.client.Get(ctx, &api.GetRequest{Key: key})
	return r, err

}

func (s *RemoteServer) ObjectToJSON(ctx context.Context, name string, opts models.ObjectToJSONOptions) (*api.ObjectToJSONResponse, error) {
	r, err := s.client.ObjectToJSON(ctx, &api.ObjectToJSONRequest{
		Name: name,
	})
	return r, err
}

func (s *RemoteServer) GetFromObject(ctx context.Context, name, key string, opts models.GetFromObjectOptions) (*api.GetFromObjectResponse, error) {
	r, err := s.client.GetFromObject(ctx, &api.GetFromObjectRequest{
		Key:    key,
		Object: name,
	})
	return r, err

}

func (s *RemoteServer) SetToObject(ctx context.Context, name, key, value string, opts models.SetToObjectOptions) error {
	_, err := s.client.SetToObject(ctx, &api.SetToObjectRequest{
		Key:    key,
		Value:  value,
		Object: name,
	})
	return err

}

func (s *RemoteServer) NewObject(ctx context.Context, name string, opts models.ObjectOptions) error {
	_, err := s.client.NewObject(ctx, &api.NewObjectRequest{
		Name: name,
	})
	return err
}

func (s *RemoteServer) Size(ctx context.Context, name string, opts models.SizeOptions) (*api.ObjectSizeResponse, error) {
	r, err := s.client.Size(ctx, &api.ObjectSizeRequest{
		Name: name,
	})
	return r, err
}

func (s *RemoteServer) DeleteObject(ctx context.Context, name string, opts models.DeleteObjectOptions) error {
	_, err := s.client.DeleteObject(ctx, &api.DeleteObjectRequest{
		Object: name,
	})
	return err
}

func (s *RemoteServer) Delete(ctx context.Context, Key string, opts models.DeleteOptions) error {
	_, err := s.client.Delete(ctx, &api.DeleteRequest{
		Key: Key,
	})
	return err
}

func (s *RemoteServer) AttachToObject(ctx context.Context, dst string, src string, opts models.AttachToObjectOptions) error {
	_, err := s.client.AttachToObject(ctx, &api.AttachToObjectRequest{
		Dst: dst,
		Src: src,
	})
	return err
}

func (s *RemoteServer) GetNumber() int32 {
	return s.number
}

func (s *RemoteServer) GetTries() uint32 {
	return s.tries.Load()
}

func (s *RemoteServer) IncTries() {
	s.tries.Add(1)
}

func (s *RemoteServer) ResetTries() {
	s.tries.Store(0)
}

func (s *RemoteServer) DeleteAttr(ctx context.Context, attr string, object string, opts models.DeleteAttrOptions) error {
	_, err := s.client.DeleteAttr(ctx, &api.DeleteAttrRequest{
		Object: object,
		Key:    attr,
	})
	return err
}

func (s *RemoteServer) Find(ctx context.Context, key string, out chan<- string, once *sync.Once, opts models.GetOptions) {
	get, err := s.Get(ctx, key, opts)
	if err != nil {
		return
	}

	once.Do(func() {
		out <- get.Value
	})
}
