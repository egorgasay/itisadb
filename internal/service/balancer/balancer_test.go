package balancer

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/pkg"
	"itisadb/pkg/api/storage"
	storagemock "itisadb/pkg/api/storage/gomocks"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestServer_find(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		out  chan string
		once *sync.Once
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		want         string
		wantErr      bool
		args         args
	}{
		{
			name: "ok",
			args: args{
				ctx:  context.Background(),
				key:  "test",
				out:  make(chan string),
				once: &sync.Once{},
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&storage.GetResponse{
					Value: "value"}, nil)
			},
			want: "value",
		},
		{
			name: "notFound",
			args: args{
				ctx:  context.Background(),
				key:  "test34",
				out:  make(chan string),
				once: &sync.Once{},
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			ctx, cancel := context.WithTimeout(tt.args.ctx, 2*time.Second)
			defer cancel()

			s := &RemoteServer{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}

			go s.find(ctx, tt.args.key, tt.args.out, tt.args.once)

			select {
			case <-ctx.Done():
				if tt.wantErr {
					return
				}
				t.Fatalf("timeout")
			case got := <-tt.args.out:
				if got != tt.want {
					t.Errorf("RemoteServer.find() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServers_AddClient(t *testing.T) {
	type args struct {
		address   string
		available uint64
		total     uint64
		server    int32
	}
	tests := []struct {
		name         string
		args         args
		want         int32
		mockBehavior func(cl *storagemock.MockStorageClient)
		wantErr      bool
	}{
		{
			name: "ok",
			args: args{
				address:   "test",
				available: 100,
				total:     100,
				server:    1,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			s := &Balancer{
				servers: map[int32]*RemoteServer{},
				freeID:  1,
				RWMutex: sync.RWMutex{},
			}

			got, err := s.AddServer(tt.args.address, tt.args.available, tt.args.total, tt.args.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddServer() error = %v, wantCode %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddServer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_DeepSearch(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior func(cl *storagemock.MockStorageClient)
		servers      map[int32]*RemoteServer
		want         string
		wantErr      bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&storage.GetResponse{
					Value: "value"}, nil).AnyTimes()
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 44,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 43,
						total:     100,
					},
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			want: "value",
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found")).AnyTimes()
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 23,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
			},
			wantErr: true,
		},
		{
			name: "badConnection",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection")).AnyTimes()
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 23,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
			},
			wantErr: true,
		},
	}
	pool := make(chan struct{}, runtime.GOMAXPROCS(0))
	defer close(pool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			for _, serv := range tt.servers {
				serv.storage = cl
			}

			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
				poolCh:  pool,
			}

			got, err := s.DeepSearch(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeepSearch() error = %v, wantCode %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeepSearch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_Disconnect(t *testing.T) {
	type args struct {
		number int32
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		servers      map[int32]*RemoteServer
		want         bool
		args         args
	}{
		{
			name: "ok",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
			},
			args: args{
				number: 1,
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 32,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
			},
			want: true,
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
			},
			args: args{
				number: 1,
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 34,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
			},
			want: false,
		},
		{
			name: "notFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
			},
			args: args{
				number: 333,
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 13,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)
			for _, s := range tt.servers {
				s.storage = cl
			}

			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
			}
			s.Disconnect(tt.args.number)

			if tt.want && s.Exists(tt.args.number) {
				t.Errorf("Disconnect() error = %v, wantCode %v", s.Exists(tt.args.number), false)
			}
		})
	}
}

func TestServers_Exists(t *testing.T) {
	type args struct {
		number int32
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		servers      map[int32]*RemoteServer
		args         args
		want         bool
	}{
		{
			name: "ok",
			args: args{
				number: 1,
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 55,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
			},
			want: true,
		},
		{
			name: "notFound",
			args: args{
				number: 3,
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 11,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
			}
			if got := s.Exists(tt.args.number); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_GetClient(t *testing.T) {
	tests := []struct {
		name    string
		servers map[int32]*RemoteServer
		want    int32
		wantRes bool
	}{
		{
			name: "ok",
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 100,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 66,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 77,
						total:     100,
					},
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			want:    1,
			wantRes: true,
		},
		{
			name: "ok2",
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 33,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 66,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 77,
						total:     100,
					},
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			want:    3,
			wantRes: true,
		},
		{
			name:    "noClients",
			servers: map[int32]*RemoteServer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
			}
			got, ok := s.GetServer()
			if ok != tt.wantRes {
				t.Errorf("GetServer() got1 = %v, want %v", ok, tt.wantRes)
				return
			}

			if !ok {
				return
			}

			if got.number != tt.want {
				t.Errorf("GetServer() got = %v, want %v", got.number, tt.want)
			}
		})
	}
}

