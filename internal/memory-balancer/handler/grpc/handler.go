package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mocks "itisadb/internal/memory-balancer/handler/mocks/usecase"
	"itisadb/internal/memory-balancer/servers"
	"itisadb/internal/memory-balancer/usecase"
	api "itisadb/pkg/api/balancer"
	"strings"
)

type Handler struct {
	api.UnimplementedBalancerServer
	logic mocks.IUseCase
}

func New(logic mocks.IUseCase) *Handler {
	return &Handler{logic: logic}
}
func (h *Handler) Set(ctx context.Context, r *api.BalancerSetRequest) (*api.BalancerSetResponse, error) {
	setTo, err := h.logic.Set(ctx, r.Key, r.Value, r.Server, r.Uniques)
	if err != nil {
		return nil, err
	}

	return &api.BalancerSetResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) SetToObject(ctx context.Context, r *api.BalancerSetToObjectRequest) (*api.BalancerSetToObjectResponse, error) {
	setTo, err := h.logic.SetToObject(ctx, r.Object, r.Key, r.Value, r.Uniques)
	if err != nil {
		return nil, err
	}

	return &api.BalancerSetToObjectResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.BalancerGetRequest) (*api.BalancerGetResponse, error) {
	value, err := h.logic.Get(ctx, r.Key, r.Server)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, usecase.ErrUnknownServer) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}

		return nil, err
	}

	return &api.BalancerGetResponse{
		Value: value,
	}, nil
}

func (h *Handler) GetFromObject(ctx context.Context, r *api.BalancerGetFromObjectRequest) (*api.BalancerGetFromObjectResponse, error) {
	value, err := h.logic.GetFromObject(ctx, r.GetObject(), r.GetKey(), r.GetServer())
	if err != nil {
		if errors.Is(err, usecase.ErrNoData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, usecase.ErrUnknownServer) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}

		return nil, err
	}

	return &api.BalancerGetFromObjectResponse{
		Value: value,
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.BalancerDeleteRequest) (*api.BalancerDeleteResponse, error) {
	err := h.logic.Delete(ctx, r.Key, r.Server)
	resp := &api.BalancerDeleteResponse{}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) AttachToObject(ctx context.Context, r *api.BalancerAttachToObjectRequest) (*api.BalancerAttachToObjectResponse, error) {
	err := h.logic.AttachToObject(ctx, r.Dst, r.Src)
	if err != nil {
		if errors.Is(err, usecase.ErrObjectNotFound) {
			return nil, status.Error(codes.ResourceExhausted, usecase.ErrObjectNotFound.Error())
		}
		return nil, err
	}

	return &api.BalancerAttachToObjectResponse{}, nil
}

func (h *Handler) DeleteObject(ctx context.Context, r *api.BalancerDeleteObjectRequest) (*api.BalancerDeleteObjectResponse, error) {
	err := h.logic.DeleteObject(ctx, r.Object)
	if err != nil {
		return nil, err
	}

	return &api.BalancerDeleteObjectResponse{}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.BalancerConnectRequest) (*api.BalancerConnectResponse, error) {
	serverNum, err := h.logic.Connect(request.GetAddress(), request.GetAvailable(), request.GetTotal(), request.Server)
	if err != nil {
		if errors.Is(err, servers.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, err
	}

	return &api.BalancerConnectResponse{
		Status:       "connected successfully",
		ServerNumber: serverNum,
	}, nil
}

func (h *Handler) Object(ctx context.Context, request *api.BalancerObjectRequest) (*api.BalancerObjectResponse, error) {
	_, err := h.logic.Object(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerObjectResponse{}, nil
}

func (h *Handler) ObjectToJSON(ctx context.Context, request *api.BalancerObjectToJSONRequest) (*api.BalancerObjectToJSONResponse, error) {
	m, err := h.logic.ObjectToJSON(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerObjectToJSONResponse{
		Object: m,
	}, nil
}

func (h *Handler) IsObject(ctx context.Context, request *api.BalancerIsObjectRequest) (*api.BalancerIsObjectResponse, error) {
	ok, err := h.logic.IsObject(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerIsObjectResponse{
		Ok: ok,
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.BalancerDeleteAttrRequest) (*api.BalancerDeleteAttrResponse, error) {
	err := h.logic.DeleteAttr(ctx, r.GetKey(), r.GetObject())
	if err != nil {
		if errors.Is(err, usecase.ErrObjectNotFound) {
			return &api.BalancerDeleteAttrResponse{}, status.Error(codes.ResourceExhausted, usecase.ErrObjectNotFound.Error())
		}

		return &api.BalancerDeleteAttrResponse{}, err
	}

	return &api.BalancerDeleteAttrResponse{}, nil
}

func (h *Handler) Size(ctx context.Context, request *api.BalancerObjectSizeRequest) (*api.BalancerObjectSizeResponse, error) {
	size, err := h.logic.Size(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerObjectSizeResponse{
		Size: size,
	}, nil
}

func (h *Handler) Disconnect(ctx context.Context, request *api.BalancerDisconnectRequest) (*api.BalancerDisconnectResponse, error) {
	err := h.logic.Disconnect(ctx, request.GetServerNumber())
	if err != nil {
		if errors.Is(err, context.Canceled) { // TODO: add everywhere
			return nil, status.Error(codes.Canceled, context.Canceled.Error())
		}
		return nil, err
	}

	return &api.BalancerDisconnectResponse{}, nil
}

func (h *Handler) Servers(ctx context.Context, request *api.BalancerServersRequest) (*api.BalancerServersResponse, error) {
	servers := h.logic.Servers()
	s := strings.Join(servers, "<br>")
	return &api.BalancerServersResponse{
		ServersInfo: s,
	}, nil
}

func (h *Handler) Authenticate(ctx context.Context, request *api.BalancerAuthRequest) (*api.BalancerAuthResponse, error) {
	token, err := h.logic.Authenticate(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &api.BalancerAuthResponse{Token: token}, nil
}
