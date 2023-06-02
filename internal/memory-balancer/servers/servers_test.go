package servers

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/pkg/api/storage"
	storagemock "itisadb/pkg/api/storage/gomocks"
	"reflect"
	"sync"
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
					Value: "value"}, nil).AnyTimes()
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
					nil, status.Error(codes.NotFound, "not found")).AnyTimes()
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

			s := &Server{
				tries:   0,
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
					t.Errorf("Server.find() = %v, want %v", got, tt.want)
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
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Servers{
				servers: map[int32]*Server{
					1: {
						tries:   0,
						storage: cl,
						ram: RAM{
							available: 100,
							total:     100,
						},
						number: 1,
						mu:     &sync.RWMutex{},
					},
				},
				freeID:  2,
				RWMutex: sync.RWMutex{},
			}

			got, err := s.AddClient(tt.args.address, tt.args.available, tt.args.total, tt.args.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_DeepSearch(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Servers{
				servers: map[int32]*Server{
					1: {
						tries: 0,
						//storage: cl,
						ram: RAM{
							available: 100,
							total:     100,
						},
						number: 1,
						mu:     &sync.RWMutex{},
					},
				},
				freeID:  2,
				RWMutex: sync.RWMutex{},
			}
			got, err := s.DeepSearch(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeepSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeepSearch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_Disconnect(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	type args struct {
		number int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Servers{
				servers: tt.fields.servers,
				freeID:  tt.fields.freeID,
				RWMutex: tt.fields.RWMutex,
			}
			s.Disconnect(tt.args.number)
		})
	}
}

func TestServers_Exists(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	type args struct {
		number int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Servers{
				servers: tt.fields.servers,
				freeID:  tt.fields.freeID,
				RWMutex: tt.fields.RWMutex,
			}
			if got := s.Exists(tt.args.number); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_GetClient(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   *Server
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Servers{
				servers: tt.fields.servers,
				freeID:  tt.fields.freeID,
				RWMutex: tt.fields.RWMutex,
			}
			got, got1 := s.GetClient()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestServers_GetClientByID(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	type args struct {
		number int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Server
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Servers{
				servers: tt.fields.servers,
				freeID:  tt.fields.freeID,
				RWMutex: tt.fields.RWMutex,
			}
			got, got1 := s.GetClientByID(tt.args.number)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientByID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetClientByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestServers_GetServers(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Servers{
				servers: tt.fields.servers,
				freeID:  tt.fields.freeID,
				RWMutex: tt.fields.RWMutex,
			}
			if got := s.GetServers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_Len(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Servers{
				servers: tt.fields.servers,
				freeID:  tt.fields.freeID,
				RWMutex: tt.fields.RWMutex,
			}
			if got := s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_SetToAll(t *testing.T) {
	type fields struct {
		servers map[int32]*Server
		freeID  int32
		RWMutex sync.RWMutex
	}
	type args struct {
		ctx     context.Context
		key     string
		val     string
		uniques bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Servers{
				servers: tt.fields.servers,
				freeID:  tt.fields.freeID,
				RWMutex: tt.fields.RWMutex,
			}
			if got := s.SetToAll(tt.args.ctx, tt.args.key, tt.args.val, tt.args.uniques); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetToAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
