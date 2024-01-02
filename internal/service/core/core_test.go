package core

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	servers2 "itisadb/internal/balancer"
	"itisadb/internal/service/balancer"
	serversmock "itisadb/internal/service/core/mocks/servers"
	gstorage "itisadb/pkg/api/storage"
	storagemock "itisadb/pkg/api/storage/gomocks"
	"itisadb/pkg/logger"
	"sync"
	"testing"
)

func TestUseCase_Connect(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}

	type args struct {
		address   string
		available uint64
		total     uint64
		server    int32
	}
	tests := []struct {
		name                string
		args                args
		want                int32
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				address:   "localhost:8080",
				available: 1,
				total:     1,
				server:    1,
			},
			want: 1,
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().AddServer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int32(1), nil)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
		},
		{
			name: "connectionError",
			args: args{
				address:   "localhost:8080",
				available: 1,
				total:     1,
				server:    1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().AddServer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int32(0),
					errors.New("connection error"))
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
	}

	c := gomock.NewController(t)
	defer c.Finish()
	loggerInstance, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("failed to inizialise logger: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := storagemock.NewMockStorageClient(c)
			srv.serv = balancer.NewRemoteServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &Core{
				balancer: s,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test_object": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			got, err := uc.Connect(tt.args.address, tt.args.available, tt.args.total, tt.args.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Connect() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Delete(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}
	type args struct {
		ctx context.Context
		key string
		num int32
	}
	tests := []struct {
		name                string
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		args                args
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				key: "test_key",
				num: 1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "serverNotFound",
			args: args{
				ctx: context.Background(),
				key: "test_key",
				num: 3,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
		{
			name: "valueNotFound",
			args: args{
				ctx: context.Background(),
				key: "test_key",
				num: 1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil,
					servers2.ErrNotFound)
			},
			wantErr: true,
		},
	}

	c := gomock.NewController(t)
	defer c.Finish()
	loggerInstance, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("failed to inizialise logger: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := storagemock.NewMockStorageClient(c)
			tt.grpcStorageBehavior(sc)

			srv.serv = balancer.NewRemoteServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &Core{
				balancer: s,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test_object": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			if err = uc.Delete(tt.args.ctx, tt.args.key, tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Disconnect(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}

	type args struct {
		number int32
	}
	tests := []struct {
		name            string
		serversBehavior serversBehavior
		args            args
	}{
		{
			name: "success",
			args: args{
				number: 1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().Disconnect(gomock.Any()).Return()
			},
		},
	}
	c := gomock.NewController(t)
	defer c.Finish()
	loggerInstance, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("failed to inizialise logger: %v", err)
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := storagemock.NewMockStorageClient(c)
			srv.serv = balancer.NewRemoteServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &Core{
				balancer: s,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test_object": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			err = uc.Disconnect(ctx, tt.args.number)
			if err != nil {
				t.Errorf("Disconnect() error = %v", err)
			}
		})
	}
}

func TestUseCase_Get(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}

	type args struct {
		ctx          context.Context
		key          string
		serverNumber int32
	}
	tests := []struct {
		name                string
		args                args
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		want                string
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx:          context.Background(),
				key:          "test_key",
				serverNumber: 1,
			},

			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(
					srv.serv, true)
				cl.EXPECT().Len().Return(int32(1))
				cl.EXPECT().Exists(gomock.Any()).Return(true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&gstorage.GetResponse{Value: "test_value"}, nil)
			},
			want: "test_value",
		},
		{
			name: "clNotFoundAndPhysNotFound",
			args: args{
				ctx:          context.Background(),
				key:          "test_key",
				serverNumber: 4,
			},

			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().Len().Return(int32(1))
				cl.EXPECT().Exists(gomock.Any()).Return(false)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			wantErr: true,
		},
		{
			name: "notFoundInGrpc",
			args: args{
				ctx:          context.Background(),
				key:          "test_key",
				serverNumber: 1,
			},

			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().Len().Return(int32(1))
				cl.EXPECT().Exists(gomock.Any()).Return(true)
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&gstorage.GetResponse{}, status.Error(codes.NotFound, "not found"))
			},
			wantErr: true,
		},
	}
	c := gomock.NewController(t)
	defer c.Finish()
	loggerInstance, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("failed to inizialise logger: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := storagemock.NewMockStorageClient(c)
			s := serversmock.NewMockIServers(c)
			srv.serv = balancer.NewRemoteServer(sc, 1)

			tt.grpcStorageBehavior(sc)
			tt.serversBehavior(s)

			uc := &Core{
				balancer: s,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test_object": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			got, err := uc.Get(tt.args.ctx, tt.args.key, tt.args.serverNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Set(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}
	type args struct {
		ctx          context.Context
		key          string
		val          string
		serverNumber int32
		uniques      bool
	}
	tests := []struct {
		name                string
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		args                args
		want                int32
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx:          context.Background(),
				key:          "test_key",
				val:          "test_value",
				serverNumber: 1,
				uniques:      false,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(
					srv.serv, true)
				cl.EXPECT().Len().Return(int32(1))
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
					&gstorage.SetResponse{}, nil)
			},
			want: 0,
		},
		{
			name: "noServers",
			args: args{
				ctx:          context.Background(),
				key:          "test_key",
				val:          "test_value",
				serverNumber: 1,
				uniques:      false,
			},

			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().Len().Return(int32(0))
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			wantErr: true,
		},
	}
	c := gomock.NewController(t)
	defer c.Finish()
	loggerInstance, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("failed to inizialise logger: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := storagemock.NewMockStorageClient(c)
			s := serversmock.NewMockIServers(c)
			srv.serv = balancer.NewRemoteServer(sc, 1)

			tt.grpcStorageBehavior(sc)
			tt.serversBehavior(s)
			uc := &Core{
				balancer: s,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test_object": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}
			got, err := uc.Set(tt.args.ctx, tt.args.key, tt.args.val, tt.args.serverNumber, tt.args.uniques)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Set() got = %v, want %v", got, tt.want)
			}
		})
	}
}
