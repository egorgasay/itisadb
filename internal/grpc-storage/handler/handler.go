package handler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/grpc-storage/storage"
	"itisadb/internal/grpc-storage/usecase"
	api "itisadb/pkg/api/storage"
)

type Handler struct {
	api.UnimplementedStorageServer
	logic usecase.IUseCase
}

func New(logic usecase.IUseCase) *Handler {
	return &Handler{logic: logic}
}

func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	ram, err := h.logic.Set(r.Key, r.Value, r.Unique)
	resp := &api.SetResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if errors.Is(err, storage.ErrAlreadyExists) {
		return resp, status.Error(codes.AlreadyExists, err.Error())
	}

	return resp, err
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	ram, value, err := h.logic.Get(r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.NotFound, err.Error())
		}
		return &api.GetResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, err
	}

	return &api.GetResponse{
		Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
		Value: value,
	}, nil
}

func (h *Handler) SetToIndex(ctx context.Context, r *api.SetToIndexRequest) (*api.SetResponse, error) {
	ram, err := h.logic.SetToIndex(r.Name, r.Key, r.Value, r.Unique)
	resp := &api.SetResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, storage.ErrIndexNotFound) {
			return resp, status.Error(codes.ResourceExhausted, err.Error())
		}

		return resp, err
	}

	return resp, nil
}

func (h *Handler) GetFromIndex(ctx context.Context, r *api.GetFromIndexRequest) (*api.GetResponse, error) {
	ram, value, err := h.logic.GetFromIndex(r.Name, r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, storage.ErrIndexNotFound) {
			return &api.GetResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.ResourceExhausted, err.Error())
		}

		return &api.GetResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, err
	}

	return &api.GetResponse{
		Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
		Value: value,
	}, nil
}

func (h *Handler) IndexToJSON(ctx context.Context, r *api.IndexToJSONRequest) (*api.IndexToJSONResponse, error) {
	ram, index, err := h.logic.IndexToJSON(r.Name)
	if err != nil {
		if errors.Is(err, storage.ErrIndexNotFound) {
			return &api.IndexToJSONResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.ResourceExhausted, err.Error())
		}
		return &api.IndexToJSONResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, err
	}
	return &api.IndexToJSONResponse{
		Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
		Index: index,
	}, nil
}

func (h *Handler) NewIndex(ctx context.Context, r *api.NewIndexRequest) (*api.NewIndexResponse, error) {
	ram, err := h.logic.NewIndex(r.Name)
	resp := &api.NewIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrEmptyIndexName) {
			return resp, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, storage.ErrSomethingExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}

		return resp, err
	}
	return resp, nil
}

func (h *Handler) AttachToIndex(ctx context.Context, r *api.AttachToIndexRequest) (*api.AttachToIndexResponse, error) {
	ram, err := h.logic.AttachToIndex(r.Dst, r.Src)
	resp := &api.AttachToIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrIndexNotFound) {
			return resp, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, storage.ErrCircularAttachment) {
			return resp, status.Error(codes.PermissionDenied, err.Error())
		}

		if errors.Is(err, storage.ErrSomethingExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}

		return resp, err
	}
	return resp, nil
}

func (h *Handler) DeleteIndex(ctx context.Context, r *api.DeleteIndexRequest) (*api.DeleteIndexResponse, error) {
	ram, err := h.logic.DeleteIndex(r.Index)
	resp := &api.DeleteIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrIndexNotFound) {
			return resp, status.Error(codes.ResourceExhausted, err.Error())
		}
		return resp, err
	}
	return resp, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	ram, err := h.logic.Delete(r.Key)
	resp := &api.DeleteResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return resp, status.Error(codes.NotFound, storage.ErrNotFound.Error())
		}
		return resp, err
	}
	return resp, nil
}

func (h *Handler) Size(ctx context.Context, r *api.IndexSizeRequest) (*api.IndexSizeResponse, error) {
	ram, size, err := h.logic.Size(r.Name)
	resp := &api.IndexSizeResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrIndexNotFound) {
			return resp, status.Error(codes.ResourceExhausted, err.Error())
		}

		if errors.Is(err, storage.ErrSomethingExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}

		return resp, err
	}
	resp.Size = size
	return resp, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.DeleteAttrRequest) (*api.DeleteAttrResponse, error) {
	ram, err := h.logic.DeleteAttr(r.Name, r.Key)
	resp := &api.DeleteAttrResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return resp, status.Error(codes.NotFound, storage.ErrNotFound.Error())
		}

		if errors.Is(err, storage.ErrIndexNotFound) {
			return resp, status.Error(codes.ResourceExhausted, storage.ErrIndexNotFound.Error())
		}

		return resp, err
	}
	return resp, nil
}
