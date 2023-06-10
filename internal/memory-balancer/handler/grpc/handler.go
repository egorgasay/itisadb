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
	logic iUseCase
}

type iUseCase interface {
	Index(ctx context.Context, name string) (int32, error)
	GetFromIndex(ctx context.Context, index string, key string, serverNumber int32) (string, error)
	SetToIndex(ctx context.Context, index string, key string, val string, uniques bool) (int32, error)
	GetIndex(ctx context.Context, name string) (map[string]string, error)
	IsIndex(ctx context.Context, name string) (bool, error)
	Size(ctx context.Context, name string) (uint64, error)
	DeleteIndex(ctx context.Context, name string) error
	AttachToIndex(ctx context.Context, dst string, src string) error
	DeleteAttr(ctx context.Context, attr string, index string) error
	Set(ctx context.Context, key string, val string, serverNumber int32, uniques bool) (int32, error)
	Get(ctx context.Context, key string, serverNumber int32) (string, error)
	Connect(address string, available uint64, total uint64, server int32) (int32, error)
	Disconnect(ctx context.Context, number int32) error
	Servers() []string
	Delete(ctx context.Context, key string, num int32) error
}

func New(logic iUseCase) *Handler {
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

func (h *Handler) SetToIndex(ctx context.Context, r *api.BalancerSetToIndexRequest) (*api.BalancerSetToIndexResponse, error) {
	setTo, err := h.logic.SetToIndex(ctx, r.Index, r.Key, r.Value, r.Uniques)
	if err != nil {
		if errors.Is(err, servers.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "")
		}
		return nil, err
	}

	return &api.BalancerSetToIndexResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.BalancerGetRequest) (*api.BalancerGetResponse, error) {
	value, err := h.logic.Get(ctx, r.Key, r.Server)
	if err != nil {
		if errors.Is(err, usecase.ErrNoData) {
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

func (h *Handler) GetFromIndex(ctx context.Context, r *api.BalancerGetFromIndexRequest) (*api.BalancerGetFromIndexResponse, error) {
	value, err := h.logic.GetFromIndex(ctx, r.GetIndex(), r.GetKey(), r.GetServer())
	if err != nil {
		if errors.Is(err, usecase.ErrNoData) {
			return &api.BalancerGetFromIndexResponse{
				Value: err.Error(),
			}, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, usecase.ErrUnknownServer) {
			return &api.BalancerGetFromIndexResponse{
				Value: err.Error(),
			}, status.Error(codes.Unavailable, err.Error())
		}

		return &api.BalancerGetFromIndexResponse{
			Value: err.Error(),
		}, err
	}

	return &api.BalancerGetFromIndexResponse{
		Value: value,
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.BalancerDeleteRequest) (*api.BalancerDeleteResponse, error) {
	err := h.logic.Delete(ctx, r.Key, r.Server)
	if err != nil {
		if errors.Is(err, usecase.ErrNoData) {
			return &api.BalancerDeleteResponse{}, status.Error(codes.NotFound, err.Error())
		}
		return &api.BalancerDeleteResponse{}, err
	}
	return &api.BalancerDeleteResponse{}, nil
}

func (h *Handler) AttachToIndex(ctx context.Context, r *api.BalancerAttachToIndexRequest) (*api.BalancerAttachToIndexResponse, error) {
	err := h.logic.AttachToIndex(ctx, r.Dst, r.Src)
	if err != nil {
		// TODO: handle error
		return &api.BalancerAttachToIndexResponse{}, err
	}
	return &api.BalancerAttachToIndexResponse{}, nil
}

func (h *Handler) DeleteIndex(ctx context.Context, r *api.BalancerDeleteIndexRequest) (*api.BalancerDeleteIndexResponse, error) {
	err := h.logic.DeleteIndex(ctx, r.Index)
	if err != nil {
		// TODO: handle error
		return &api.BalancerDeleteIndexResponse{}, err
	}
	return &api.BalancerDeleteIndexResponse{}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.BalancerConnectRequest) (*api.BalancerConnectResponse, error) {
	serverNum, err := h.logic.Connect(request.GetAddress(), request.GetAvailable(), request.GetTotal(), request.Server)
	if err != nil {
		return nil, err
	}

	return &api.BalancerConnectResponse{
		Status:       "connected successfully",
		ServerNumber: serverNum,
	}, nil
}

func (h *Handler) Index(ctx context.Context, request *api.BalancerIndexRequest) (*api.BalancerIndexResponse, error) {
	_, err := h.logic.Index(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerIndexResponse{}, nil
}

func (h *Handler) GetIndex(ctx context.Context, request *api.BalancerGetIndexRequest) (*api.BalancerGetIndexResponse, error) {
	m, err := h.logic.GetIndex(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerGetIndexResponse{
		Index: m,
	}, nil
}

func (h *Handler) IsIndex(ctx context.Context, request *api.BalancerIsIndexRequest) (*api.BalancerIsIndexResponse, error) {
	ok, err := h.logic.IsIndex(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerIsIndexResponse{
		Ok: ok,
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.BalancerDeleteAttrRequest) (*api.BalancerDeleteAttrResponse, error) {
	err := h.logic.DeleteAttr(ctx, r.GetKey(), r.GetIndex())
	if err != nil {
		if errors.Is(err, servers.ErrNotFound) {
			return &api.BalancerDeleteAttrResponse{}, status.Error(codes.ResourceExhausted, "")
		}
		if errors.Is(err, servers.ErrUnavailable) {
			return &api.BalancerDeleteAttrResponse{}, status.Error(codes.Unavailable, "")
		}
		if errors.Is(err, usecase.ErrIndexNotFound) {
			return &api.BalancerDeleteAttrResponse{}, status.Error(codes.NotFound, "")
		}
		return &api.BalancerDeleteAttrResponse{}, err
	}
	return &api.BalancerDeleteAttrResponse{}, nil
}

func (h *Handler) Size(ctx context.Context, request *api.BalancerIndexSizeRequest) (*api.BalancerIndexSizeResponse, error) {
	size, err := h.logic.Size(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.BalancerIndexSizeResponse{
		Size: size,
	}, nil
}

func (h *Handler) Disconnect(ctx context.Context, request *api.BalancerDisconnectRequest) (*api.BalancerDisconnectResponse, error) {
	h.logic.Disconnect(ctx, request.GetServerNumber())

	return &api.BalancerDisconnectResponse{}, nil
}

func (h *Handler) Servers(ctx context.Context, request *api.BalancerServersRequest) (*api.BalancerServersResponse, error) {
	servers := h.logic.Servers()
	s := strings.Join(servers, "<br>")
	return &api.BalancerServersResponse{
		ServersInfo: s,
	}, nil
}
