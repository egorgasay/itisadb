package session

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
)

type Session struct {
	storage   domains.Storage
	generator domains.Generator
	logger    *zap.Logger

	key []byte
}

func New(storage domains.Storage, generator domains.Generator, l *zap.Logger) domains.Session {
	return Session{
		storage:   storage,
		generator: generator,
		logger:    l,
		key:       []byte("CHANGE_ME"), // TODO: change me
	}
}

func (s Session) AuthByPassword(ctx context.Context, username, password string) (string, error) {
	id, user, err := s.storage.GetUserByName(username)
	if err != nil {
		return "", err
	}

	if user.Password != password {
		return "", constants.ErrInvalidPassword
	}

	token, _, err := s.generator.AccessToken(ctx, id, s.key, constants.AccessTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Session) AuthByToken(ctx context.Context, token string) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	id, err := s.guidFromJWT(token)
	if err != nil {
		return 0, err
	}

	_, err = s.storage.GetUserByID(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Session) Create(ctx context.Context, guid int) (string, error) {
	token, _, err := s.generator.AccessToken(ctx, guid, s.key, constants.AccessTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Session) guidFromJWT(token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		s.logger.Error("can't parse token", zap.Error(err))
		return 0, constants.ErrInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("can't extract claims from token", zap.Error(err))
		return 0, constants.ErrInvalidToken
	}

	guid, ok := claims[constants.GUID]
	if !ok {
		s.logger.Error("can't extract guid from token", zap.Error(err))
		return 0, constants.ErrInvalidToken
	}

	guidInt, ok := guid.(float64)
	if !ok {
		s.logger.Error("can't convert guid to float64", zap.Error(err))
		return 0, constants.ErrInvalidToken
	}

	return int(guidInt), nil
}
