package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/memory-balancer/servers"
	repo "itisadb/internal/memory-balancer/storage"
	serversmock "itisadb/internal/memory-balancer/usecase/mocks/servers"
	"itisadb/pkg/api/storage"
	"itisadb/pkg/api/storage/gomocks"
	"itisadb/pkg/logger"
	"reflect"
	"sync"
	"testing"
)

type serversBehavior func(cl *serversmock.MockIServers)
type gStorageBehavior func(cl *storagemock.MockStorageClient)

func TestUseCase_AttachToIndex(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx context.Context
		dst string
		src string
	}
	tests := []struct {
		name                string
		args                args
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dst: "test1",
				src: "test2",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "clientNotFound",
			args: args{
				ctx: context.Background(),
				dst: "test1",
				src: "test2",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				dst: "test3",
				src: "test2",
			},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
		{
			name: "badConnection",
			args: args{
				ctx: context.Background(),
				dst: "test1",
				src: "test2",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
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
					"test1": 1,
					"test2": 2,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			if err = uc.AttachToIndex(tt.args.ctx, tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("AttachToIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_DeleteAttr(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx   context.Context
		attr  string
		index string
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
				ctx:   context.Background(),
				attr:  "test1",
				index: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx:   context.Background(),
				attr:  "test1",
				index: "test_index2",
			},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
		{
			name: "attrNotFound",
			args: args{
				ctx:   context.Background(),
				attr:  "test1",
				index: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "attr was not found"))
			},
			wantErr: true,
		},
		{
			name: "servNotFound",
			args: args{
				ctx:   context.Background(),
				attr:  "test1",
				index: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
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
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			if err = uc.DeleteAttr(tt.args.ctx, tt.args.attr, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAttr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_DeleteIndex(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx  context.Context
		name string
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
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_index33",
			},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
		{
			name: "badConnection",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			wantErr: true,
		},
		{
			name: "servNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
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
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}
			if err := uc.DeleteIndex(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_GetFromIndex(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx          context.Context
		index        string
		key          string
		serverNumber int32
	}
	tests := []struct {
		name                string
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		args                args
		want                string
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx:          context.Background(),
				index:        "test_index",
				key:          "test1",
				serverNumber: 1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).Return(
					&storage.GetResponse{Value: "test1"}, nil)
			},
			want:    "test1",
			wantErr: false,
		},
		{
			name: "indexNotFound(fromServer)",
			args: args{
				ctx:          context.Background(),
				index:        "test_index",
				key:          "test1",
				serverNumber: 1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).Return(
					&storage.GetResponse{}, status.Error(codes.NotFound, "index not found"))
			},
			wantErr: true,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx:          context.Background(),
				index:        "test_ind33ex33",
				key:          "test1",
				serverNumber: 0,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			wantErr: true,
		},
		{
			name: "serverNotFound",
			args: args{
				ctx:          context.Background(),
				index:        "test_index",
				key:          "test1",
				serverNumber: 3,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
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
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			got, err := uc.GetFromIndex(tt.args.ctx, tt.args.index, tt.args.key, tt.args.serverNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetIndex(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name                string
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		args                args
		want                map[string]string
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetIndex(gomock.Any(), gomock.Any()).Return(
					&storage.GetIndexResponse{
						Index: map[string]string{
							"test1": "test1",
							"test2": "test2",
						},
					}, nil)
			},
			want: map[string]string{
				"test1": "test1",
				"test2": "test2",
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_index33",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			wantErr: true,
		},
		{
			name: "indexNotFound(remoteStorage)",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetIndex(gomock.Any(), gomock.Any()).Return(
					&storage.GetIndexResponse{}, status.Error(codes.NotFound, "index not found"))
			},
			wantErr: true,
		},
		{
			name: "serverNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
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
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			got, err := uc.GetIndex(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Index(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx  context.Context
		name string
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
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().NewIndex(gomock.Any(), gomock.Any()).Return(
					&storage.NewIndexResponse{}, nil)
			},
			want: 0,
		},
		{
			name: "create",
			args: args{
				ctx:  context.Background(),
				name: "test_index2",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServer().Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().NewIndex(gomock.Any(), gomock.Any()).Return(
					&storage.NewIndexResponse{}, nil)
			},
			want: 0,
		}, {
			name: "noActiveClients",
			args: args{
				ctx:  context.Background(),
				name: "test_index2",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServer().Return(nil, false)
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

	st, err := repo.New()
	if err != nil {
		t.Fatalf("failed to inizialise repo: %v", err)
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
				storage: st,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}
			got, err := uc.Index(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Index() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_IsIndex(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name                string
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		args                args
		want                bool
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			want: true,
		},
		{
			name: "notFound",
			args: args{
				ctx:  context.Background(),
				name: "test_index1",
			},
			want: false,
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

			uc := &UseCase{
				servers: s,
				logger:  logger.New(loggerInstance),
				indexes: map[string]int32{
					"test_index": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}
			got, err := uc.IsIndex(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_SetToIndex(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx     context.Context
		index   string
		key     string
		val     string
		uniques bool
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
				ctx:     context.Background(),
				index:   "test_index",
				key:     "test_key",
				val:     "test_val",
				uniques: true,
			},
			want: 1,
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().SetToIndex(gomock.Any(), gomock.Any()).Return(
					&storage.SetResponse{}, nil)
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx:     context.Background(),
				index:   "test_index1",
				key:     "test_key",
				val:     "test_val",
				uniques: true,
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			wantErr:             true,
		},
		{
			name: "serverNotFound",
			args: args{
				ctx:     context.Background(),
				index:   "test_index",
				key:     "test_key",
				val:     "test_val",
				uniques: true,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			wantErr: true,
		},
		{
			name: "storageError",
			args: args{
				ctx:     context.Background(),
				index:   "test_index",
				key:     "test_key",
				val:     "test_val",
				uniques: true,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().SetToIndex(gomock.Any(), gomock.Any()).Return(
					nil, servers.ErrNotFound)
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
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			got, err := uc.SetToIndex(tt.args.ctx, tt.args.index, tt.args.key, tt.args.val, tt.args.uniques)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetToIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetToIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Size(t *testing.T) {
	srv := struct {
		serv *servers.Server
	}{}

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name                string
		serversBehavior     serversBehavior
		grpcStorageBehavior gStorageBehavior
		args                args
		want                uint64
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			want: 6,
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Size(gomock.Any(), gomock.Any()).Return(
					&storage.IndexSizeResponse{Size: 6}, nil)
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_index1",
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			wantErr:             true,
		},
		{
			name: "serverNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(nil, false)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			wantErr: true,
		},
		{
			name: "storageError",
			args: args{
				ctx:  context.Background(),
				name: "test_index",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Size(gomock.Any(), gomock.Any()).Return(
					nil, servers.ErrNotFound)
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
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			got, err := uc.Size(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}
