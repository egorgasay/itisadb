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

func (h *Handler) SetToObject(ctx context.Context, r *api.SetToObjectRequest) (*api.SetResponse, error) {
	ram, err := h.logic.SetToObject(r.Name, r.Key, r.Value, r.Unique)
	resp := &api.SetResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, storage.ErrObjectNotFound) {
			return resp, status.Error(codes.ResourceExhausted, err.Error())
		}

		return resp, err
	}

	return resp, nil
}

func (h *Handler) GetFromObject(ctx context.Context, r *api.GetFromObjectRequest) (*api.GetResponse, error) {
	ram, value, err := h.logic.GetFromObject(r.Name, r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, storage.ErrObjectNotFound) {
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

func (h *Handler) ObjectToJSON(ctx context.Context, r *api.ObjectToJSONRequest) (*api.ObjectToJSONResponse, error) {
	ram, object, err := h.logic.ObjectToJSON(r.Name)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
			return &api.ObjectToJSONResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.ResourceExhausted, err.Error())
		}
		return &api.ObjectToJSONResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, err
	}
	return &api.ObjectToJSONResponse{
		Ram:    &api.Ram{Total: ram.Total, Available: ram.Available},
		Object: object,
	}, nil
}

func (h *Handler) NewObject(ctx context.Context, r *api.NewObjectRequest) (*api.NewObjectResponse, error) {
	ram, err := h.logic.NewObject(r.Name)
	resp := &api.NewObjectResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrEmptyObjectName) {
			return resp, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, storage.ErrSomethingExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}

		return resp, err
	}
	return resp, nil
}

func (h *Handler) AttachToObject(ctx context.Context, r *api.AttachToObjectRequest) (*api.AttachToObjectResponse, error) {
	ram, err := h.logic.AttachToObject(r.Dst, r.Src)
	resp := &api.AttachToObjectResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
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

func (h *Handler) DeleteObject(ctx context.Context, r *api.DeleteObjectRequest) (*api.DeleteObjectResponse, error) {
	ram, err := h.logic.DeleteObject(r.Object)
	resp := &api.DeleteObjectResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
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

func (h *Handler) Size(ctx context.Context, r *api.ObjectSizeRequest) (*api.ObjectSizeResponse, error) {
	ram, size, err := h.logic.Size(r.Name)
	resp := &api.ObjectSizeResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
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

		if errors.Is(err, storage.ErrObjectNotFound) {
			return resp, status.Error(codes.ResourceExhausted, storage.ErrObjectNotFound.Error())
		}

		return resp, err
	}
	return resp, nil
}
