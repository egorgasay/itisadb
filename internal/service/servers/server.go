package servers

import (
	"context"
	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk"
	"itisadb/internal/models"
	"sync"
	"sync/atomic"
)

// =============== server ====================== //

func NewServer(cl *itisadb.Client, number int32) *RemoteServer {
	return &RemoteServer{
		sdk:    cl,
		number: number,
		mu:     &sync.RWMutex{},
		tries:  atomic.Uint32{},
	}
}

type RemoteServer struct {
	tries  atomic.Uint32
	ram    models.RAM
	number int32
	mu     *sync.RWMutex

	sdk *itisadb.Client
}

func (s *RemoteServer) Number() int32 {
	return s.number
}

func (s *RemoteServer) Tries() uint32 {
	return s.tries.Load()
}

func (s *RemoteServer) GetOne(ctx context.Context, key string, opts ...itisadb.GetOptions) (res gost.Result[string]) {
	return s.sdk.GetOne(ctx, key, opts...)
}

func (s *RemoteServer) DelOne(ctx context.Context, key string, opts ...itisadb.DeleteOptions) gost.Result[gost.Nothing] {
	return s.sdk.DelOne(ctx, key, opts...)
}

func (s *RemoteServer) SetOne(ctx context.Context, key string, val string, opts ...itisadb.SetOptions) (res gost.Result[int32]) {
	return s.sdk.SetOne(ctx, key, val, opts...)
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

//func (s *RemoteServer) Set(ctx context.Context, key, value string, opts models.SetOptions) error {
//	_, err := s.client.Set(ctx, &api.SetRequest{
//		Key:   key,
//		Value: value,
//	})
//	return err
//}
//
//func (s *RemoteServer) Get(ctx context.Context, key string, opts models.GetOptions) (*api.GetResponse, error) {
//	r, err := s.client.Get(ctx, &api.GetRequest{Key: key})
//	return r, err
//
//}
//
//func (s *RemoteServer) ObjectToJSON(ctx context.Context, name string, opts models.ObjectToJSONOptions) (*api.ObjectToJSONResponse, error) {
//	r, err := s.client.ObjectToJSON(ctx, &api.ObjectToJSONRequest{
//		Name: name,
//	})
//	return r, err
//}
//
//func (s *RemoteServer) GetFromObject(ctx context.Context, name, key string, opts models.GetFromObjectOptions) (*api.GetFromObjectResponse, error) {
//	r, err := s.client.GetFromObject(ctx, &api.GetFromObjectRequest{
//		Key:    key,
//		Object: name,
//	})
//	return r, err
//
//}
//
//func (s *RemoteServer) SetToObject(ctx context.Context, name, key, value string, opts models.SetToObjectOptions) error {
//	_, err := s.client.SetToObject(ctx, &api.SetToObjectRequest{
//		Key:    key,
//		Value:  value,
//		Object: name,
//	})
//	return err
//
//}
//
//func (s *RemoteServer) NewObject(ctx context.Context, name string, opts models.ObjectOptions) error {
//	_, err := s.client.NewObject(ctx, &api.NewObjectRequest{
//		Name: name,
//	})
//	return err
//}
//
//func (s *RemoteServer) Size(ctx context.Context, name string, opts models.SizeOptions) (*api.ObjectSizeResponse, error) {
//	r, err := s.client.Size(ctx, &api.ObjectSizeRequest{
//		Name: name,
//	})
//	return r, err
//}
//
//func (s *RemoteServer) DeleteObject(ctx context.Context, name string, opts models.DeleteObjectOptions) error {
//	_, err := s.client.DeleteObject(ctx, &api.DeleteObjectRequest{
//		Object: name,
//	})
//	return err
//}
//
//func (s *RemoteServer) Delete(ctx context.Context, Key string, opts models.DeleteOptions) error {
//	_, err := s.client.Delete(ctx, &api.DeleteRequest{
//		Key: Key,
//	})
//	return err
//}
//
//func (s *RemoteServer) AttachToObject(ctx context.Context, dst string, src string, opts models.AttachToObjectOptions) error {
//	_, err := s.client.AttachToObject(ctx, &api.AttachToObjectRequest{
//		Dst: dst,
//		Src: src,
//	})
//	return err
//}
//
//func (s *RemoteServer) GetNumber() int32 {
//	return s.number
//}
//
//func (s *RemoteServer) GetTries() uint32 {
//	return s.tries.Load()
//}

func (s *RemoteServer) IncTries() {
	s.tries.Add(1)
}

func (s *RemoteServer) ResetTries() {
	s.tries.Store(0)
}

//func (s *RemoteServer) DeleteAttr(ctx context.Context, attr string, object string, opts models.DeleteAttrOptions) error {
//	_, err := s.client.DeleteAttr(ctx, &api.DeleteAttrRequest{
//		Object: object,
//		Key:    attr,
//	})
//	return err
//}

func (s *RemoteServer) Find(ctx context.Context, key string, out chan<- string, once *sync.Once, _ models.GetOptions) {
	r := s.sdk.GetOne(ctx, key)
	if err := r.Error(); err != nil {
		return
	}

	once.Do(func() {
		out <- r.Unwrap()
	})
}

func (s *RemoteServer) NewObject(ctx context.Context, name string, opts models.ObjectOptions) interface{} {
	// TODO: implement me
	panic("TODO: implement me")
	return nil
}

func (s *RemoteServer) GetFromObject(ctx context.Context, object string, key string, opts models.GetFromObjectOptions) (interface{}, interface{}) {
	// TODO: implement me
	panic("TODO: implement me")
	return nil, nil
}

func (s *RemoteServer) SetToObject(ctx context.Context, object string, key string, value string, opts models.SetToObjectOptions) gost.Result[gost.Nothing] {
	// return s.sdk.Object(ctx, object).Unwrap().Set(ctx, key, value, opts.ToSDK())
	// TODO: implement me
	panic("TODO: implement me")
	return gost.Result[gost.Nothing]{}
}
