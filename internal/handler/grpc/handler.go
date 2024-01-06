package grpc

import (
	"context"
	"strings"

	"github.com/egorgasay/gost"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/handler/converterr"
	"itisadb/internal/models"

	api "github.com/egorgasay/itisadb-shared-proto/go"
)

type Handler struct {
	api.UnimplementedItisaDBServer
	core     domains.Balancer
	logger   *zap.Logger
	session  domains.Session
	security config.SecurityConfig
}

func New(
	logic domains.Balancer,
	l *zap.Logger,
	session domains.Session,
	conf config.SecurityConfig,
) *Handler {
	return &Handler{core: logic, logger: l, session: session, security: conf}
}

func (h *Handler) claimsFromContext(ctx context.Context) (opt gost.Option[models.UserClaims]) {
	value := ctx.Value(constants.UserKey)
	if value == nil {
		return opt.None()
	}

	claims, ok := value.(models.UserClaims)
	if !ok {
		h.logger.Warn("failed to cast userID", zap.Any("value", value))
		return opt.None()
	}

	return opt.Some(claims)
}

func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.SetOptions{}
	if err := copier.Copy(&opts, gost.SafeDeref(r.Options)); err != nil {
		h.logger.Warn("failed to copy SetOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	setTo, err := h.core.Set(ctx, claims, r.Key, r.Value, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.SetResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) SetToObject(ctx context.Context, r *api.SetToObjectRequest) (*api.SetToObjectResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.SetToObjectOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy SetToObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	setTo, err := h.core.SetToObject(ctx, claims, r.Object, r.Key, r.Value, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.SetToObjectResponse{
		SavedTo: setTo,
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.GetOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy GetOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	res, err := h.core.Get(ctx, claims, r.Key, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.GetResponse{
		Value:    res.Value,
		ReadOnly: res.ReadOnly,
		Level:    api.Level(res.Level),
	}, nil
}

func (h *Handler) GetFromObject(ctx context.Context, r *api.GetFromObjectRequest) (*api.GetFromObjectResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.GetFromObjectOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy GetFromObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	value, err := h.core.GetFromObject(ctx, claims, r.Object, r.Key, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.GetFromObjectResponse{
		Value: value,
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.DeleteOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy DeleteOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.Delete(ctx, claims, r.Key, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteResponse{}, nil
}

func (h *Handler) AttachToObject(ctx context.Context, r *api.AttachToObjectRequest) (*api.AttachToObjectResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.AttachToObjectOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy AttachToObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.AttachToObject(ctx, claims, r.Dst, r.Src, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.AttachToObjectResponse{}, nil
}

func (h *Handler) DeleteObject(ctx context.Context, r *api.DeleteObjectRequest) (*api.DeleteObjectResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.DeleteObjectOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy DeleteObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.DeleteObject(ctx, claims, r.Object, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteObjectResponse{}, nil
}

func (h *Handler) Connect(ctx context.Context, request *api.ConnectRequest) (*api.ConnectResponse, error) {
	serverNum, err := h.core.Connect(ctx, request.GetAddress(), request.GetAvailable(), request.GetTotal())
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ConnectResponse{
		Status: "connected successfully",
		Server: serverNum,
	}, nil
}

func (h *Handler) Object(ctx context.Context, r *api.ObjectRequest) (*api.ObjectResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.ObjectOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy ObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	_, err = h.core.Object(ctx, claims, r.Name, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ObjectResponse{}, nil
}

func (h *Handler) ObjectToJSON(ctx context.Context, r *api.ObjectToJSONRequest) (*api.ObjectToJSONResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.ObjectToJSONOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy ObjectToJSONOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	m, err := h.core.ObjectToJSON(ctx, claims, r.Name, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ObjectToJSONResponse{
		Object: m,
	}, nil
}

func (h *Handler) IsObject(ctx context.Context, r *api.IsObjectRequest) (*api.IsObjectResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.IsObjectOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy IsObjectOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	ok, err := h.core.IsObject(ctx, claims, r.Name, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.IsObjectResponse{
		Ok: ok,
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.DeleteAttrRequest) (*api.DeleteAttrResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.DeleteAttrOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy DeleteAttrOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.DeleteAttr(ctx, claims, r.Key, r.Object, opts)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteAttrResponse{}, nil
}

func (h *Handler) Size(ctx context.Context, r *api.ObjectSizeRequest) (*api.ObjectSizeResponse, error) {
	claims := h.claimsFromContext(ctx)

	opts := models.SizeOptions{}
	err := copier.Copy(&opts, gost.SafeDeref(r.Options))
	if err != nil {
		h.logger.Warn("failed to copy SizeOptions", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	size, err := h.core.Size(ctx, claims, r.Name, opts)
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
	claims := h.claimsFromContext(ctx)

	user := models.User{}
	err := copier.Copy(&user, gost.SafeDeref(r.User))
	if err != nil {
		h.logger.Warn("failed to copy User", zap.Error(err))
		return nil, converterr.ToGRPC(err)
	}

	err = h.core.CreateUser(ctx, claims, user)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.CreateUserResponse{}, nil
}
func (h *Handler) DeleteUser(ctx context.Context, r *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	claims := h.claimsFromContext(ctx)

	err := h.core.DeleteUser(ctx, claims, r.Login)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.DeleteUserResponse{}, nil
}
func (h *Handler) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	claims := h.claimsFromContext(ctx)

	err := h.core.ChangePassword(ctx, claims, r.Login, r.NewPassword)
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ChangePasswordResponse{}, nil
}
func (h *Handler) ChangeLevel(ctx context.Context, r *api.ChangeLevelRequest) (*api.ChangeLevelResponse, error) {
	claims := h.claimsFromContext(ctx)

	err := h.core.ChangeLevel(ctx, claims, r.Login, models.Level(r.Level))
	if err != nil {
		return nil, converterr.ToGRPC(err)
	}

	return &api.ChangeLevelResponse{}, nil
}

func (h *Handler) GetRam(ctx context.Context, r *api.GetRamRequest) (*api.GetRamResponse, error) {
	res := h.core.CalculateRAM(ctx)
	if res.IsErr() {
		return nil, converterr.ToGRPC(res.Error())
	}

	ram := res.Unwrap()

	return &api.GetRamResponse{Ram: &api.Ram{
		Total:     ram.Total,
		Available: ram.Available,
	}}, nil
}
