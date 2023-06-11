package grpc

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	mockusecase "itisadb/internal/memory-balancer/handler/grpc/mocks/usecase"
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
		name        string
		args        args
		mockUseCase mockUseCase
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

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
		name    string
		args    args
		want    *api.BalancerConnectResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.Connect(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerDeleteResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.Delete(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerDeleteAttrResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.DeleteAttr(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerDeleteIndexResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.DeleteIndex(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerDisconnectResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.Disconnect(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerGetResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.Get(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerGetFromIndexResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.GetFromIndex(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerGetIndexResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.GetIndex(tt.args.ctx, tt.args.request)
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
		name    string
		args    args
		want    *api.BalancerIndexResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.Index(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerIsIndexResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.IsIndex(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
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
		name    string
		args    args
		want    *api.BalancerServersResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
		name    string
		args    args
		want    *api.BalancerSetResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
		name    string
		args    args
		want    *api.BalancerSetToIndexResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
		name    string
		args    args
		want    *api.BalancerIndexSizeResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
