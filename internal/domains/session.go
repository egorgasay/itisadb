package domains

import (
	"context"

	"itisadb/internal/models"
)

type Session interface {
	AuthByToken(ctx context.Context, token string) (models.UserClaims, error)
	AuthByPassword(ctx context.Context, username, password string) (string, error)
	Create(ctx context.Context, userID int, level models.Level) (string, error)
}
