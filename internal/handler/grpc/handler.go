package grpc

import (
	"context"
	"fmt"
	"itisadb/internal/domains"
	"itisadb/internal/handler/converterr"
	"itisadb/pkg/api"
	"strings"
)

type Handler struct {
	api.UnimplementedItisaDBServer
	core domains.Core
}

func New(logic domains.Core) *Handler {
	return &Handler{core: logic}
}
func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	setTo, err := h.core.Set(ctx, r.Server, r.Key, r.Value, r.Uniques)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.SetResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) SetToObject(ctx context.Context, r *api.SetToObjectRequest) (*api.SetToObjectResponse, error) {
	setTo, err := h.core.SetToObject(ctx, r.Server, r.Object, r.Key, r.Value, r.Uniques)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.SetToObjectResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	value, err := h.core.Get(ctx, r.Server, r.Key)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.GetResponse{
		Value: value,
	}, nil
}

func (h *Handler) GetFromObject(ctx context.Context, r *api.GetFromObjectRequest) (*api.GetFromObjectResponse, error) {
	value, err := h.core.GetFromObject(ctx, r.Server, r.GetObject(), r.GetKey())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.GetFromObjectResponse{
		Value: value,
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	err := h.core.Delete(ctx, r.Server, r.Key)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteResponse{}, nil
}

func (h *Handler) AttachToObject(ctx context.Context, r *api.AttachToObjectRequest) (*api.AttachToObjectResponse, error) {
	err := h.core.AttachToObject(ctx, r.Server, r.Dst, r.Src)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.AttachToObjectResponse{}, nil
}

func (h *Handler) DeleteObject(ctx context.Context, r *api.DeleteObjectRequest) (*api.DeleteObjectResponse, error) {
	err := h.core.DeleteObject(ctx, r.Server, r.Object)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteObjectResponse{}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.ConnectRequest) (*api.ConnectResponse, error) {
	serverNum, err := h.core.Connect(request.GetAddress(), request.GetAvailable(), request.GetTotal())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ConnectResponse{
		Status: "connected successfully",
		Server: serverNum,
	}, nil
}

func (h *Handler) Object(ctx context.Context, r *api.ObjectRequest) (*api.ObjectResponse, error) {
	_, err := h.core.Object(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ObjectResponse{}, nil
}

func (h *Handler) ObjectToJSON(ctx context.Context, r *api.ObjectToJSONRequest) (*api.ObjectToJSONResponse, error) {
	m, err := h.core.ObjectToJSON(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ObjectToJSONResponse{
		Object: m,
	}, nil
}

func (h *Handler) IsObject(ctx context.Context, r *api.IsObjectRequest) (*api.IsObjectResponse, error) {
	ok, err := h.core.IsObject(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.IsObjectResponse{
		Ok: ok,
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.DeleteAttrRequest) (*api.DeleteAttrResponse, error) {
	err := h.core.DeleteAttr(ctx, r.Server, r.GetKey(), r.GetObject())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteAttrResponse{}, nil
}

func (h *Handler) Size(ctx context.Context, r *api.ObjectSizeRequest) (*api.ObjectSizeResponse, error) {
	size, err := h.core.Size(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ObjectSizeResponse{
		Size: size,
	}, nil
}

func (h *Handler) Disconnect(ctx context.Context, r *api.DisconnectRequest) (*api.DisconnectResponse, error) {
	err := h.core.Disconnect(ctx, r.Server)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DisconnectResponse{}, nil
}

func (h *Handler) Servers(ctx context.Context, r *api.ServersRequest) (*api.ServersResponse, error) {
	t, err := getToken(ctx)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	fmt.Println(t)

	servers := h.core.Servers()
	s := strings.Join(servers, "\n")
	return &api.ServersResponse{
		ServersInfo: s,
	}, nil
}

func (h *Handler) Authenticate(ctx context.Context, request *api.AuthRequest) (*api.AuthResponse, error) {
	token, err := h.core.Authenticate(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.AuthResponse{Token: token}, nil
}
