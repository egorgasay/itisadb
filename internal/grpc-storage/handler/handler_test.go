package handler

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/grpc-storage/storage"
	"itisadb/internal/grpc-storage/usecase"
	mockusecase "itisadb/internal/grpc-storage/usecase/mocks/usecase"
	api "itisadb/pkg/api/storage"
	"reflect"
	"testing"
)

type useCaseMock func(usecase.IUseCase)

var ram = usecase.RAM{}
var apiRam = &api.Ram{
	Total:     0,
	Available: 0,
}

func TestHandler_AttachToIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx context.Context
		r   *api.AttachToIndexRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.AttachToIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToIndexRequest{
					Dst: "test1",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any()).Return(ram, nil)
			},
			want: &api.AttachToIndexResponse{
				Ram: apiRam,
			},
		},
		{
			name: "dstNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToIndexRequest{
					Dst: "test3",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any()).Return(ram, storage.ErrIndexNotFound)
			},
			want: &api.AttachToIndexResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.NotFound, storage.ErrIndexNotFound.Error()),
		},
		{
			name: "circularAttachment",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToIndexRequest{
					Dst: "test3",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any()).Return(ram, storage.ErrCircularAttachment)
			},
			want: &api.AttachToIndexResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.PermissionDenied, storage.ErrCircularAttachment.Error()),
		},
		{
			name: "somethingExists",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToIndexRequest{
					Dst: "test4",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any()).Return(ram, storage.ErrSomethingExists)
			},
			want: &api.AttachToIndexResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)
			got, err := h.AttachToIndex(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AttachToIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttachToIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx context.Context
		r   *api.DeleteRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.DeleteResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteRequest{
					Key: "test1",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Delete(gomock.Any()).Return(ram, nil)
			},
			want: &api.DeleteResponse{
				Ram: apiRam,
			},
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteRequest{
					Key: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Delete(gomock.Any()).Return(ram, storage.ErrNotFound)
			},
			want: &api.DeleteResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.Delete(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_DeleteAttr(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx context.Context
		r   *api.DeleteAttrRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *api.DeleteAttrResponse
		useCaseMock useCaseMock
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteAttrRequest{
					Key:  "test1",
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any()).Return(ram, nil)
			},
			want: &api.DeleteAttrResponse{
				Ram: apiRam,
			},
		},
		{
			name: "keyNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteAttrRequest{
					Key:  "test2",
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any()).Return(ram, storage.ErrNotFound)
			},
			want: &api.DeleteAttrResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteAttrRequest{
					Key:  "test2",
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any()).Return(ram, storage.ErrIndexNotFound)
			},
			want: &api.DeleteAttrResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.DeleteAttr(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteAttr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAttr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_DeleteIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.DeleteIndexRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.DeleteIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteIndexRequest{
					Index: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteIndex(gomock.Any()).Return(ram, nil)
			},
			want: &api.DeleteIndexResponse{
				Ram: apiRam,
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteIndexRequest{
					Index: "index2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteIndex(gomock.Any()).Return(ram, storage.ErrIndexNotFound)
			},
			want: &api.DeleteIndexResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.DeleteIndex(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.GetRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.GetResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.GetRequest{
					Key: "test",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Get(gomock.Any()).Return(ram, "test", nil)
			},
			want: &api.GetResponse{
				Ram:   apiRam,
				Value: "test",
			},
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				r: &api.GetRequest{
					Key: "test",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Get(gomock.Any()).Return(ram, "", storage.ErrNotFound)
			},
			want: &api.GetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.Get(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GetFromIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.GetFromIndexRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.GetResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.GetFromIndexRequest{
					Key:  "test",
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).Return(ram, "test", nil)
			},
			want: &api.GetResponse{
				Ram:   apiRam,
				Value: "test",
			},
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				r: &api.GetFromIndexRequest{
					Key:  "test",
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).Return(ram, "", storage.ErrNotFound)
			},
			want: &api.GetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.GetFromIndexRequest{
					Key:  "test",
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().GetFromIndex(gomock.Any(), gomock.Any()).Return(ram, "", storage.ErrIndexNotFound)
			},
			want: &api.GetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.GetFromIndex(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetFromIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFromIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_IndexToJSON(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.IndexToJSONRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.IndexToJSONResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.IndexToJSONRequest{
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().IndexToJSON(gomock.Any()).
					Return(ram,
						"{\n\t\"isIndex\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
						nil)
			},
			want: &api.IndexToJSONResponse{
				Ram:   apiRam,
				Index: "{\n\t\"isIndex\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.IndexToJSONRequest{
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().IndexToJSON(gomock.Any()).Return(ram, "", storage.ErrIndexNotFound)
			},
			want: &api.IndexToJSONResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.IndexToJSON(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("IndexToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_NewIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.NewIndexRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.NewIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.NewIndexRequest{
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().NewIndex(gomock.Any()).Return(ram, nil)
			},

			want: &api.NewIndexResponse{
				Ram: apiRam,
			},
		},
		{
			name: "emptyIndexName",
			args: args{
				ctx: context.Background(),
				r: &api.NewIndexRequest{
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().NewIndex(gomock.Any()).Return(ram, storage.ErrEmptyIndexName)
			},

			want: &api.NewIndexResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.InvalidArgument, storage.ErrEmptyIndexName.Error()),
		},
		{
			name: "somethingExists",
			args: args{
				ctx: context.Background(),
				r: &api.NewIndexRequest{
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().NewIndex(gomock.Any()).Return(ram, storage.ErrSomethingExists)
			},

			want: &api.NewIndexResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.NewIndex(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Set(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.SetRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.SetResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.SetRequest{
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(ram, nil)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
		},
		{
			name: "alreadyExists",
			args: args{
				ctx: context.Background(),
				r: &api.SetRequest{
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(ram, storage.ErrAlreadyExists)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrAlreadyExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.Set(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_SetToIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.SetToIndexRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.SetResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.SetToIndexRequest{
					Name:   "index",
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ram, nil)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
		},
		{
			name: "alreadyExists",
			args: args{
				ctx: context.Background(),
				r: &api.SetToIndexRequest{
					Name:   "index",
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ram, storage.ErrAlreadyExists)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrAlreadyExists.Error()),
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.SetToIndexRequest{
					Name:   "index",
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ram, storage.ErrIndexNotFound)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.SetToIndex(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetToIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetToIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Size(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.IndexSizeRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.IndexSizeResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.IndexSizeRequest{
					Name: "index",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Size(gomock.Any()).Return(ram, uint64(1), nil)
			},
			want: &api.IndexSizeResponse{
				Ram:  apiRam,
				Size: 1,
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.IndexSizeRequest{
					Name: "index2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Size(gomock.Any()).Return(ram, uint64(0), storage.ErrIndexNotFound)
			},
			want: &api.IndexSizeResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
		{
			name: "notAnIndex",
			args: args{
				ctx: context.Background(),
				r: &api.IndexSizeRequest{
					Name: "index2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Size(gomock.Any()).Return(ram, uint64(0), storage.ErrSomethingExists)
			},
			want: &api.IndexSizeResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.Size(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}
