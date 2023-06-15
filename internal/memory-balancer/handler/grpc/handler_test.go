package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/grpc-storage/storage"
	mockusecase "itisadb/internal/memory-balancer/handler/grpc/mocks/usecase"
	"itisadb/internal/memory-balancer/servers"
	"itisadb/internal/memory-balancer/usecase"
	api "itisadb/pkg/api/balancer"
	"reflect"
	"testing"
)

type mockUseCase func(*mockusecase.MockIUseCase)

func TestHandler_AttachToIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx context.Context
		r   *api.BalancerAttachToIndexRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerAttachToIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToIndexRequest{
					Dst: "index1",
					Src: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerAttachToIndexResponse{},
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToIndexRequest{
					Dst: "index3",
					Src: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrIndexNotFound)
			},
			want:    &api.BalancerAttachToIndexResponse{},
			wantErr: status.Error(codes.ResourceExhausted, usecase.ErrIndexNotFound.Error()),
		},
		{
			name: "notFound#2",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToIndexRequest{
					Dst: "index7",
					Src: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					status.Error(codes.ResourceExhausted, storage.ErrNotFound.Error()))
			},
			want:    &api.BalancerAttachToIndexResponse{},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrNotFound.Error()),
		},
		{
			name: "circularAttachment",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerAttachToIndexRequest{
					Dst: "index3",
					Src: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					status.Error(codes.PermissionDenied, storage.ErrCircularAttachment.Error()))
			},
			want:    &api.BalancerAttachToIndexResponse{},
			wantErr: status.Error(codes.PermissionDenied, storage.ErrCircularAttachment.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.AttachToIndex(tt.args.ctx, tt.args.r)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AttachToIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttachToIndex() got = %v, want %v", got, tt.want)
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
					Key:   "key",
					Index: "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerDeleteAttrResponse{},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteAttrRequest{
					Key:   "key",
					Index: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteAttr(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()))
			},
			want:    &api.BalancerDeleteAttrResponse{},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
		{
			name: "attrNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteAttrRequest{
					Key:   "key",
					Index: "index2",
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

func TestHandler_DeleteIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.BalancerDeleteIndexRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerDeleteIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteIndexRequest{
					Index: "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &api.BalancerDeleteIndexResponse{},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteIndexRequest{
					Index: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()))
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
		{
			name: "unavailable",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerDeleteIndexRequest{
					Index: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.Unavailable, servers.ErrUnavailable.Error()))
			},
			wantErr: status.Error(codes.Unavailable, servers.ErrUnavailable.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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

func TestHandler_GetFromIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.BalancerGetFromIndexRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerGetFromIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerGetFromIndexRequest{
					Index:  "index",
					Key:    "qwe",
					Server: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetFromIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("1", nil)
			},
			want: &api.BalancerGetFromIndexResponse{
				Value: "1",
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerGetFromIndexRequest{
					Index:  "index",
					Key:    "qwe",
					Server: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetFromIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()))
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
		{
			name: "attrNotFound",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerGetFromIndexRequest{
					Index:  "index",
					Key:    "qwe",
					Server: 1,
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetFromIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.NotFound, storage.ErrNotFound.Error()))
			},
			wantErr: status.Error(codes.NotFound, storage.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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

func TestHandler_GetIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerGetIndexRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerGetIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerGetIndexRequest{
					Name: "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetIndex(gomock.Any(), gomock.Any()).Return(map[string]string{"index": ""}, nil)
			},
			want: &api.BalancerGetIndexResponse{
				Index: map[string]string{"index": ""},
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerGetIndexRequest{
					Name: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().GetIndex(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()))
			},
			wantErr: status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.GetIndex(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Index(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerIndexRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIndexRequest{
					Name: "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Index(gomock.Any(), gomock.Any()).Return(int32(0), nil)
			},
			want: &api.BalancerIndexResponse{},
		},
		{
			name: "ErrSomethingExists",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIndexRequest{
					Name: "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Index(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()))
			},
			wantErr: status.Error(codes.AlreadyExists, storage.ErrSomethingExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.Index(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_IsIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx     context.Context
		request *api.BalancerIsIndexRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerIsIndexResponse
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIsIndexRequest{
					Name: "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().IsIndex(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			want: &api.BalancerIsIndexResponse{
				Ok: true,
			},
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIsIndexRequest{
					Name: "index2",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().IsIndex(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			want: &api.BalancerIsIndexResponse{
				Ok: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.IsIndex(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("IsIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsIndex() got = %v, want %v", got, tt.want)
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
				t.Errorf("Servers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Servers() got = %v, want %v", got, tt.want)
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

func TestHandler_SetToIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mockusecase.NewMockIUseCase(c)
	h := New(logicmock)
	type args struct {
		ctx context.Context
		r   *api.BalancerSetToIndexRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerSetToIndexResponse
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerSetToIndexRequest{
					Key:     "key",
					Value:   "value",
					Uniques: false,
					Index:   "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			want: &api.BalancerSetToIndexResponse{
				SavedTo: 1,
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				r: &api.BalancerSetToIndexRequest{
					Key:     "key",
					Value:   "value",
					Uniques: false,
					Index:   "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), errors.New("unexpected error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			got, err := h.SetToIndex(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		ctx     context.Context
		request *api.BalancerIndexSizeRequest
	}
	tests := []struct {
		mockUseCase mockUseCase
		name        string
		args        args
		want        *api.BalancerIndexSizeResponse
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIndexSizeRequest{
					Name: "index",
				},
			},
			mockUseCase: func(*mockusecase.MockIUseCase) {
				logicmock.EXPECT().Size(gomock.Any(), gomock.Any()).
					Return(uint64(1), nil)
			},
			want: &api.BalancerIndexSizeResponse{
				Size: 1,
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				request: &api.BalancerIndexSizeRequest{
					Name: "index2",
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
