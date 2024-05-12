package constants

import (
	"errors"
	"fmt"

	"github.com/egorgasay/gost"
)

var (
	ErrServerNotFound   = gost.NewErrX(0, "server not found")
	ErrWrongCredentials = errors.New("wrong credentials")

	ErrNoData        = errors.New("the value is not found")
	ErrUnknownServer = errors.New("unknown server")

	ErrAlreadyExists = gost.NewErrX(0, "already exists")
	ErrUnavailable   = errors.New("server is unavailable")

	ErrNotFound      = gost.NewErrX(0, "not found")
	ErrObjectNotFound   = ErrNotFound.Extend(0, "object not found")

	ErrCircularAttachment = gost.NewErrX(0, "circular attachment")
	ErrInternal           = gost.NewErrX(0, "internal error")
	ErrInvalidName        = fmt.Errorf("invalid name")

	ErrSomethingExists = ErrAlreadyExists.Extend(0, "something with this name already exists")
	ErrEmptyObjectName = gost.NewErrX(0, "object name is empty")

	/*
	 JWT Errors
	*/

	ErrInvalidGUID   = errors.New("invalid GUID")
	ErrSignToken     = errors.New("can't sign token")
	ErrInvalidToken  = errors.New("invalid token")
	ErrGenerateToken = errors.New("can't generate token")

	/*
		Session Errors
	*/

	ErrInvalidPassword = errors.New("invalid password")

	ErrForbidden = gost.NewErrX(0, "forbidden")
)
