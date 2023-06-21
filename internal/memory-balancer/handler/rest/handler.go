package rest

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"itisadb/internal/memory-balancer/handler/converterr"
	mocks "itisadb/internal/memory-balancer/handler/mocks/usecase"
	"itisadb/internal/memory-balancer/schema"
	"strings"
)

type Handler struct {
	logic mocks.IUseCase
}

func New(logic mocks.IUseCase) *Handler {
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
			h.del(ctx)
		}
	case "/index": // handle index endpoint
		switch string(ctx.Method()) {
		case fasthttp.MethodGet: // get index
			h.getIndex(ctx)
		case fasthttp.MethodPost: // create index if not exists
			h.index(ctx)
		case fasthttp.MethodDelete: // delete index
			h.delIndex(ctx)
		}
	case "/index/": // handle values in index endpoint
		switch string(ctx.Method()) {
		case fasthttp.MethodGet: // get from index
			h.getFromIndex(ctx)
		case fasthttp.MethodPost: // set to index
			h.setToIndex(ctx)
		case fasthttp.MethodDelete: // delete attribute
			h.delFromIndex(ctx)
		}
	case "/index/size": // size of index
		h.indexSize(ctx)
	case "/index/is": // is index exists
		h.isIndex(ctx)
	case "/index/attach": // attach one index to another
		h.attachIndex(ctx)
	case "/connect": // connect to balancer
		h.connect(ctx)
	case "/disconnect": // disconnect from balancer
		h.disconnect(ctx)
	case "/servers": // servers info
		h.servers(ctx)
	}
}

func dataFromRequest[T any](r *fasthttp.Request) (t T, err error) {
	b := r.Body()
	if err = BindJSON(b, &r); err != nil {
		return t, err
	}

	return t, nil
}

func (h *Handler) get(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.GetRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	value, err := h.logic.Get(ctx, r.Key, r.Server)
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
	ctx.SetBody([]byte(value))
}

func (h *Handler) set(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.SetRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	setTo, err := h.logic.Set(ctx, r.Key, r.Value, r.Server, r.Uniques)
	if err != nil {
		err = converterr.Set(err)
		if errors.Is(err, converterr.ErrExists) {
			ctx.Error(fmt.Sprint(setTo), fasthttp.StatusConflict)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(converterr.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(fmt.Sprint(setTo)))
}

func (h *Handler) del(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.DelRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = h.logic.Delete(ctx, r.Key, r.Server)
	if err != nil {
		err = converterr.Del(err)
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

func (h *Handler) getFromIndex(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.GetFromIndexRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	value, err := h.logic.GetFromIndex(ctx, r.Index, r.Key, r.Server)
	if err != nil {
		err = converterr.GetFromIndex(err)
		if errors.Is(err, converterr.ErrNotFound) {
			ctx.Error(converterr.ErrNotFound.Error(), fasthttp.StatusNotFound)
		} else if errors.Is(err, converterr.ErrIndexNotFound) {
			ctx.Error(converterr.ErrIndexNotFound.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(converterr.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(value))
}

func (h *Handler) setToIndex(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.SetToIndexRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	v, err := h.logic.SetToIndex(ctx, r.Index, r.Key, r.Value, r.Uniques)
	if err != nil {
		err = converterr.SetToIndex(err)
		if errors.Is(err, converterr.ErrExists) {
			ctx.Error(fmt.Sprint(v), fasthttp.StatusConflict)
		} else if errors.Is(err, converterr.ErrIndexNotFound) {
			ctx.Error(converterr.ErrIndexNotFound.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(converterr.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(fmt.Sprint(v)))
}

func (h *Handler) delFromIndex(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.DelFromIndexRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = h.logic.DeleteAttr(ctx, r.Index, r.Key)
	if err != nil {
		err = converterr.DelFromIndex(err)
		if errors.Is(err, converterr.ErrNotFound) {
			ctx.Error(converterr.ErrNotFound.Error(), fasthttp.StatusNotFound)
		} else if errors.Is(err, converterr.ErrIndexNotFound) {
			ctx.Error(converterr.ErrIndexNotFound.Error(), fasthttp.StatusGone)
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

func (h *Handler) getIndex(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.GetIndexRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	value, err := h.logic.GetIndex(ctx, r.Index)
	if err != nil {
		err = converterr.GetIndex(err)
		if errors.Is(err, converterr.ErrIndexNotFound) {
			ctx.Error(converterr.ErrIndexNotFound.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(converterr.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	v, err := json.Marshal(value)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(v)

}

func (h *Handler) index(ctx *fasthttp.RequestCtx) {

}

func (h *Handler) delIndex(ctx *fasthttp.RequestCtx) {

}

func (h *Handler) indexSize(ctx *fasthttp.RequestCtx) {

}

func (h *Handler) isIndex(ctx *fasthttp.RequestCtx) {

}

func (h *Handler) attachIndex(ctx *fasthttp.RequestCtx) {

}

func (h *Handler) connect(ctx *fasthttp.RequestCtx) {

}

func (h *Handler) disconnect(ctx *fasthttp.RequestCtx) {

}
