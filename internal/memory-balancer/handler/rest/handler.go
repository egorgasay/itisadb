package rest

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"itisadb/internal/memory-balancer/handler/converterr"
	"itisadb/internal/memory-balancer/schema"
	"itisadb/internal/memory-balancer/usecase"
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
	case "/": // key:value endpoint
		switch string(ctx.Method()) {
		case fasthttp.MethodGet: // get value
			h.get(ctx)
		case fasthttp.MethodPost: // set value
			h.set(ctx)
		case fasthttp.MethodDelete: // delete value
			h.get(ctx)
		}
	case "/index": // handle index endpoint
		switch string(ctx.Method()) {
		case fasthttp.MethodGet: // get index

		case fasthttp.MethodPost: // create index
		case fasthttp.MethodDelete: // delete index
		}
	case "/index/": // handle values in index endpoint
		switch string(ctx.Method()) {
		case fasthttp.MethodGet: // get from index
		case fasthttp.MethodPost: // set to index
		case fasthttp.MethodDelete: // delete attribute
		}
	case "/index/size": // size of index
		//
	case "/index/is": // is index exists
		//
	case "/index/attach": // attach one index to another
		//
	case "/connect": // connect to balancer
		//
	case "/disconnect": // disconnect from balancer
		//
	case "/servers": // servers info
		h.servers(ctx)
	}
}

func getDataFromRequest[T any](r *fasthttp.Request) (t T, err error) {
	b := r.Body()
	if err = BindJSON(b, &r); err != nil {
		return t, err
	}

	return t, nil
}

func (h *Handler) get(ctx *fasthttp.RequestCtx) {
	r, err := getDataFromRequest[schema.GetRequest](&ctx.Request)
	if err != nil {
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

func (h *Handler) set(ctx *fasthttp.RequestCtx) {
	r, err := getDataFromRequest[schema.SetRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	setTo, err := h.logic.Set(ctx, r.Key, r.Value, r.Server, r.Uniques)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(fmt.Sprint(setTo)))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(fmt.Sprint(setTo)))
}

func (h *Handler) del(ctx *fasthttp.RequestCtx) {
	r, err := getDataFromRequest[schema.DelRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = h.logic.Delete(ctx, r.Key, r.Server)
	if err != nil {
		err = converterr.Get(err)
		if errors.Is(err, converterr.ErrNotFound) {
			ctx.Error(converterr.ErrNotFound.Error(), fasthttp.StatusNotFound)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(converterr.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) servers(ctx *fasthttp.RequestCtx) {
	servers := h.logic.Servers()
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(strings.Join(servers, "<br>")))
}
