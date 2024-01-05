package rest

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"itisadb/internal/domains"
)

type Handler struct {
	core domains.Balancer
}

func New(logic domains.Balancer) *Handler {
	return &Handler{core: logic}
}

func BindJSON(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}

func (h *Handler) ServeHTTP(ctx *fasthttp.RequestCtx) {
	//switch string(ctx.Path()) {
	//case "/": // key:value endpoint
	//	switch string(ctx.Method()) {
	//	case fasthttp.MethodGet: // get value
	//		h.get(ctx)
	//	case fasthttp.MethodPost: // set value
	//		h.set(ctx)
	//	case fasthttp.MethodDelete: // delete value
	//		h.del(ctx)
	//	}
	//case "/object": // handle object endpoint
	//	switch string(ctx.Method()) {
	//	case fasthttp.MethodGet: // get object
	//		h.ObjectToJSON(ctx)
	//	case fasthttp.MethodPost: // create object if not exists
	//		h.object(ctx)
	//	case fasthttp.MethodDelete: // delete object
	//		h.delObject(ctx)
	//	}
	//case "/object/": // handle values in object endpoint
	//	switch string(ctx.Method()) {
	//	case fasthttp.MethodGet: // get from object
	//		h.getFromObject(ctx)
	//	case fasthttp.MethodPost: // set to object
	//		h.setToObject(ctx)
	//	case fasthttp.MethodDelete: // delete attribute
	//		h.delFromObject(ctx)
	//	}
	//case "/object/size": // size of object
	//	h.objectSize(ctx)
	//case "/object/is": // is object exists
	//	h.isObject(ctx)
	//case "/object/attach": // attach one object to another
	//	h.attachObject(ctx)
	//case "/connect": // connect to balancer
	//	h.connect(ctx)
	//case "/disconnect": // disconnect from balancer
	//	h.disconnect(ctx)
	//case "/balancer": // balancer info
	//	h.balancer(ctx)
	//}
}

