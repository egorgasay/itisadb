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
	"testing"
)

/*
	ALL TEST ARE SIMPLE NOW,
	TODO: ADD EDGE CASES
*/

type mockBehavior func(r *storagemock.MockStorageClient)

func TestServer_AttachToIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		dst string
		src string
	}
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      error
	}{
		{
			name: "Success",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().AttachToIndex(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
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
				cl.EXPECT().AttachToIndex(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				dst: "test2",
				src: "inner2",
			},
			wantErr: ErrUnavailable,
		},
		{
			name: "indexNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().AttachToIndex(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "index not found")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				dst: "test3",
				src: "inner3",
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

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
			if err := s.AttachToIndex(tt.args.ctx, tt.args.dst, tt.args.src); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("AttachToIndex() error = %v, wantErr %v", err, tt.wantErr)
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
				cl.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
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
					nil, status.Error(codes.Unavailable, "bad connection")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: ErrUnavailable,
		},
		{
			name: "keyNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "key was not found")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

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
			if err := s.Delete(tt.args.ctx, tt.args.Key); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_DeleteAttr(t *testing.T) {
	type args struct {
		ctx   context.Context
		attr  string
		index string
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
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
			args: args{
				ctx:   context.Background(),
				attr:  "test",
				index: "inner",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection")).AnyTimes()
			},
			args: args{
				ctx:   context.Background(),
				attr:  "test2",
				index: "inner2",
			},
			wantErr: ErrUnavailable,
		},
		{
			name: "indexNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "index not found")).AnyTimes()
			},
			args: args{
				ctx:   context.Background(),
				attr:  "test2",
				index: "inner2",
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

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
			if err := s.DeleteAttr(tt.args.ctx, tt.args.attr, tt.args.index); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("DeleteAttr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_DeleteIndex(t *testing.T) {
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
				cl.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
			args: args{
				ctx:  context.Background(),
				name: "test",
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection")).AnyTimes()
			},
			args: args{
				ctx:  context.Background(),
				name: "test2",
			},
			wantErr: ErrUnavailable,
		},
		{
			name: "indexNotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "index not found")).AnyTimes()
			},
			args: args{
				ctx:  context.Background(),
				name: "test2",
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

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
			if err := s.DeleteIndex(tt.args.ctx, tt.args.name); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
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
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&storage.GetResponse{Value: "test"}, nil).AnyTimes()
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
					nil, status.Error(codes.Unavailable, "bad connection")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: ErrUnavailable,
		},
		{
			name: "NotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			tt.mockBehavior(cl)

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
			got, err := s.Get(tt.args.ctx, tt.args.Key)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetFromIndex(t *testing.T) {
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
				cl.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).
					Return(&storage.GetResponse{Value: "test"}, nil).AnyTimes()
			},
		},
		{
			name: "badConnection",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.Unavailable, "bad connection")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				Key: "test2",
			},
			wantErr: ErrUnavailable,
		},
		{
			name: "NotFound",
			mockBehavior: func(cl *storagemock.MockStorageClient) {
				cl.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "not found")).AnyTimes()
			},
			args: args{
				ctx: context.Background(),
				Key: "test3",
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		c := gomock.NewController(t)
		defer c.Finish()
		cl := storagemock.NewMockStorageClient(c)
		tt.mockBehavior(cl)
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
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFromIndex(tt.args.ctx, tt.args.name, tt.args.Key)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("GetFromIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && tt.want != nil {
				t.Errorf("GetFromIndex() got = %v, want %v", got, tt.want)
				return
			} else if got == nil || tt.want == nil {
				return
			}

			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("GetFromIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetIndex(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *storage.GetIndexResponse
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				tries: 0,
				//storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			got, err := s.GetIndex(tt.args.ctx, tt.args.name)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetNumber(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
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
			s := &Server{
				tries: 0,
				//storage: cl,
				ram: RAM{
					available: 100,
					total:     100,
				},
				number: 1,
				mu:     &sync.RWMutex{},
			}
			if got := s.GetNumber(); got != tt.want {
				t.Errorf("GetNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetRAM(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   RAM
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			if got := s.GetRAM(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRAM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetTries(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			if got := s.GetTries(); got != tt.want {
				t.Errorf("GetTries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_IncTries(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			s.IncTries()
		})
	}
}

func TestServer_NewIndex(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			if err := s.NewIndex(tt.args.ctx, tt.args.name); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("NewIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_ResetTries(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			s.ResetTries()
		})
	}
}

func TestServer_Set(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	type args struct {
		ctx    context.Context
		Key    string
		Value  string
		unique bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			if err := s.Set(tt.args.ctx, tt.args.Key, tt.args.Value, tt.args.unique); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_SetToIndex(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	type args struct {
		ctx    context.Context
		name   string
		Key    string
		Value  string
		unique bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			if err := s.SetToIndex(tt.args.ctx, tt.args.name, tt.args.Key, tt.args.Value, tt.args.unique); (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("SetToIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Size(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *storage.IndexSizeResponse
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

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
			got, err := s.Size(tt.args.ctx, tt.args.name)
			if (err != nil) && (!errors.Is(err, tt.wantErr)) {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_setRAM(t *testing.T) {
	type fields struct {
		tries   uint
		storage storage.StorageClient
		ram     RAM
		number  int32
		mu      *sync.RWMutex
	}
	type args struct {
		ram *storage.Ram
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
			//c := gomock.NewController(t)
			//defer c.Finish()
			//cl := storagemock.NewMockStorageClient(c)
			//tt.mockBehavior(cl)

			//s := &Server{
			//	tries:   0,
			//	storage: cl,
			//	ram: RAM{
			//		available: 100,
			//		total:     100,
			//	},
			//	number: 1,
			//	mu:     &sync.RWMutex{},
			//}
		})
	}
}
