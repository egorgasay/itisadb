package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/memory-balancer/servers"
	"itisadb/internal/memory-balancer/usecase"
	api "itisadb/pkg/api/balancer"
	"strings"
)

type Handler struct {
	api.UnimplementedBalancerServer
	logic *usecase.UseCase
}

func New(logic *usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}
func (h *Handler) Set(ctx context.Context, r *api.BalancerSetRequest) (*api.BalancerSetResponse, error) {
	setTo, err := h.logic.Set(ctx, r.Key, r.Value, r.Server, r.Uniques)
	if err != nil {
		if errors.Is(err, servers.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "")
		}
		return nil, err
	}

	return &api.BalancerSetResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) SetToIndex(ctx context.Context, r *api.BalancerSetToIndexRequest) (*api.BalancerSetResponse, error) {
	setTo, err := h.logic.SetToIndex(ctx, r.Index, r.Key, r.Value, r.Server, r.Uniques)
	if err != nil {
		if errors.Is(err, servers.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "")
		}
		return nil, err
	}

	return &api.BalancerSetResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.BalancerGetRequest) (*api.BalancerGetResponse, error) {
	value, err := h.logic.Get(ctx, r.Key, r.Server)

	if err != nil {
		if errors.Is(err, usecase.ErrNoData) {
			return &api.BalancerGetResponse{
				Value: err.Error(),
			}, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, usecase.ErrUnknownServer) {
			return &api.BalancerGetResponse{
				Value: err.Error(),
			}, status.Error(codes.Unavailable, err.Error())
		}

		return &api.BalancerGetResponse{
			Value: err.Error(),
		}, err
	}

	return &api.BalancerGetResponse{
		Value: value,
	}, nil
}

func (h *Handler) GetFromIndex(ctx context.Context, r *api.BalancerGetFromIndexRequest) (*api.BalancerGetResponse, error) {
	value, err := h.logic.GetFromIndex(ctx, r.GetIndex(), r.GetKey(), r.GetServer())

	if err != nil {
		if errors.Is(err, usecase.ErrNoData) {
			return &api.BalancerGetResponse{
				Value: err.Error(),
			}, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, usecase.ErrUnknownServer) {
			return &api.BalancerGetResponse{
				Value: err.Error(),
			}, status.Error(codes.Unavailable, err.Error())
		}

		return &api.BalancerGetResponse{
			Value: err.Error(),
		}, err
	}

	return &api.BalancerGetResponse{
		Value: value,
	}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.ConnectRequest) (*api.ConnectResponse, error) {
	serverNum, err := h.logic.Connect(request.GetAddress(), request.GetAvailable(), request.GetTotal(), request.Server)
	if err != nil {
		return nil, err
	}

	return &api.ConnectResponse{
		Status:       "connected successfully",
		ServerNumber: serverNum,
	}, nil
}

func (h *Handler) Index(ctx context.Context, request *api.IndexRequest) (*api.IndexResponse, error) {
	_, err := h.logic.Index(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.IndexResponse{}, nil
}

func (h *Handler) Disconnect(ctx context.Context, request *api.DisconnectRequest) (*api.DisconnectResponse, error) {
	h.logic.Disconnect(request.GetServerNumber())

	return &api.DisconnectResponse{}, nil
}

func (h *Handler) Servers(ctx context.Context, request *api.ServersRequest) (*api.ServersResponse, error) {
	servers := h.logic.Servers()
	s := strings.Join(servers, "<br>")
	return &api.ServersResponse{
		ServersInfo: s,
	}, nil
}
