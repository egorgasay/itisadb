package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"itisadb/internal/memory-balancer/servers"
	"itisadb/internal/memory-balancer/storage"
	serversmock "itisadb/internal/memory-balancer/usecase/mocks/servers"
	repomock "itisadb/internal/memory-balancer/usecase/mocks/storage"
	gstorage "itisadb/pkg/api/storage"
	storagemock "itisadb/pkg/api/storage/gomocks"
	"itisadb/pkg/logger"
	"sync"
	"testing"
)

func TestUseCase_Connect(t *testing.T) {
	srv := struct {
		serv *servers.Server
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
		storageBehavior     storageBehavior
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
			srv.serv = servers.NewServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &UseCase{
				servers: s,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu: sync.RWMutex{},
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
		serv *servers.Server
	}{}
	type args struct {
		ctx context.Context
		key string
		num int32
	}
	tests := []struct {
		name                string
		storageBehavior     storageBehavior
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
					servers.ErrNotFound)
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

			srv.serv = servers.NewServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &UseCase{
				servers: s,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu: sync.RWMutex{},
			}

			if err = uc.Delete(tt.args.ctx, tt.args.key, tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Disconnect(t *testing.T) {
	srv := struct {
		serv *servers.Server
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := storagemock.NewMockStorageClient(c)
			srv.serv = servers.NewServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &UseCase{
				servers: s,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu: sync.RWMutex{},
			}

			uc.Disconnect(tt.args.number)
		})
	}
}

func TestUseCase_FindInDB(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name                string
		storageBehavior     storageBehavior
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		args                args
		want                string
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				key: "test_key",
			},
			storageBehavior: func(cl *repomock.MockIStorage) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					"test_value", nil)
			},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			want:                "test_value",
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				key: "test_key",
			},
			storageBehavior: func(cl *repomock.MockIStorage) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					"", storage.ErrNotFound)
			},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
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
			srv.serv = servers.NewServer(sc, 1)
			tt.grpcStorageBehavior(sc)
			rm := repomock.NewMockIStorage(c)
			tt.storageBehavior(rm)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &UseCase{
				servers: s,
				storage: rm,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu: sync.RWMutex{},
			}

			got, err := uc.FindInDB(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindInDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindInDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Get(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx          context.Context
		key          string
		serverNumber int32
	}
	tests := []struct {
		name                string
		args                args
		storageBehavior     storageBehavior
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
			storageBehavior: func(cl *repomock.MockIStorage) {
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
			storageBehavior: func(cl *repomock.MockIStorage) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					"", storage.ErrNotFound)
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
			name: "findInPhysDB",
			args: args{
				ctx:          context.Background(),
				key:          "test_key",
				serverNumber: 4,
			},
			storageBehavior: func(cl *repomock.MockIStorage) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					"test_value", nil)
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().Len().Return(int32(1))
				cl.EXPECT().Exists(gomock.Any()).Return(false)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			want: "test_value",
		},
		{
			name: "notFoundInGrpc",
			args: args{
				ctx:          context.Background(),
				key:          "test_key",
				serverNumber: 4,
			},
			storageBehavior: func(cl *repomock.MockIStorage) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					"test_value", nil)
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().Len().Return(int32(1))
				cl.EXPECT().Exists(gomock.Any()).Return(true)
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&gstorage.GetResponse{}, storage.ErrNotFound)
			},
			want: "test_value",
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
			rm := repomock.NewMockIStorage(c)
			s := serversmock.NewMockIServers(c)
			srv.serv = servers.NewServer(sc, 1)

			tt.storageBehavior(rm)
			tt.grpcStorageBehavior(sc)
			tt.serversBehavior(s)

			uc := &UseCase{
				servers: s,
				storage: rm,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu: sync.RWMutex{},
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
		serv *servers.Server
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
		storageBehavior     storageBehavior
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
			storageBehavior: func(cl *repomock.MockIStorage) {},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(
					srv.serv, true)
				cl.EXPECT().Len().Return(int32(1))
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
					&gstorage.SetResponse{}, nil)
			},
			want: 1,
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
			storageBehavior: func(cl *repomock.MockIStorage) {
				cl.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil)
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().Len().Return(int32(0))
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			want: -1,
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
			rm := repomock.NewMockIStorage(c)
			s := serversmock.NewMockIServers(c)
			srv.serv = servers.NewServer(sc, 1)

			tt.storageBehavior(rm)
			tt.grpcStorageBehavior(sc)
			tt.serversBehavior(s)
			uc := &UseCase{
				servers: s,
				storage: rm,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu: sync.RWMutex{},
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
