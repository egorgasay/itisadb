package servers

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/pkg/api/storage"
	storagemock "itisadb/pkg/api/storage/gomocks"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

func TestServer_AttachToObject(t *testing.T) {
	type args struct {
		ctx context.Context
		dst string
		src string
	}
	tests := []struct {
		name         string
		mockBehavior func(r *storagemock.MockStorageClient)
		args         args
		wantCode     error
	}{
		{
			name: "Success",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().AttachToObject(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			args: args{
				ctx: context.Background(),
				dst: "test",
				src: "inner",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().AttachToObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx: context.Background(),
				dst: "test2",
				src: "inner2",
			},
			wantCode: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "objectNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().AttachToObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "object not found"))
			},
			args: args{
				ctx: context.Background(),
				dst: "test3",
				src: "inner3",
			},
			wantCode: status.Error(codes.NotFound, "object not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if err := s.AttachToObject(tt.args.ctx, tt.args.dst, tt.args.src); (err != nil) && (!errors.Is(err, tt.wantCode)) {
				t.Errorf("AttachToObject() error = %v, wantCode %v", err, tt.wantCode)
			}
		})
	}
}

func TestServer_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		Key string
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		wantErr      error
	}{
		{
			name: "Success",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			args: args{
				ctx: context.Background(),
				Key: "test",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "keyNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "key was not found"))
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: status.Error(codes.NotFound, "key was not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if err := s.Delete(tt.args.ctx, tt.args.Key); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Delete() error = %v, wantCode %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_DeleteAttr(t *testing.T) {
	type args struct {
		ctx    context.Context
		attr   string
		object string
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		wantErr      error
	}{
		{
			name: "Success",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			args: args{
				ctx:    context.Background(),
				attr:   "test",
				object: "inner",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx:    context.Background(),
				attr:   "test2",
				object: "inner2",
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "objectNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "object not found"))
			},
			args: args{
				ctx:    context.Background(),
				attr:   "test2",
				object: "inner2",
			},
			wantErr: status.Error(codes.NotFound, "object not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if err := s.DeleteAttr(tt.args.ctx, tt.args.attr, tt.args.object); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("DeleteAttr() error = %v, wantCode %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_DeleteObject(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		wantErr      error
	}{
		{
			name: "Success",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			args: args{
				ctx:  context.Background(),
				name: "test",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx:  context.Background(),
				name: "test2",
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "objectNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "object not found"))
			},
			args: args{
				ctx:  context.Background(),
				name: "test2",
			},
			wantErr: status.Error(codes.NotFound, "object not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if err := s.DeleteObject(tt.args.ctx, tt.args.name); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("DeleteObject() error = %v, wantCode %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		Key string
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		want         *storage.GetResponse
		wantErr      error
	}{
		{
			name: "Success",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&storage.GetResponse{Value: "test"}, nil)
			},
			args: args{
				ctx: context.Background(),
				Key: "test",
			},
			want: &storage.GetResponse{
				Value: "test",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "NotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found"))
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			got, err := s.Get(tt.args.ctx, tt.args.Key)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Get() error = %v, wantCode %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetFromObject(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		Key  string
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		want         *storage.GetResponse
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "test",
				Key:  "test",
			},
			want: &storage.GetResponse{
				Value: "test",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).
					Return(&storage.GetResponse{Value: "test"}, nil)
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "NotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found"))
			},
			args: args{
				ctx: context.Background(),
				Key: "test3",
			},
			wantErr: status.Error(codes.NotFound, "not found"),
		},
	}
	c := gomock.NewController(t)
	defer c.Finish()
	cl := storagemock.NewMockStorageClient(c)
	for _, tt := range tests {
		tt.mockBehavior(cl)
		s := &Server{
			tries:   atomic.Uint32{},
			storage: cl,
			ram: RAM{
				available: 100,
				total:     100,
			},
			number: 1,
			mu:     &sync.RWMutex{},
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFromObject(tt.args.ctx, tt.args.name, tt.args.Key)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("GetFromObject() error = %v, wantCode %v", err, tt.wantErr)
				return
			}

			if got == nil && tt.want != nil {
				t.Errorf("GetFromObject() got = %v, want %v", got, tt.want)
				return
			} else if got == nil || tt.want == nil {
				return
			}

			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("GetFromObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ObjectToJSON(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		want         *storage.ObjectToJSONResponse
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "TestServer_ObjectToJSON1",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).Return(
					&storage.ObjectToJSONResponse{
						Object: "{\n\t\"isObject\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
					}, nil,
				)
			},
			want: &storage.ObjectToJSONResponse{
				Object: "{\n\t\"isObject\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx:  context.Background(),
				name: "TestServer_ObjectToJSON2",
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "NotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found"))
			},
			args: args{
				ctx:  context.Background(),
				name: "TestServer_ObjectToJSON3",
			},
			wantErr: status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)
			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			got, err := s.ObjectToJSON(tt.args.ctx, tt.args.name)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("ObjectToJSON() error = %v, wantCode %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				return
			}

			if got != nil && tt.want != nil {
				if !reflect.DeepEqual(*got, *tt.want) {
					t.Errorf("ObjectToJSON() got = %v, want %v", got, tt.want)
				}
			} else {
				t.Errorf("ObjectToJSON() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestServer_NewObject(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior func(cl *storagemock.MockStorageClient)
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "TestServer_NewObject1",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().NewObject(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().NewObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			args: args{
				ctx:  context.Background(),
				name: "TestServer_NewObject2",
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if err := s.NewObject(tt.args.ctx, tt.args.name); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("NewObject() error = %v, wantCode %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Set(t *testing.T) {
	type args struct {
		ctx    context.Context
		Key    string
		Value  string
		unique bool
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior func(cl *storagemock.MockStorageClient)
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				Key:    "Key_Set",
				Value:  "test",
				unique: false,
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
		},
		{
			name: "badConnection",
			args: args{
				ctx:    context.Background(),
				Key:    "Key_Set",
				Value:  "test",
				unique: false,
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "AlreadyExists",
			args: args{
				ctx:    context.Background(),
				Key:    "Key_Set",
				Value:  "test",
				unique: true,
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.AlreadyExists, "already exists"))
			},
			wantErr: status.Error(codes.AlreadyExists, "already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if err := s.Set(tt.args.ctx, tt.args.Key, tt.args.Value, tt.args.unique); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Set() error = %v, wantCode %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_SetToObject(t *testing.T) {
	type args struct {
		ctx    context.Context
		name   string
		Key    string
		Value  string
		unique bool
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				name:   "TestServer_SetToObject",
				Key:    "Key_Set",
				Value:  "test",
				unique: false,
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().SetToObject(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
		},
		{
			name: "badConnection",
			args: args{
				ctx:    context.Background(),
				name:   "TestServer_SetToObject",
				Key:    "Key_Set",
				Value:  "test",
				unique: false,
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().SetToObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "AlreadyExists",
			args: args{
				ctx:    context.Background(),
				name:   "TestServer_SetToObject",
				Key:    "Key_Set",
				Value:  "test",
				unique: true,
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().SetToObject(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.AlreadyExists, "already exists"))
			},
			wantErr: status.Error(codes.AlreadyExists, "already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if err := s.SetToObject(tt.args.ctx, tt.args.name, tt.args.Key, tt.args.Value, tt.args.unique); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("SetToObject() error = %v, wantCode %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Size(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name         string
		mockBehavior func(cl *storagemock.MockStorageClient)
		args         args
		want         *storage.ObjectSizeResponse
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: "TestServer_Size",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Size(gomock.Any(), gomock.Any()).Return(
					&storage.ObjectSizeResponse{
						Size: 100,
					}, nil)
			},
			want: &storage.ObjectSizeResponse{
				Size: 100,
			},
		},
		{
			name: "badConnection",
			args: args{
				ctx:  context.Background(),
				name: "TestServer_Size",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Size(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection"))
			},
			wantErr: status.Error(codes.Unavailable, "bad connection"),
		},
		{
			name: "NotFound",
			args: args{
				ctx:  context.Background(),
				name: "TestServer_Size",
			},
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Size(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found"))
			},
			wantErr: status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

			s := &Server{
				tries:   atomic.Uint32{},
				storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			got, err := s.Size(tt.args.ctx, tt.args.name)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Size() error = %v, wantCode %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_setRAM(t *testing.T) {
	type args struct {
		ram *storage.AttachToObjectResponse
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				ram: &storage.AttachToObjectResponse{
					Ram: &storage.Ram{
						Available: 100,
						Total:     100,
					},
				},
			},
		},
		{
			name: "nil",
			args: args{
				ram: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				tries: atomic.Uint32{},
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			s.setRAM(tt.args.ram)

			if tt.args.ram == nil {
				return
			}
			if s.ram.available != tt.args.ram.Ram.Available {
				t.Errorf("setRAM() = %v, want %v", s.ram.available, tt.args.ram.Ram.Available)
			}
			if s.ram.total != tt.args.ram.Ram.Total {
				t.Errorf("setRAM() = %v, want %v", s.ram.total, tt.args.ram.Ram.Total)
			}
		})
	}
}
