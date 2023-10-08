package grpc

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"itisadb/internal/domains"
	"itisadb/internal/handler/converterr"
	"itisadb/internal/models"
	"itisadb/pkg/api"
	"strings"
)

type Handler struct {
	api.UnimplementedItisaDBServer
	core   domains.Core
	logger *zap.Logger
}

func New(logic domains.Core, l *zap.Logger) *Handler {
	return &Handler{core: logic, logger: l}
}

func safeDeref[T any](x *T) T {
	if x == nil {
		return *new(T)
	}
	return *x
}

func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	opts := models.SetOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy SetOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	setTo, err := h.core.Set(ctx, r.Key, r.Value, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.SetResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) SetToObject(ctx context.Context, r *api.SetToObjectRequest) (*api.SetToObjectResponse, error) {
	opts := models.SetToObjectOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy SetToObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	setTo, err := h.core.SetToObject(ctx, r.Object, r.Key, r.Value, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.SetToObjectResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	opts := models.GetOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy GetOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	value, err := h.core.Get(ctx, r.Key, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.GetResponse{
		Value: value,
	}, nil
}

func (h *Handler) GetFromObject(ctx context.Context, r *api.GetFromObjectRequest) (*api.GetFromObjectResponse, error) {
	opts := models.GetFromObjectOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy GetFromObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	value, err := h.core.GetFromObject(ctx, r.Object, r.Key, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.GetFromObjectResponse{
		Value: value,
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	opts := models.DeleteOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy DeleteOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.Delete(ctx, r.Key, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteResponse{}, nil
}

func (h *Handler) AttachToObject(ctx context.Context, r *api.AttachToObjectRequest) (*api.AttachToObjectResponse, error) {
	opts := models.AttachToObjectOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy AttachToObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.AttachToObject(ctx, r.Dst, r.Src, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.AttachToObjectResponse{}, nil
}

func (h *Handler) DeleteObject(ctx context.Context, r *api.DeleteObjectRequest) (*api.DeleteObjectResponse, error) {
	opts := models.DeleteObjectOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy DeleteObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.DeleteObject(ctx, r.Object, opts)
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
	opts := models.ObjectOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy ObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	_, err = h.core.Object(ctx, r.Name, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ObjectResponse{}, nil
}

func (h *Handler) ObjectToJSON(ctx context.Context, r *api.ObjectToJSONRequest) (*api.ObjectToJSONResponse, error) {
	opts := models.ObjectToJSONOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy ObjectToJSONOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	m, err := h.core.ObjectToJSON(ctx, r.Name, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ObjectToJSONResponse{
		Object: m,
	}, nil
}

func (h *Handler) IsObject(ctx context.Context, r *api.IsObjectRequest) (*api.IsObjectResponse, error) {
	opts := models.IsObjectOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy IsObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	ok, err := h.core.IsObject(ctx, r.Name, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.IsObjectResponse{
		Ok: ok,
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.DeleteAttrRequest) (*api.DeleteAttrResponse, error) {
	opts := models.DeleteAttrOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy DeleteAttrOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.DeleteAttr(ctx, r.Key, r.Object, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteAttrResponse{}, nil
}

func (h *Handler) Size(ctx context.Context, r *api.ObjectSizeRequest) (*api.ObjectSizeResponse, error) {
	opts := models.SizeOptions{}

	err := copier.Copy(&opts, safeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy SizeOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	size, err := h.core.Size(ctx, r.Name, opts)
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

func (h *Handler) CreateUser(ctx context.Context, r *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	user := models.User{}

	err := copier.Copy(&user, safeDeref(r.User))
	if err != nil {
		h.logger.Warn("failed to copy User", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.CreateUser(ctx, user)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.CreateUserResponse{}, nil
}
func (h *Handler) DeleteUser(ctx context.Context, r *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	err := h.core.DeleteUser(ctx, r.Login)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteUserResponse{}, nil
}
func (h *Handler) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	err := h.core.ChangePassword(ctx, r.Login, r.NewPassword)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ChangePasswordResponse{}, nil
}
func (h *Handler) ChangeLevel(ctx context.Context, r *api.ChangeLevelRequest) (*api.ChangeLevelResponse, error) {
	err := h.core.ChangeLevel(ctx, r.Login, int8(r.Level))
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ChangeLevelResponse{}, nil
}
