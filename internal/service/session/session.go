package session

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
)

type Session struct {
	storage      domains.Storage
	generator    domains.Generator
	logger       *zap.Logger
	validateUser bool

	key []byte
}

func New(config config.Config, storage domains.Storage, generator domains.Generator, l *zap.Logger) domains.Session {
	return Session{
		storage:      storage,
		generator:    generator,
		logger:       l,
		validateUser: config.Balancer.On,
		key:          []byte("CHANGE_ME"), // TODO: change me
	}
}

func (s Session) AuthByPassword(ctx context.Context, username, password string) (string, error) {
	r := s.storage.GetUserByName(username)
	if r.IsErr() {
		if errors.Is(r.Error(), constants.ErrNotFound) {
			return "", constants.ErrWrongCredentials
		}
		return "", r.Error()
	}

	user := r.Unwrap()
	if user.Password != password {
		return "", constants.ErrInvalidPassword
	}

	token, _, err := s.generator.AccessToken(ctx, user.ExtractClaims(), s.key, constants.AccessTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Session) AuthByToken(ctx context.Context, token string) (models.UserClaims, error) {
	if ctx.Err() != nil {
		return models.UserClaims{}, ctx.Err()
	}

	claims, err := s.infoFromJWT(token)
	if err != nil {
		return models.UserClaims{}, err
	}

	if r := s.storage.GetUserByID(claims.ID); r.IsErr() {
		return models.UserClaims{}, r.Error()
	}

	return claims, nil
}

func (s Session) Create(ctx context.Context, userID int, level models.Level) (string, error) {
	token, _, err := s.generator.AccessToken(ctx, models.UserClaims{
		ID:    userID,
		Level: level,
	}, s.key, constants.AccessTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Session) infoFromJWT(token string) (models.UserClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		s.logger.Error("can't parse token", zap.Error(err))
		return models.UserClaims{}, constants.ErrInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("can't extract claims from token", zap.Error(err))
		return models.UserClaims{}, constants.ErrInvalidToken
	}

	guid, ok := claims[constants.GUID]
	if !ok {
		s.logger.Error("can't extract guid from token", zap.Error(err))
		return models.UserClaims{}, constants.ErrInvalidToken
	}

	guidInt, ok := guid.(float64)
	if !ok {
		s.logger.Error("can't convert guid to float64", zap.Error(err))
		return models.UserClaims{}, constants.ErrInvalidToken
	}

	levelRaw, ok := claims[constants.LEVEL]
	if !ok {
		s.logger.Error("can't extract level from token", zap.Error(err))
		return models.UserClaims{}, constants.ErrInvalidToken
	}

	level, ok := levelRaw.(float64)
	if !ok {
		s.logger.Error("can't convert meta to models.UserMeta", zap.Error(err))
		return models.UserClaims{}, constants.ErrInvalidToken
	}

	return models.UserClaims{
		ID:    int(guidInt),
		Level: models.Level(level),
	}, nil
}
