package constants

import (
	"errors"
)

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrServerNotFound = errors.New("server not found")

	ErrNoServers        = errors.New("no servers available")
	ErrWrongCredentials = errors.New("wrong credentials")

	ErrNoData        = errors.New("the value is not found")
	ErrUnknownServer = errors.New("unknown server")

	ErrAlreadyExists = errors.New("already exists")
	ErrUnavailable   = errors.New("server is unavailable")
	ErrNotFound      = errors.New("not found")
)
