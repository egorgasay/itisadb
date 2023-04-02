package handler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-storage/internal/memory-balancer/usecase"
	api "grpc-storage/pkg/api/balancer"
	"strings"
)

type Handler struct {
	api.UnsafeBalancerServer
	logic *usecase.UseCase
}

func New(logic *usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}

func (h *Handler) Set(ctx context.Context, r *api.BalancerSetRequest) (*api.BalancerSetResponse, error) {
	setTo, err := h.logic.Set(r.Key, r.Value)

	if err != nil {
		return nil, err
	}

	return &api.BalancerSetResponse{
		Status:  "ok",
		SavedTo: int32(setTo),
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.BalancerGetRequest) (*api.BalancerGetResponse, error) {
	value, err := h.logic.Get(r.Key, r.Server)
	if errors.Is(err, usecase.ErrNoData) {
		return &api.BalancerGetResponse{
			Value: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return &api.BalancerGetResponse{
			Value: err.Error(),
		}, err
	}

	return &api.BalancerGetResponse{
		Value: value,
	}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.ConnectRequest) (*api.ConnectResponse, error) {
	serverNum, err := h.logic.Connect(request.GetAddress(), request.GetAvailable(), request.GetTotal())
	if err != nil {
		return nil, err
	}

	return &api.ConnectResponse{
		Status:       "connected successfully",
		ServerNumber: serverNum,
	}, nil
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