//func dataFromRequest[T any](r *fasthttp.Request) (t T, err error) {
//	b := r.Body()
//	if err = BindJSON(b, &r); err != nil {
//		return t, err
//	}
//
//	return t, nil
//}
//
//func (h *Handler) get(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[GetRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	value, err := h.core.Get(ctx, r.Server, r.Key)
//	if err != nil {
//		if errors.Is(err, constants.ErrNotFound) {
//			ctx.Error(constants.ErrNotFound.Error(), fasthttp.StatusNotFound)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(constants.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.SetBody([]byte(value))
//}
//
//func (h *Handler) set(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[SetRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	setTo, err := h.core.Set(ctx, r.Server, r.Key, r.Value, r.Uniques)
//	if err != nil {
//		if errors.Is(err, constants.ErrAlreadyExists) {
//			ctx.Error(fmt.Sprint(setTo), fasthttp.StatusConflict)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(constants.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.SetBody([]byte(fmt.Sprint(setTo)))
//}
//
//func (h *Handler) del(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[DelRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	err = h.core.Delete(ctx, r.Server, r.Key)
//	if err != nil {
//		if errors.Is(err, constants.ErrNotFound) {
//			ctx.Error(constants.ErrNotFound.Error(), fasthttp.StatusNotFound)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(constants.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//}
//
//func (h *Handler) getFromObject(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[GetFromObjectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	value, err := h.core.GetFromObject(ctx, r.Server, r.Object, r.Key)
//	if err != nil {
//		if errors.Is(err, constants.ErrNotFound) {
//			ctx.Error(constants.ErrNotFound.Error(), fasthttp.StatusNotFound)
//		} else if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(constants.ErrObjectNotFound.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(constants.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.SetBody([]byte(value))
//}
//
//func (h *Handler) setToObject(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[SetToObjectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	v, err := h.core.SetToObject(ctx, r.Server, r.Object, r.Key, r.Value, r.Uniques)
//	if err != nil {
//		if errors.Is(err, constants.ErrAlreadyExists) {
//			ctx.Error(fmt.Sprint(v), fasthttp.StatusConflict)
//		} else if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(constants.ErrObjectNotFound.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(constants.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.SetBody([]byte(fmt.Sprint(v)))
//}
//
//func (h *Handler) delFromObject(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[DelFromObjectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	err = h.core.DeleteAttr(ctx, r.Server, r.Object, r.Key)
//	if err != nil {
//		if errors.Is(err, constants.ErrNotFound) {
//			ctx.Error(constants.ErrNotFound.Error(), fasthttp.StatusNotFound)
//		} else if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(constants.ErrObjectNotFound.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(constants.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//}
//
//func (h *Handler) balancer(ctx *fasthttp.RequestCtx) {
//	balancer := h.core.Balancer()
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.SetBody([]byte(strings.Join(balancer, "<br>")))
//}
//
//func (h *Handler) ObjectToJSON(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[ObjectToJSONRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	value, err := h.core.ObjectToJSON(ctx, r.Server, r.Object)
//	if err != nil {
//		if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(constants.ErrObjectNotFound.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(constants.ErrUnavailable.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.Response.Header.Set("Content-Type", "application/json")
//
//	v, err := json.Marshal(value)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		return
//	}
//
//	ctx.SetBody(v)
//
//}
//
//func (h *Handler) object(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[ObjectToJSONRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	_, err = h.core.Object(ctx, r.Server, r.Object)
//	if err != nil {
//		if errors.Is(err, constants.ErrAlreadyExists) {
//			ctx.Error(err.Error(), fasthttp.StatusConflict)
//		} else if errors.Is(err, constants.ErrInvalidName) {
//			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//}
//
//func (h *Handler) delObject(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[DelObjectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	err = h.core.DeleteObject(ctx, r.Server, r.Object)
//	if err != nil {
//		if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(err.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//}
//
//func (h *Handler) objectSize(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[SizeObjectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	size, err := h.core.Size(ctx, r.Server, r.Object)
//	if err != nil {
//		if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(err.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.Write([]byte(fmt.Sprint(size)))
//}
//
//func (h *Handler) isObject(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[IsObjectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	is, err := h.core.IsObject(ctx, r.Server, r.Name)
//	if err != nil {
//		if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(err.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.Response.Header.Set("Content-Type", "application/json")
//	ctx.Write([]byte(fmt.Sprintf(`{"isObject":%v}`, is)))
//}
//
//func (h *Handler) attachObject(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[AttachRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	err = h.core.AttachToObject(ctx, r.Server, r.Dst, r.Src)
//	if err != nil {
//		if errors.Is(err, constants.ErrObjectNotFound) {
//			ctx.Error(err.Error(), fasthttp.StatusGone)
//		} else if errors.Is(err, constants.ErrUnavailable) {
//			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
//		} else if errors.Is(err, constants.ErrCircularAttachment) {
//			ctx.Error(err.Error(), fasthttp.StatusForbidden)
//		} else {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//}
//
//func (h *Handler) connect(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[ConnectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	snum, err := h.core.Connect(r.Address, r.Available, r.Total)
//	if err != nil {
//		if errors.Is(err, constants.ErrAlreadyExists) {
//			ctx.Error(err.Error(), fasthttp.StatusConflict)
//		}
//		if errors.Is(err, constants.ErrInternal) {
//			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.Write([]byte(fmt.Sprint(snum)))
//}
//
//func (h *Handler) disconnect(ctx *fasthttp.RequestCtx) {
//	r, err := dataFromRequest[DisconnectRequest](&ctx.Request)
//	if err != nil {
//		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
//		return
//	}
//
//	err = h.core.Disconnect(ctx, r.Server)
//	if err != nil {
//		if errors.Is(err, context.DeadlineExceeded) {
//			ctx.Error(err.Error(), fasthttp.StatusRequestTimeout)
//		}
//		return
//	}
//
//	ctx.SetStatusCode(fasthttp.StatusOK)
//}
