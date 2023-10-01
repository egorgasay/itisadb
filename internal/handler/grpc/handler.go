package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/servers"
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
	setTo, ram, err := h.core.Set(ctx, r.Server, r.Key, r.Value, r.Uniques)
	if err != nil {
		return nil, err
	}

	return &api.SetResponse{
		SavedTo: setTo,
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) SetToObject(ctx context.Context, r *api.SetToObjectRequest) (*api.SetToObjectResponse, error) {
	setTo, ram, err := h.core.SetToObject(ctx, r.Server, r.Object, r.Key, r.Value, r.Uniques)
	if err != nil {
		return nil, err
	}

	return &api.SetToObjectResponse{
		SavedTo: setTo,
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	value, ram, err := h.core.Get(ctx, r.Server, r.Key)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, constants.ErrUnknownServer) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}

		return nil, err
	}

	return &api.GetResponse{
		Value: value,
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) GetFromObject(ctx context.Context, r *api.GetFromObjectRequest) (*api.GetFromObjectResponse, error) {
	value, ram, err := h.core.GetFromObject(ctx, r.Server, r.GetObject(), r.GetKey())
	if err != nil {
		if errors.Is(err, constants.ErrNoData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, constants.ErrUnknownServer) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}

		return nil, err
	}

	return &api.GetFromObjectResponse{
		Value: value,
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	ram, err := h.core.Delete(ctx, r.Server, r.Key)
	resp := &api.DeleteResponse{
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) AttachToObject(ctx context.Context, r *api.AttachToObjectRequest) (*api.AttachToObjectResponse, error) {
	ram, err := h.core.AttachToObject(ctx, r.Server, r.Dst, r.Src)
	if err != nil {
		if errors.Is(err, constants.ErrObjectNotFound) {
			return nil, status.Error(codes.ResourceExhausted, constants.ErrObjectNotFound.Error())
		}
		return nil, err
	}

	return &api.AttachToObjectResponse{
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) DeleteObject(ctx context.Context, r *api.DeleteObjectRequest) (*api.DeleteObjectResponse, error) {
	ram, err := h.core.DeleteObject(ctx, r.Server, r.Object)
	if err != nil {
		return nil, err
	}

	return &api.DeleteObjectResponse{
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.ConnectRequest) (*api.ConnectResponse, error) {
	serverNum, err := h.core.Connect(request.GetAddress(), request.GetAvailable(), request.GetTotal())
	if err != nil {
		if errors.Is(err, servers.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, err
	}

	return &api.ConnectResponse{
		Status: "connected successfully",
		Server: serverNum,
	}, nil
}

func (h *Handler) Object(ctx context.Context, r *api.ObjectRequest) (*api.ObjectResponse, error) {
	_, ram, err := h.core.Object(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, err
	}

	return &api.ObjectResponse{
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) ObjectToJSON(ctx context.Context, r *api.ObjectToJSONRequest) (*api.ObjectToJSONResponse, error) {
	m, ram, err := h.core.ObjectToJSON(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, err
	}

	return &api.ObjectToJSONResponse{
		Object: m,
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) IsObject(ctx context.Context, r *api.IsObjectRequest) (*api.IsObjectResponse, error) {
	ok, ram, err := h.core.IsObject(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, err
	}

	return &api.IsObjectResponse{
		Ok: ok,
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.DeleteAttrRequest) (*api.DeleteAttrResponse, error) {
	ram, err := h.core.DeleteAttr(ctx, r.Server, r.GetKey(), r.GetObject())
	if err != nil {
		if errors.Is(err, constants.ErrObjectNotFound) {
			return &api.DeleteAttrResponse{}, status.Error(codes.ResourceExhausted, constants.ErrObjectNotFound.Error())
		}

		return &api.DeleteAttrResponse{}, err
	}

	return &api.DeleteAttrResponse{
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) Size(ctx context.Context, r *api.ObjectSizeRequest) (*api.ObjectSizeResponse, error) {
	size, ram, err := h.core.Size(ctx, r.Server, r.GetName())
	if err != nil {
		return nil, err
	}

	return &api.ObjectSizeResponse{
		Size: size,
		Ram: &api.Ram{
			Available: ram.Available,
			Total:     ram.Total,
		},
	}, nil
}

func (h *Handler) Disconnect(ctx context.Context, r *api.DisconnectRequest) (*api.DisconnectResponse, error) {
	err := h.core.Disconnect(ctx, r.Server)
	if err != nil {
		if errors.Is(err, context.Canceled) { // TODO: add everywhere
			return nil, status.Error(codes.Canceled, context.Canceled.Error())
		}
		return nil, err
	}

	return &api.DisconnectResponse{}, nil
}

func (h *Handler) Servers(ctx context.Context, r *api.ServersRequest) (*api.ServersResponse, error) {
	t, err := getToken(ctx)
	if err != nil {
		return nil, err
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
		if errors.Is(err, constants.ErrWrongCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, err
	}

	return &api.AuthResponse{Token: token}, nil
}
