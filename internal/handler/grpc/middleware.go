package grpc

import (
	"context"
	"google.golang.org/grpc"
	"itisadb/internal/handler/converterr"
)

func (h *Handler) AuthMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t, err := getToken(ctx)
	if err != nil && info.FullMethod != "/api.ItisaDB/Authenticate" {
		return nil, converterr.ToGRPC(err)
	}

	_ = t
	// if h.core.ValidateToken

	resp, err := handler(ctx, req)

	// h.logger.Info("AuthMiddleware:", zap.String("method", info.FullMethod), zap.String("token", t))

	return resp, err
}
