package rest

import (
	"github.com/golang/mock/gomock"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mocks "itisadb/internal/memory-balancer/handler/mocks/usecase"
	"testing"
)

type mockUseCase func(*mocks.MockIUseCase)

func TestHandler_get(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return("qwe", nil)
			},
			rJSON:    `{"key":"qwe"}`,
			want:     "qwe",
			wantCode: 200,
		},
		{
			name: "notFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.NotFound, "not found"))
			},
			rJSON:    `{"key":"qwe1"}`,
			wantCode: 404,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"key":"qwe1"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.get(tt.args.ctx)

			if tt.want != "" {
				if b := string(tt.args.ctx.Response.Body()); b != tt.want {
					t.Errorf("Want {%s} got {%s}", tt.want, b)
				}
			}

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_set(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			rJSON:    `{"key":"qwe", "value":"qwe", "uniques": true}`,
			want:     "1",
			wantCode: 200,
		},
		{
			name: "exists",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), status.Error(codes.AlreadyExists, "exists"))
			},
			rJSON:    `{"key":"qwe1"}`,
			want:     "1",
			wantCode: 409,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"key":"qwe1"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.set(tt.args.ctx)

			if tt.want != "" {
				if b := string(tt.args.ctx.Response.Body()); b != tt.want {
					t.Errorf("Want {%s} got {%s}", tt.want, b)
				}
			}

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_del(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			rJSON:    `{"key":"qwe"}`,
			wantCode: 200,
		},
		{
			name: "notFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.NotFound, "not found"))
			},
			rJSON:    `{"key":"qwe1"}`,
			wantCode: 404,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"key":"qwe1"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.del(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_getFromIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetFromIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("qwe", nil)
			},
			rJSON:    `{"key":"qwe", "uniques": true, "index": "q"}`,
			want:     "qwe",
			wantCode: 200,
		},
		{
			name: "notFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetFromIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.NotFound, "not found"))
			},
			rJSON:    `{"key":"qwe1", "uniques": true, "index": "q"}`,
			wantCode: 404,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetFromIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.ResourceExhausted, "index not found"))
			},
			rJSON:    `{"key":"qwe", "uniques": true, "index": "q2"}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetFromIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"key":"qwe1"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.getFromIndex(tt.args.ctx)

			if tt.want != "" {
				if b := string(tt.args.ctx.Response.Body()); b != tt.want {
					t.Errorf("Want {%s} got {%s}", tt.want, b)
				}
			}

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_setToIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			rJSON:    `{"key":"qwe", "value":"qwe", "uniques": true, "index": "q"}`,
			want:     "1",
			wantCode: 200,
		},
		{
			name: "exists",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), status.Error(codes.AlreadyExists, "exists"))
			},
			rJSON:    `{"key":"qwe1", "value":"qwe", "uniques": true, "index": "q"}`,
			want:     "1",
			wantCode: 409,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.ResourceExhausted, "index not found"))
			},
			rJSON:    `{"key":"qwe", "value":"qwe", "uniques": true, "index": "q2"}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().SetToIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"key":"qwe1"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.setToIndex(tt.args.ctx)

			if tt.want != "" {
				if b := string(tt.args.ctx.Response.Body()); b != tt.want {
					t.Errorf("Want {%s} got {%s}", tt.want, b)
				}
			}

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_getIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetIndex(gomock.Any(), gomock.Any()).
					Return(map[string]string{"i": "1", "v": "2"}, nil)
			},
			rJSON:    `{"index":"qwe"}`,
			want:     `{"i":"1","v":"2"}`,
			wantCode: 200,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetIndex(gomock.Any(), gomock.Any()).
					Return(map[string]string{}, status.Error(codes.ResourceExhausted, "index not found"))
			},
			rJSON:    `{"index":"qwe2"}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetIndex(gomock.Any(), gomock.Any()).
					Return(map[string]string{}, status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"index":"qwe"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.getIndex(tt.args.ctx)

			if tt.want != "" {
				if b := string(tt.args.ctx.Response.Body()); b != tt.want {
					t.Errorf("Want {%s} got {%s}", tt.want, b)
				}
			}

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_index(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Index(gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			rJSON:    `{"index":"qwe"}`,
			wantCode: 200,
		},
		{
			name: "somethingExists",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Index(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.AlreadyExists, "something exists"))
			},
			rJSON:    `{"index":"qwe2"}`,
			wantCode: 409,
		},
		{
			name: "invalidName",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Index(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.InvalidArgument, "invalid name"))
			},
			rJSON:    `{"index":""}`,
			wantCode: 400,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Index(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"index":"qwe"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.index(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_delIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			rJSON:    `{"index":"qwe"}`,
			wantCode: 200,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.ResourceExhausted, "not found"))
			},
			rJSON:    `{"index":""}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().DeleteIndex(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"index":"qwe"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.delIndex(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_sizeIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Size(gomock.Any(), gomock.Any()).
					Return(uint64(33), nil)
			},
			rJSON:    `{"index":"qwe"}`,
			want:     "33",
			wantCode: 200,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Size(gomock.Any(), gomock.Any()).
					Return(uint64(0), status.Error(codes.ResourceExhausted, "not found"))
			},
			rJSON:    `{"index":""}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Size(gomock.Any(), gomock.Any()).
					Return(uint64(0), status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"index":"qwe"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.indexSize(tt.args.ctx)

			if tt.want != "" {
				if b := string(tt.args.ctx.Response.Body()); b != tt.want {
					t.Errorf("Want {%s} got {%s}", tt.want, b)
				}
			}

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_isIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().IsIndex(gomock.Any(), gomock.Any()).
					Return(true, nil)
			},
			rJSON:    `{"name":"qwe"}`,
			want:     `{"isIndex":true}`,
			wantCode: 200,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"key":"qwe1"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.isIndex(tt.args.ctx)

			if tt.want != "" {
				if b := string(tt.args.ctx.Response.Body()); b != tt.want {
					t.Errorf("Want {%s} got {%s}", tt.want, b)
				}
			}

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_attachIndex(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	logicmock := mocks.NewMockIUseCase(c)
	h := New(logicmock)

	type args struct {
		ctx *fasthttp.RequestCtx
	}

	tests := []struct {
		name        string
		args        args
		mockUseCase mockUseCase
		rJSON       string
		want        string
		wantCode    int
	}{
		{
			name: "ok",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			rJSON:    `{"dst":"qwe", "src":"qwe3"}`,
			wantCode: 200,
		},
		{
			name: "indexNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.ResourceExhausted, "not found"))
			},
			rJSON:    `{"dst":"qwe2", "src":"qwe"}`,
			wantCode: 410,
		},
		{
			name: "circularAttachment",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.PermissionDenied, "circular attachment not allowed"))
			},
			rJSON:    `{"dst":"qwe", "src":"qwe"}`,
			wantCode: 403,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().AttachToIndex(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"dst":"qwe", "src":"qwe"}`,
			wantCode: 503,
		},
		{
			name: "badRequest",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {},
			rJSON:       `{"dst":"qwe", "src":"qwe"}"`,
			wantCode:    400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockUseCase(logicmock)

			tt.args.ctx.Request.AppendBodyString(tt.rJSON)

			h.attachIndex(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}
