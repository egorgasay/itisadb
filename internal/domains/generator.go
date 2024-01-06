package domains

import (
	"context"
	"time"

	"itisadb/internal/models"
)

type Generator interface {
	AccessToken(ctx context.Context, claims models.UserClaims, key []byte, accessTTL time.Duration) (access string, exp int64, err error)
	RefreshToken(ctx context.Context, refreshTTL time.Duration) (token string, exp int64, err error)
}
