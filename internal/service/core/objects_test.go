package core

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	servers2 "itisadb/internal/balancer"
	"itisadb/internal/service/balancer"
	serversmock "itisadb/internal/service/core/mocks/servers"
	repo "itisadb/internal/storage"
	"itisadb/pkg/api/storage"
	"itisadb/pkg/logger"
	"reflect"
	"sync"
	"testing"
)

type serversBehavior func(cl *serversmock.MockIServers)
type gStorageBehavior func(cl *storagemock.MockStorageClient)

func TestUseCase_AttachToObject(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
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
				cl.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
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
			name: "objectNotFound",
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
				cl.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).Return(
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
			srv.serv = balancer.NewRemoteServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &Core{
				balancer: s,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test1": 1,
					"test2": 2,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}

			if err = uc.AttachToObject(tt.args.ctx, tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("AttachToObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_DeleteAttr(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}

	type args struct {
		ctx    context.Context
		attr   string
		object string
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
				ctx:    context.Background(),
				attr:   "test1",
				object: "test_object",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx:    context.Background(),
				attr:   "test1",
				object: "test_object2",
			},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
		{
			name: "attrNotFound",
			args: args{
				ctx:    context.Background(),
				attr:   "test1",
				object: "test_object",
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
				ctx:    context.Background(),
				attr:   "test1",
				object: "test_object",
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

			if err = uc.DeleteAttr(tt.args.ctx, tt.args.attr, tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAttr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_DeleteObject(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
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
				name: "test_object",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_object33",
			},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			wantErr:             true,
		},
		{
			name: "badConnection",
			args: args{
				ctx:  context.Background(),
				name: "test_object",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			wantErr: true,
		},
		{
			name: "servNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_object",
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
			if err := uc.DeleteObject(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_GetFromObject(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}

	type args struct {
		ctx          context.Context
		object       string
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
				object:       "test_object",
				key:          "test1",
				serverNumber: 1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).Return(
					&storage.GetResponse{Value: "test1"}, nil)
			},
			want:    "test1",
			wantErr: false,
		},
		{
			name: "objectNotFound(fromServer)",
			args: args{
				ctx:          context.Background(),
				object:       "test_object",
				key:          "test1",
				serverNumber: 1,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).Return(
					&storage.GetResponse{}, status.Error(codes.NotFound, "object not found"))
			},
			wantErr: true,
		},
		{
			name: "objectNotFound",
			args: args{
				ctx:          context.Background(),
				object:       "test_ind33ex33",
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
				object:       "test_object",
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

			got, err := uc.GetFromObject(tt.args.ctx, tt.args.object, tt.args.key, tt.args.serverNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_ObjectToJSON(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
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
		want                string
		wantErr             bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "test_object",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).Return(
					&storage.ObjectToJSONResponse{
						Object: "[\n  {\n    \"name\": \"object1\",\n    \"isObject\": true,\n    \"values\": [\n      {\n        \"name\": \"key1\",\n        \"value\": \"value1\"\n      },\n      {\n        \"name\": \"key2\",\n        \"value\": \"value2\"\n      }\n    ]\n  }\n]",
					}, nil)
			},
			want: "[\n  {\n    \"name\": \"object1\",\n    \"isObject\": true,\n    \"values\": [\n      {\n        \"name\": \"key1\",\n        \"value\": \"value1\"\n      },\n      {\n        \"name\": \"key2\",\n        \"value\": \"value2\"\n      }\n    ]\n  }\n]",
		},
		{
			name: "objectNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_object33",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
			},
			wantErr: true,
		},
		{
			name: "objectNotFound(remoteStorage)",
			args: args{
				ctx:  context.Background(),
				name: "test_object",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).Return(
					&storage.ObjectToJSONResponse{}, status.Error(codes.NotFound, "object not found"))
			},
			wantErr: true,
		},
		{
			name: "serverNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_object",
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

			got, err := uc.ObjectToJSON(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ObjectToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Object(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
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
				name: "test_object",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().NewObject(gomock.Any(), gomock.Any()).Return(
					&storage.NewObjectResponse{}, nil)
			},
			want: 0,
		},
		{
			name: "create",
			args: args{
				ctx:  context.Background(),
				name: "test_object2",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServer().Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().NewObject(gomock.Any(), gomock.Any()).Return(
					&storage.NewObjectResponse{}, nil)
			},
			want: 0,
		}, {
			name: "noActiveClients",
			args: args{
				ctx:  context.Background(),
				name: "test_object2",
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
			srv.serv = balancer.NewRemoteServer(sc, 1)

			s := serversmock.NewMockIServers(c)
			tt.serversBehavior(s)

			uc := &Core{
				balancer: s,
				storage:  st,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test_object": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}
			got, err := uc.Object(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Object() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Object() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_IsObject(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
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
				name: "test_object",
			},
			want: true,
		},
		{
			name: "notFound",
			args: args{
				ctx:  context.Background(),
				name: "test_object1",
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
			srv.serv = balancer.NewRemoteServer(sc, 1)
			s := serversmock.NewMockIServers(c)

			uc := &Core{
				balancer: s,
				logger:   logger.New(loggerInstance),
				objects: map[string]int32{
					"test_object": 1,
				},
				mu:   sync.RWMutex{},
				pool: make(chan struct{}, 30000),
			}
			got, err := uc.IsObject(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_SetToObject(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
	}{}

	type args struct {
		ctx     context.Context
		object  string
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
				object:  "test_object",
				key:     "test_key",
				val:     "test_val",
				uniques: true,
			},
			want: 1,
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().SetToObject(gomock.Any(), gomock.Any()).Return(
					&storage.SetResponse{}, nil)
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx:     context.Background(),
				object:  "test_object1",
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
				object:  "test_object",
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
				object:  "test_object",
				key:     "test_key",
				val:     "test_val",
				uniques: true,
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().SetToObject(gomock.Any(), gomock.Any()).Return(
					nil, servers2.ErrNotFound)
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

			got, err := uc.SetToObject(tt.args.ctx, tt.args.object, tt.args.key, tt.args.val, tt.args.uniques)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetToObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetToObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Size(t *testing.T) {
	srv := struct {
		serv *balancer.RemoteServer
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
				name: "test_object",
			},
			want: 6,
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Size(gomock.Any(), gomock.Any()).Return(
					&storage.ObjectSizeResponse{Size: 6}, nil)
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_object1",
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {},
			serversBehavior:     func(cl *serversmock.MockIServers) {},
			wantErr:             true,
		},
		{
			name: "serverNotFound",
			args: args{
				ctx:  context.Background(),
				name: "test_object",
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
				name: "test_object",
			},
			serversBehavior: func(cl *serversmock.MockIServers) {
				cl.EXPECT().GetServerByID(gomock.Any()).Return(srv.serv, true)
			},
			grpcStorageBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Size(gomock.Any(), gomock.Any()).Return(
					nil, servers2.ErrNotFound)
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
