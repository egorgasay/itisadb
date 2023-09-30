package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"itisadb/internal/handler/converterr"
	"itisadb/internal/handler/mocks/usecase"
	"itisadb/internal/schema"
	servers2 "itisadb/internal/servers"
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
	case "/object": // handle object endpoint
		switch string(ctx.Method()) {
		case fasthttp.MethodGet: // get object
			h.ObjectToJSON(ctx)
		case fasthttp.MethodPost: // create object if not exists
			h.object(ctx)
		case fasthttp.MethodDelete: // delete object
			h.delObject(ctx)
		}
	case "/object/": // handle values in object endpoint
		switch string(ctx.Method()) {
		case fasthttp.MethodGet: // get from object
			h.getFromObject(ctx)
		case fasthttp.MethodPost: // set to object
			h.setToObject(ctx)
		case fasthttp.MethodDelete: // delete attribute
			h.delFromObject(ctx)
		}
	case "/object/size": // size of object
		h.objectSize(ctx)
	case "/object/is": // is object exists
		h.isObject(ctx)
	case "/object/attach": // attach one object to another
		h.attachObject(ctx)
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

func (h *Handler) getFromObject(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.GetFromObjectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	value, err := h.logic.GetFromObject(ctx, r.Object, r.Key, r.Server)
	if err != nil {
		err = converterr.GetFromObject(err)
		if errors.Is(err, converterr.ErrNotFound) {
			ctx.Error(converterr.ErrNotFound.Error(), fasthttp.StatusNotFound)
		} else if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(converterr.ErrObjectNotFound.Error(), fasthttp.StatusGone)
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

func (h *Handler) setToObject(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.SetToObjectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	v, err := h.logic.SetToObject(ctx, r.Object, r.Key, r.Value, r.Uniques)
	if err != nil {
		err = converterr.SetToObject(err)
		if errors.Is(err, converterr.ErrExists) {
			ctx.Error(fmt.Sprint(v), fasthttp.StatusConflict)
		} else if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(converterr.ErrObjectNotFound.Error(), fasthttp.StatusGone)
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

func (h *Handler) delFromObject(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.DelFromObjectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = h.logic.DeleteAttr(ctx, r.Object, r.Key)
	if err != nil {
		err = converterr.DelFromObject(err)
		if errors.Is(err, converterr.ErrNotFound) {
			ctx.Error(converterr.ErrNotFound.Error(), fasthttp.StatusNotFound)
		} else if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(converterr.ErrObjectNotFound.Error(), fasthttp.StatusGone)
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

func (h *Handler) ObjectToJSON(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.ObjectToJSONRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	value, err := h.logic.ObjectToJSON(ctx, r.Object)
	if err != nil {
		err = converterr.ObjectToJSON(err)
		if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(converterr.ErrObjectNotFound.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(converterr.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")

	v, err := json.Marshal(value)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(v)

}

func (h *Handler) object(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.ObjectToJSONRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	_, err = h.logic.Object(ctx, r.Object)
	if err != nil {
		err = converterr.Object(err)
		if errors.Is(err, converterr.ErrExists) {
			ctx.Error(err.Error(), fasthttp.StatusConflict)
		} else if errors.Is(err, converterr.ErrInvalidName) {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) delObject(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.DelObjectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = h.logic.DeleteObject(ctx, r.Object)
	if err != nil {
		err = converterr.DelObject(err)
		if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(err.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) objectSize(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.SizeObjectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	size, err := h.logic.Size(ctx, r.Object)
	if err != nil {
		err = converterr.SizeObject(err)
		if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(err.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte(fmt.Sprint(size)))
}

func (h *Handler) isObject(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.IsObjectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	is, err := h.logic.IsObject(ctx, r.Name)
	if err != nil {
		err = converterr.IsObject(err)
		if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(err.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write([]byte(fmt.Sprintf(`{"isObject":%v}`, is)))
}

func (h *Handler) attachObject(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.AttachRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = h.logic.AttachToObject(ctx, r.Dst, r.Src)
	if err != nil {
		err = converterr.AttachObject(err)
		if errors.Is(err, converterr.ErrObjectNotFound) {
			ctx.Error(err.Error(), fasthttp.StatusGone)
		} else if errors.Is(err, converterr.ErrUnavailable) {
			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
		} else if errors.Is(err, converterr.ErrCircularAttachment) {
			ctx.Error(err.Error(), fasthttp.StatusForbidden)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) connect(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.ConnectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	snum, err := h.logic.Connect(r.Address, r.Available, r.Total, r.Server)
	if err != nil {
		if errors.Is(err, servers2.ErrAlreadyExists) {
			ctx.Error(err.Error(), fasthttp.StatusConflict)
		}
		if errors.Is(err, servers2.ErrInternal) {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte(fmt.Sprint(snum)))
}

func (h *Handler) disconnect(ctx *fasthttp.RequestCtx) {
	r, err := dataFromRequest[schema.DisconnectRequest](&ctx.Request)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = h.logic.Disconnect(ctx, r.Server)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			ctx.Error(err.Error(), fasthttp.StatusRequestTimeout)
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
