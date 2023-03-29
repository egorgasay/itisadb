package handler

import (
	"context"
	"github.com/egorgasay/grpc-storage/internal/memory-balancer/usecase"
	api "github.com/egorgasay/grpc-storage/pkg/api/balancer"
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
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.BalancerGetRequest) (*api.BalancerGetResponse, error) {
	value, err := h.logic.Get(r.Key)
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
	err := h.logic.Connect(request.GetAddress())
	if err != nil {
		return nil, err
	}

	return &api.ConnectResponse{
		Status: "connected successfully",
	}, nil
}
