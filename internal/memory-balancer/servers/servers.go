package servers

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-storage/pkg/api/storage"
	"log"
	"sync"
)

type Servers struct {
	servers map[int32]*Server
	sync.RWMutex
}

var ErrNotFound = errors.New("the value was not found")

func New() *Servers {
	s := make(map[int32]*Server, 10)
	return &Servers{
		servers: s,
	}
}

func (s *Servers) GetClient() (*Server, bool) {
	s.RLock()
	defer s.RUnlock()

	max := 0.0
	var serverNumber int32 = 0

	for num, cl := range s.servers {
		if float64(cl.Available)/float64(cl.Total)*100 > max {
			serverNumber = num
		}
	}

	cl, ok := s.servers[serverNumber]
	return cl, ok
}

func (s *Servers) Len() int32 {
	s.RLock()
	defer s.RUnlock()
	return int32(len(s.servers))
}

func (s *Servers) AddClient(address string, available, total uint64) (int32, error) {
	s.Lock()
	defer s.Unlock()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return 0, err
	}

	cl := storage.NewStorageClient(conn)
	var stClient = &Server{
		storage:   cl,
		Available: available,
		Total:     total,
	}

	stClient.Number = int32(len(s.servers) + 1)
	s.servers[stClient.Number] = stClient

	return stClient.Number, nil
}

func (s *Servers) Disconnect(number int32) {
	s.Lock()
	defer s.Unlock()
	delete(s.servers, number)
	//s.servers[number] = nil
}

func (s *Servers) GetServers() []string {
	s.RLock()
	defer s.RUnlock()

	var servers = make([]string, 0, 5)
	for num, cl := range s.servers {
		servers = append(servers, fmt.Sprintf("s#%d Avaliable: %d MB, Total: %d MB", num, cl.Available, cl.Total))
	}

	return servers
}

func (s *Servers) DeepSearch(ctx context.Context, key string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	var out = make(chan string, 1)
	defer close(out)

	var wg sync.WaitGroup
	wg.Add(len(s.servers))

	// TODO: Add pull of goroutines
	for _, cl := range s.servers {
		go cl.find(ctxCancel, cancel, &wg, key, out)
	}

	wg.Wait()
	select {
	case v := <-out:
		return v, nil
	default:
		return "", ErrNotFound
	}
}

func (s *Server) find(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, key string, out chan<- string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("the panic in find has been restored, probably sending via a closed channel")
		}
	}()
	get, err := s.Get(ctx, key)
	wg.Done()
	if err != nil {
		return
	}

	cancel()
	out <- get.Value
}

func (s *Servers) GetClientByID(number int32) (*Server, bool) {
	s.RLock()
	defer s.RUnlock()
	srv, ok := s.servers[number]
	return srv, ok
}

func (s *Servers) Exists(number int32) bool {
	_, ok := s.servers[number]
	return ok
}
