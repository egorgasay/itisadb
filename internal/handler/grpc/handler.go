package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/core"
	"itisadb/internal/handler/mocks/usecase"
	"itisadb/internal/servers"
	"itisadb/pkg/api"
	"strings"
)

type Handler struct {
	api.UnimplementedItisaDBServer
	logic mocks.IUseCase
}

func New(logic mocks.IUseCase) *Handler {
	return &Handler{logic: logic}
}
func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	setTo, err := h.logic.Set(ctx, r.Server, r.Key, r.Value, r.Uniques)
	if err != nil {
		return nil, err
	}

	return &api.SetResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) SetToObject(ctx context.Context, r *api.SetToObjectRequest) (*api.SetToObjectResponse, error) {
	setTo, err := h.logic.SetToObject(ctx, r.Server, r.Object, r.Key, r.Value, r.Uniques)
	if err != nil {
		return nil, err
	}

	return &api.SetToObjectResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	value, err := h.logic.Get(ctx, r.Server, r.Key)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, core.ErrUnknownServer) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}

		return nil, err
	}

	return &api.GetResponse{
		Value: value,
	}, nil
}

func (h *Handler) GetFromObject(ctx context.Context, r *api.GetFromObjectRequest) (*api.GetFromObjectResponse, error) {
	value, err := h.logic.GetFromObject(ctx, r.Server, r.GetObject(), r.GetKey())
	if err != nil {
		if errors.Is(err, core.ErrNoData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, core.ErrUnknownServer) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}

		return nil, err
	}

	return &api.GetFromObjectResponse{
		Value: value,
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	err := h.logic.Delete(ctx, r.GetServer(), r.Key)
	resp := &api.DeleteResponse{}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) AttachToObject(ctx context.Context, r *api.AttachToObjectRequest) (*api.AttachToObjectResponse, error) {
	err := h.logic.AttachToObject(ctx, r.Dst, r.Src)
	if err != nil {
		if errors.Is(err, core.ErrObjectNotFound) {
			return nil, status.Error(codes.ResourceExhausted, core.ErrObjectNotFound.Error())
		}
		return nil, err
	}

	return &api.AttachToObjectResponse{}, nil
}

func (h *Handler) DeleteObject(ctx context.Context, r *api.DeleteObjectRequest) (*api.DeleteObjectResponse, error) {
	err := h.logic.DeleteObject(ctx, r.Object)
	if err != nil {
		return nil, err
	}

	return &api.DeleteObjectResponse{}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.ConnectRequest) (*api.ConnectResponse, error) {
	serverNum, err := h.logic.Connect(request.GetAddress(), request.GetAvailable(), request.GetTotal(), request.Server)
	if err != nil {
		if errors.Is(err, servers.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, err
	}

	return &api.ConnectResponse{
		Status:       "connected successfully",
		ServerNumber: serverNum,
	}, nil
}

func (h *Handler) Object(ctx context.Context, request *api.ObjectRequest) (*api.ObjectResponse, error) {
	_, err := h.logic.Object(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.ObjectResponse{}, nil
}

func (h *Handler) ObjectToJSON(ctx context.Context, request *api.ObjectToJSONRequest) (*api.ObjectToJSONResponse, error) {
	m, err := h.logic.ObjectToJSON(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.ObjectToJSONResponse{
		Object: m,
	}, nil
}

func (h *Handler) IsObject(ctx context.Context, request *api.IsObjectRequest) (*api.IsObjectResponse, error) {
	ok, err := h.logic.IsObject(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.IsObjectResponse{
		Ok: ok,
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.DeleteAttrRequest) (*api.DeleteAttrResponse, error) {
	err := h.logic.DeleteAttr(ctx, r.GetKey(), r.GetObject())
	if err != nil {
		if errors.Is(err, core.ErrObjectNotFound) {
			return &api.DeleteAttrResponse{}, status.Error(codes.ResourceExhausted, core.ErrObjectNotFound.Error())
		}

		return &api.DeleteAttrResponse{}, err
	}

	return &api.DeleteAttrResponse{}, nil
}

func (h *Handler) Size(ctx context.Context, request *api.ObjectSizeRequest) (*api.ObjectSizeResponse, error) {
	size, err := h.logic.Size(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.ObjectSizeResponse{
		Size: size,
	}, nil
}

func (h *Handler) Disconnect(ctx context.Context, request *api.DisconnectRequest) (*api.DisconnectResponse, error) {
	err := h.logic.Disconnect(ctx, request.GetServerNumber())
	if err != nil {
		if errors.Is(err, context.Canceled) { // TODO: add everywhere
			return nil, status.Error(codes.Canceled, context.Canceled.Error())
		}
		return nil, err
	}

	return &api.DisconnectResponse{}, nil
}

func (h *Handler) Servers(ctx context.Context, request *api.ServersRequest) (*api.ServersResponse, error) {
	t, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(t)

	servers := h.logic.Servers()
	s := strings.Join(servers, "\n")
	return &api.ServersResponse{
		ServersInfo: s,
	}, nil
}

func (h *Handler) Authenticate(ctx context.Context, request *api.AuthRequest) (*api.AuthResponse, error) {
	token, err := h.logic.Authenticate(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		if errors.Is(err, core.ErrWrongCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, err
	}

	return &api.AuthResponse{Token: token}, nil
}