func TestServers_GetClientByID(t *testing.T) {
	type args struct {
		number int32
	}
	tests := []struct {
		name    string
		servers map[int32]*RemoteServer
		args    args
		want    int32
		ok      bool
	}{
		{
			name: "ok",
			servers: map[int32]*RemoteServer{
				1: {
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			args: args{
				number: 1,
			},
			want: 1,
			ok:   true,
		},
		{
			name: "notFound",
			servers: map[int32]*RemoteServer{
				1: {
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			args: args{
				number: 33,
			},
			ok: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
			}
			got, ok := s.GetServerByID(tt.args.number)
			if ok != tt.ok {
				t.Errorf("GetServerByID() got1 = %v, want %v", ok, tt.ok)
				return
			} else if !ok {
				return
			}

			if got.number != tt.want {
				t.Errorf("GetServerByID() got = %v, want %v", got.number, tt.want)
			}

		})
	}
}

func TestServers_GetServers(t *testing.T) {
	tests := []struct {
		name    string
		servers map[int32]*RemoteServer
		want    []string
	}{
		{
			name: "ok",
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 33,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 66,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 77,
						total:     100,
					},
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			want: []string{
				"s#1 Avaliable: 33 MB, Total: 100 MB",
				"s#2 Avaliable: 66 MB, Total: 100 MB",
				"s#3 Avaliable: 77 MB, Total: 100 MB",
			},
		},
		{
			name:    "noServers",
			servers: map[int32]*RemoteServer{},
			want:    []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
			}
			if got := s.GetServers(); !pkg.IsTheSameArray(got, tt.want) {
				t.Errorf("GetServers() = \n%v,\nwant \n%v", got, tt.want)
			}
		})
	}
}

func TestServers_Len(t *testing.T) {
	tests := []struct {
		name    string
		servers map[int32]*RemoteServer
		want    int32
	}{
		{
			name: "ok",
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 33,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 66,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 77,
						total:     100,
					},
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			want: 3,
		},
		{
			name:    "noServers",
			servers: map[int32]*RemoteServer{},
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
			}
			if got := s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_SetToAll(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		val     string
		uniques bool
	}
	tests := []struct {
		name         string
		mockBehavior func(r *storagemock.MockStorageClient)
		args         args
		servers      map[int32]*RemoteServer
		want         []int32
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				key: "key",
				val: "val",
			},
			mockBehavior: func(r *storagemock.MockStorageClient) {
				r.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil, nil).Times(3)
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 33,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 66,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 77,
						total:     100,
					},
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			want: []int32{},
		},
		{
			name: "badConnection",
			args: args{
				ctx: context.Background(),
				key: "key",
				val: "val",
			},
			mockBehavior: func(r *storagemock.MockStorageClient) {
				r.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection")).Times(3)
			},
			servers: map[int32]*RemoteServer{
				1: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 33,
						total:     100,
					},
					number: 1,
					mu:     &sync.RWMutex{},
				},
				2: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 66,
						total:     100,
					},
					number: 2,
					mu:     &sync.RWMutex{},
				},
				3: {
					tries: atomic.Uint32{},
					ram: RAM{
						available: 77,
						total:     100,
					},
					number: 3,
					mu:     &sync.RWMutex{},
				},
			},
			want: []int32{1, 2, 3},
		},
	}

	pool := make(chan struct{}, runtime.GOMAXPROCS(0))
	defer close(pool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)
			for _, s := range tt.servers {
				s.storage = cl
			}

			s := &Balancer{
				servers: tt.servers,
				freeID:  int32(len(tt.servers) + 1),
				RWMutex: sync.RWMutex{},
				poolCh:  pool,
			}

			got := s.SetToAll(tt.args.ctx, tt.args.key, tt.args.val, tt.args.uniques)
			if !pkg.IsTheSameArray(got, tt.want) {
				t.Errorf("SetToAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
