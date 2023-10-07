package generator

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"time"
)

type Generator struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) domains.Generator {
	return &Generator{
		logger: logger,
	}
}

const (
	_guid = "guid"
	_iat  = "iat"
)

func (g *Generator) AccessToken(
	ctx context.Context,
	guid int, key []byte,
	accessTTL time.Duration,
) (access string, exp int64, err error) {

	if ctx.Err() != nil {
		return "", 0, ctx.Err()
	}

	exp = time.Now().Add(accessTTL).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			_guid: guid,
			_iat:  exp,
		})

	access, err = t.SignedString(key)
	if err != nil {
		g.logger.Error("can't sign token", zap.Error(err))
		return "", 0, constants.ErrSignToken
	}

	return access, exp, nil
}

func (g *Generator) RefreshToken(
	ctx context.Context,
	refreshTTL time.Duration,
) (token string, exp int64, err error) {
	if ctx.Err() != nil {
		return "", 0, ctx.Err()
	}

	uuidObj, err := uuid.NewUUID()
	if err != nil {
		g.logger.Error("can't generate uuid", zap.Error(err))
		return "", 0, constants.ErrGenerateToken
	}

	return uuidObj.String(), time.Now().Add(refreshTTL).Unix(), nil
}

// bcryptHashFrom generates a bcrypt hash from a given token.
func bcryptHashFrom(token []byte) ([]byte, error) {
	if len(token) == 0 {
		return nil, constants.ErrInvalidToken
	}

	bcryptHash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return bcryptHash, nil
}
