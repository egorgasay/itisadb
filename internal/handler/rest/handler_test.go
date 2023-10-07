package rest

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/handler/mocks/usecase"
	servers2 "itisadb/internal/servers"
	"itisadb/internal/service/servers"
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

func TestHandler_getFromObject(t *testing.T) {
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
				c.EXPECT().GetFromObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("qwe", nil)
			},
			rJSON:    `{"key":"qwe", "uniques": true, "object": "q"}`,
			want:     "qwe",
			wantCode: 200,
		},
		{
			name: "notFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetFromObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.NotFound, "not found"))
			},
			rJSON:    `{"key":"qwe1", "uniques": true, "object": "q"}`,
			wantCode: 404,
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetFromObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.ResourceExhausted, "object not found"))
			},
			rJSON:    `{"key":"qwe", "uniques": true, "object": "q2"}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().GetFromObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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

			h.getFromObject(tt.args.ctx)

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

func TestHandler_setToObject(t *testing.T) {
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
				c.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			rJSON:    `{"key":"qwe", "value":"qwe", "uniques": true, "object": "q"}`,
			want:     "1",
			wantCode: 200,
		},
		{
			name: "exists",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(1), status.Error(codes.AlreadyExists, "exists"))
			},
			rJSON:    `{"key":"qwe1", "value":"qwe", "uniques": true, "object": "q"}`,
			want:     "1",
			wantCode: 409,
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.ResourceExhausted, "object not found"))
			},
			rJSON:    `{"key":"qwe", "value":"qwe", "uniques": true, "object": "q2"}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().SetToObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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

			h.setToObject(tt.args.ctx)

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

func TestHandler_getObject(t *testing.T) {
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
				c.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).
					Return(`{"object":"qwe","values":"2"}`, nil)
			},
			rJSON:    `{"object":"qwe"}`,
			want:     `"{\"object\":\"qwe\",\"values\":\"2\"}"`,
			wantCode: 200,
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.ResourceExhausted, "object not found"))
			},
			rJSON:    `{"object":"qwe2"}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().ObjectToJSON(gomock.Any(), gomock.Any()).
					Return("", status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"object":"qwe"}`,
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

			h.ObjectToJSON(tt.args.ctx)

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

func TestHandler_object(t *testing.T) {
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
				c.EXPECT().Object(gomock.Any(), gomock.Any()).
					Return(int32(1), nil)
			},
			rJSON:    `{"object":"qwe"}`,
			wantCode: 200,
		},
		{
			name: "somethingExists",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Object(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.AlreadyExists, "something exists"))
			},
			rJSON:    `{"object":"qwe2"}`,
			wantCode: 409,
		},
		{
			name: "invalidName",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Object(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.InvalidArgument, "invalid name"))
			},
			rJSON:    `{"object":""}`,
			wantCode: 400,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Object(gomock.Any(), gomock.Any()).
					Return(int32(0), status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"object":"qwe"}`,
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

			h.object(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_delObject(t *testing.T) {
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
				c.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			rJSON:    `{"object":"qwe"}`,
			wantCode: 200,
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.ResourceExhausted, "not found"))
			},
			rJSON:    `{"object":""}`,
			wantCode: 410,
		},
		{
			name: "unavailable",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).
					Return(status.Error(codes.Unavailable, "service unavailable"))
			},
			rJSON:    `{"object":"qwe"}`,
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

			h.delObject(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_sizeObject(t *testing.T) {
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
			rJSON:    `{"object":"qwe"}`,
			want:     "33",
			wantCode: 200,
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Size(gomock.Any(), gomock.Any()).
					Return(uint64(0), status.Error(codes.ResourceExhausted, "not found"))
			},
			rJSON:    `{"object":""}`,
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
			rJSON:    `{"object":"qwe"}`,
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

			h.objectSize(tt.args.ctx)

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

func TestHandler_isObject(t *testing.T) {
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
				c.EXPECT().IsObject(gomock.Any(), gomock.Any()).
					Return(true, nil)
			},
			rJSON:    `{"name":"qwe"}`,
			want:     `{"isObject":true}`,
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

			h.isObject(tt.args.ctx)

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

func TestHandler_attachObject(t *testing.T) {
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
				c.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			rJSON:    `{"dst":"qwe", "src":"qwe3"}`,
			wantCode: 200,
		},
		{
			name: "objectNotFound",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).
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
				c.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).
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
				c.EXPECT().AttachToObject(gomock.Any(), gomock.Any(), gomock.Any()).
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

			h.attachObject(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_connect(t *testing.T) {
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
				c.EXPECT().Connect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), nil)
			},
			rJSON:    `{"address":"127.0.0.1:897", "total":100, "available":100, "server":1}`,
			wantCode: 200,
		},
		{
			name: "Exists",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Connect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), servers2.ErrAlreadyExists)
			},
			rJSON:    `{"address":"127.0.0.1:897", "total":100, "available":100, "server":1}`,
			wantCode: 409,
		},
		{
			name: "internal",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Connect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int32(0), servers.ErrInternal)
			},
			rJSON:    `{"address":"127.0.0.1:897", "total":100, "available":100, "server":1}`,
			wantCode: 500,
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

			h.connect(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}

func TestHandler_disconnect(t *testing.T) {
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
				c.EXPECT().Disconnect(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			rJSON:    `{"server":1}`,
			wantCode: 200,
		},
		{
			name: "contextDeadlineExceed",
			args: args{
				ctx: &fasthttp.RequestCtx{Request: fasthttp.Request{}},
			},
			mockUseCase: func(c *mocks.MockIUseCase) {
				c.EXPECT().Disconnect(gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
			},
			rJSON:    `{"server":2}`,
			wantCode: fasthttp.StatusRequestTimeout,
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

			h.disconnect(tt.args.ctx)

			if code := tt.args.ctx.Response.StatusCode(); code != tt.wantCode {
				t.Errorf("Want code {%d} got {%d}", tt.wantCode, code)
			}
		})
	}
}
