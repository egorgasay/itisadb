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

func TestHandler_AttachToObject(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx context.Context
		r   *api.AttachToObjectRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.AttachToObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToObjectRequest{
					Dst: "test1",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any()).Return(ram, nil)
			},
			want: &api.AttachToObjectResponse{
				Ram: apiRam,
			},
		},
		{
			name: "dstNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToObjectRequest{
					Dst: "test3",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any()).Return(ram, storage.ErrObjectNotFound)
			},
			want: &api.AttachToObjectResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.NotFound, storage.ErrObjectNotFound.Error()),
		},
		{
			name: "circularAttachment",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToObjectRequest{
					Dst: "test3",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any()).Return(ram, storage.ErrCircularAttachment)
			},
			want: &api.AttachToObjectResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.PermissionDenied, storage.ErrCircularAttachment.Error()),
		},
		{
			name: "somethingExists",
			args: args{
				ctx: context.Background(),
				r: &api.AttachToObjectRequest{
					Dst: "test4",
					Src: "test2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any()).Return(ram, storage.ErrSomethingExists)
			},
			want: &api.AttachToObjectResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)
			got, err := h.AttachToObject(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AttachToObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttachToObject() got = %v, want %v", got, tt.want)
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
					Name: "object",
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
					Name: "object",
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
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteAttrRequest{
					Key:  "test2",
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any()).Return(ram, storage.ErrObjectNotFound)
			},
			want: &api.DeleteAttrResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
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

func TestHandler_DeleteObject(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.DeleteObjectRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.DeleteObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteObjectRequest{
					Object: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteObject(gomock.Any()).Return(ram, nil)
			},
			want: &api.DeleteObjectResponse{
				Ram: apiRam,
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.DeleteObjectRequest{
					Object: "object2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().DeleteObject(gomock.Any()).Return(ram, storage.ErrObjectNotFound)
			},
			want: &api.DeleteObjectResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.DeleteObject(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteObject() got = %v, want %v", got, tt.want)
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

func TestHandler_GetFromObject(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.GetFromObjectRequest
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
				r: &api.GetFromObjectRequest{
					Key:  "test",
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).Return(ram, "test", nil)
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
				r: &api.GetFromObjectRequest{
					Key:  "test",
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).Return(ram, "", storage.ErrNotFound)
			},
			want: &api.GetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.GetFromObjectRequest{
					Key:  "test",
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().GetFromObject(gomock.Any(), gomock.Any()).Return(ram, "", storage.ErrObjectNotFound)
			},
			want: &api.GetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.GetFromObject(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetFromObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFromObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_ObjectToJSON(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.ObjectToJSONRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.ObjectToJSONResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.ObjectToJSONRequest{
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().ObjectToJSON(gomock.Any()).
					Return(ram,
						"{\n\t\"isObject\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
						nil)
			},
			want: &api.ObjectToJSONResponse{
				Ram:    apiRam,
				Object: "{\n\t\"isObject\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.ObjectToJSONRequest{
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().ObjectToJSON(gomock.Any()).Return(ram, "", storage.ErrObjectNotFound)
			},
			want: &api.ObjectToJSONResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.ObjectToJSON(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ObjectToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_NewObject(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.NewObjectRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.NewObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.NewObjectRequest{
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().NewObject(gomock.Any()).Return(ram, nil)
			},

			want: &api.NewObjectResponse{
				Ram: apiRam,
			},
		},
		{
			name: "emptyObjectName",
			args: args{
				ctx: context.Background(),
				r: &api.NewObjectRequest{
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().NewObject(gomock.Any()).Return(ram, storage.ErrEmptyObjectName)
			},

			want: &api.NewObjectResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.InvalidArgument, storage.ErrEmptyObjectName.Error()),
		},
		{
			name: "somethingExists",
			args: args{
				ctx: context.Background(),
				r: &api.NewObjectRequest{
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().NewObject(gomock.Any()).Return(ram, storage.ErrSomethingExists)
			},

			want: &api.NewObjectResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.NewObject(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewObject() got = %v, want %v", got, tt.want)
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

func TestHandler_SetToObject(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.SetToObjectRequest
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
				r: &api.SetToObjectRequest{
					Name:   "object",
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ram, nil)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
		},
		{
			name: "alreadyExists",
			args: args{
				ctx: context.Background(),
				r: &api.SetToObjectRequest{
					Name:   "object",
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ram, storage.ErrAlreadyExists)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrAlreadyExists.Error()),
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.SetToObjectRequest{
					Name:   "object",
					Key:    "key",
					Value:  "value",
					Unique: true,
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ram, storage.ErrObjectNotFound)
			},
			want: &api.SetResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.useCaseMock(logicmock)

			got, err := h.SetToObject(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetToObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetToObject() got = %v, want %v", got, tt.want)
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
		r   *api.ObjectSizeRequest
	}
	tests := []struct {
		name        string
		args        args
		useCaseMock useCaseMock
		want        *api.ObjectSizeResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.ObjectSizeRequest{
					Name: "object",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Size(gomock.Any()).Return(ram, uint64(1), nil)
			},
			want: &api.ObjectSizeResponse{
				Ram:  apiRam,
				Size: 1,
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.ObjectSizeRequest{
					Name: "object2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Size(gomock.Any()).Return(ram, uint64(0), storage.ErrObjectNotFound)
			},
			want: &api.ObjectSizeResponse{
				Ram: apiRam,
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
		{
			name: "notAnObject",
			args: args{
				ctx: context.Background(),
				r: &api.ObjectSizeRequest{
					Name: "object2",
				},
			},
			useCaseMock: func(usecase.IUseCase) {
				logicmock.EXPECT().Size(gomock.Any()).Return(ram, uint64(0), storage.ErrSomethingExists)
			},
			want: &api.ObjectSizeResponse{
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
