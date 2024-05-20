package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	servers2 "itisadb/internal/balancer"
	"itisadb/internal/core"
	"itisadb/internal/grpc-storage/storage"
	mockusecase "itisadb/internal/handler/mocks/usecase"
	"itisadb/internal/service/servers"
	api "itisadb/pkg/api/balancer"
	"reflect"
	"testing"
)

type mockUseCase func(*mockusecase.MockIUseCase)

func TestHandler_AttachToObject(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx context.Context
		r   *api.BalancerAttachToObjectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerAttachToObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToObjectRequest{
					Dst: "object1",
					Src: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerAttachToObjectResponse{},
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToObjectRequest{
					Dst: "object3",
					Src: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					core.ErrObjectNotFound)
			},
			want:    &api.BalancerAttachToObjectResponse{},
			wantErr: status.Error(codes.ResourceExhausted, core.ErrObjectNotFound.Error()),
		},
		{
			name: "notFound#2",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToObjectRequest{
					Dst: "object7",
					Src: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					status.Error(codes.ResourceExhausted, storage.ErrNotFound.Error()))
			},
			want:    &api.BalancerAttachToObjectResponse{},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrNotFound.Error()),
		},
		{
			name: "circularAttachment",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToObjectRequest{
					Dst: "object3",
					Src: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					status.Error(codes.PermissionDenied, storage.ErrCircularAttachment.Error()))
			},
			want:    &api.BalancerAttachToObjectResponse{},
			wantErr: status.Error(codes.PermissionDenied, storage.ErrCircularAttachment.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.AttachToObject(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AttachToObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttachToObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Connect(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerConnectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerConnectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerConnectRequest{
					Address:   "192.168.0.22:890",
					Total:     33,
					Available: 22,
					Server:    1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Connect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int32(1), nil)
			},
			want: &api.BalancerConnectResponse{
				Status:       "connected successfully",
				ServerNumber: 1,
			},
		},
		{
			name: "fileErr",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerConnectRequest{
					Address:   "192.168.0.22:890",
					Total:     33,
					Available: 32,
					Server:    1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Connect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), servers.ErrInternal)
			},
			wantErr: status.Error(codes.Internal, servers.ErrInternal.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.Connect(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connect() got = %v, want %v", got, tt.want)
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
		r   *api.BalancerDeleteRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerDeleteResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteRequest{
					Key:    "key",
					Server: 21,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerDeleteResponse{},
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteRequest{
					Key:    "key3",
					Server: 21,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.NotFound, storage.ErrNotFound.Error()))
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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
		r   *api.BalancerDeleteAttrRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerDeleteAttrResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteAttrRequest{
					Key:    "key",
					Object: "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerDeleteAttrResponse{},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteAttrRequest{
					Key:    "key",
					Object: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()))
			},
			want:    &api.BalancerDeleteAttrResponse{},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
		{
			name: "attrNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteAttrRequest{
					Key:    "key",
					Object: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.NotFound, storage.ErrNotFound.Error()))
			},
			want:    &api.BalancerDeleteAttrResponse{},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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
		r   *api.BalancerDeleteObjectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerDeleteObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteObjectRequest{
					Object: "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerDeleteObjectResponse{},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteObjectRequest{
					Object: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()))
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
		{
			name: "unavailable",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteObjectRequest{
					Object: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.Unavailable, servers2.ErrUnavailable.Error()))
			},
			wantErr: status.Error(codes.Unavailable, servers2.ErrUnavailable.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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

func TestHandler_Disconnect(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerDisconnectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerDisconnectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerDisconnectRequest{
					ServerNumber: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Disconnect(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerDisconnectResponse{},
		},
		{
			name: "ctxCancelled",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerDisconnectRequest{
					ServerNumber: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Disconnect(gomock.Any(), gomock.Any()).Return(context.Canceled)
			},
			wantErr: status.Error(codes.Canceled, context.Canceled.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.Disconnect(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Disconnect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Disconnect() got = %v, want %v", got, tt.want)
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
		r   *api.BalancerGetRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerGetResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r:   &api.BalancerGetRequest{},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)
			},
			want: &api.BalancerGetResponse{
				Value: "1",
			},
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				r:   &api.BalancerGetRequest{},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.NotFound, storage.ErrNotFound.Error()))
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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
		r   *api.BalancerGetFromObjectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerGetFromObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerGetFromObjectRequest{
					Object: "object",
					Key:    "qwe",
					Server: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetFromObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("1", nil)
			},
			want: &api.BalancerGetFromObjectResponse{
				Value: "1",
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerGetFromObjectRequest{
					Object: "object",
					Key:    "qwe",
					Server: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetFromObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()))
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
		{
			name: "attrNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerGetFromObjectRequest{
					Object: "object",
					Key:    "qwe",
					Server: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetFromObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.NotFound, storage.ErrNotFound.Error()))
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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
		ctx     context.Context
		request *api.BalancerObjectToJSONRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerObjectToJSONRequest
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerObjectToJSONRequest{
					Name: "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).
					Return(`{"object":"qwe","values":"2"}`, nil)
			},
			want: &api.BalancerObjectToJSONRequest{
				Name: "{\"object\":\"qwe\",\"values\":\"2\"}",
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerObjectToJSONRequest{
					Name: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()))
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.ObjectToJSON(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ObjectToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && (got.Object != tt.want.Name) {
				t.Errorf("ObjectToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Object(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerObjectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerObjectRequest{
					Name: "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Object(gomock.Any(), gomock.Any()).Return(int32(0), nil)
			},
			want: &api.BalancerObjectResponse{},
		},
		{
			name: "ErrSomethingExists",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerObjectRequest{
					Name: "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Object(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()))
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.Object(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Object() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Object() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_IsObject(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerIsObjectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerIsObjectResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIsObjectRequest{
					Name: "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().IsObject(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			want: &api.BalancerIsObjectResponse{
				Ok: true,
			},
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIsObjectRequest{
					Name: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().IsObject(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			want: &api.BalancerIsObjectResponse{
				Ok: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.IsObject(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("IsObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Servers(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerServersRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerServersResponse
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				request: &api.BalancerServersRequest{},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Servers().Return([]string{
					fmt.Sprintf("s#%d Avaliable: %d MB, Total: %d MB", 1, 1, 1),
					fmt.Sprintf("s#%d Avaliable: %d MB, Total: %d MB", 2, 2, 2),
					fmt.Sprintf("s#%d Avaliable: %d MB, Total: %d MB", 3, 3, 3),
				})
			},
			want: &api.BalancerServersResponse{
				ServersInfo: "s#1 Avaliable: 1 MB, Total: 1 MB<br>s#2 Avaliable: 2 MB, Total: 2 MB<br>s#3 Avaliable: 3 MB, Total: 3 MB",
			},
		},
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				request: &api.BalancerServersRequest{},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Servers().Return([]string{
					fmt.Sprintf("s#%d Avaliable: %d MB, Total: %d MB", 1, 1, 1),
				})
			},
			want: &api.BalancerServersResponse{
				ServersInfo: "s#1 Avaliable: 1 MB, Total: 1 MB",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.Servers(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Balancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Balancer() got = %v, want %v", got, tt.want)
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
		r   *api.BalancerSetRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerSetResponse
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerSetRequest{
					Key:     "key",
					Value:   "value",
					Server:  0,
					Uniques: false,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			want: &api.BalancerSetResponse{
				SavedTo: 1,
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerSetRequest{
					Key:     "key",
					Value:   "value",
					Server:  0,
					Uniques: false,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.Set(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		r   *api.BalancerSetToObjectRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerSetToObjectResponse
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerSetToObjectRequest{
					Key:     "key",
					Value:   "value",
					Uniques: false,
					Object:  "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			want: &api.BalancerSetToObjectResponse{
				SavedTo: 1,
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerSetToObjectRequest{
					Key:     "key",
					Value:   "value",
					Uniques: false,
					Object:  "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.SetToObject(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		ctx     context.Context
		request *api.BalancerObjectSizeRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerObjectSizeResponse
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerObjectSizeRequest{
					Name: "object",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Size(gomock.Any(), gomock.Any()).
					Return(uint64(1), nil)
			},
			want: &api.BalancerObjectSizeResponse{
				Size: 1,
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerObjectSizeRequest{
					Name: "object2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Size(gomock.Any(), gomock.Any()).
					Return(uint64(0), errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.Size(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}
