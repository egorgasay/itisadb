package rest

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"grpc-storage/internal/memory-balancer/schema"
	"grpc-storage/internal/memory-balancer/usecase"
	"log"
	"strings"
)

type Handler struct {
	logic *usecase.UseCase
}

func New(logic *usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}

func BindJSON(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}

func (h *Handler) ServeHTTP(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/get":
		h.get(ctx)
	case "/set":
		h.set(ctx, false)
	case "/unique-set":
		h.set(ctx, true)
	case "/servers":
		h.servers(ctx)
	}
}

func (h *Handler) get(ctx *fasthttp.RequestCtx) {
	var r schema.GetRequest
	b := ctx.Request.Body()
	log.Println(string(b))
	if err := BindJSON(b, &r); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	value, err := h.logic.Get(ctx, r.Key, r.Server)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(value))
}

func (h *Handler) set(ctx *fasthttp.RequestCtx, uniques bool) {
	var r schema.SetRequest
	b := ctx.Request.Body()
	if err := BindJSON(b, &r); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	setTo, err := h.logic.Set(ctx, r.Key, r.Value, r.Server, uniques)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(fmt.Sprint(setTo)))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(fmt.Sprint(setTo)))
}

func (h *Handler) servers(ctx *fasthttp.RequestCtx) {
	servers := h.logic.Servers()
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(strings.Join(servers, "<br>")))
}
