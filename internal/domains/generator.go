package domains

import (
	"context"
	"time"
)

type Generator interface {
	AccessToken(ctx context.Context, guid int, key []byte, accessTTL time.Duration) (access string, exp int64, err error)
	RefreshToken(ctx context.Context, refreshTTL time.Duration) (token string, exp int64, err error)
}
