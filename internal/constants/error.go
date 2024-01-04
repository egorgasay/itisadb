package constants

import (
	"errors"
	"fmt"
	"github.com/egorgasay/gost"
)

var (
	ErrObjectNotFound   = fmt.Errorf("object not found")
	ErrServerNotFound   = errors.New("server not found")
	ErrWrongCredentials = errors.New("wrong credentials")

	ErrNoData        = errors.New("the value is not found")
	ErrUnknownServer = errors.New("unknown server")

	ErrAlreadyExists = gost.NewError(0, 0, "already exists")
	ErrUnavailable   = errors.New("server is unavailable")
	ErrNotFound      = gost.NewError(0, 0, "not found")

	ErrCircularAttachment = fmt.Errorf("circular attachment")
	ErrInternal           = fmt.Errorf("internal error")
	ErrInvalidName        = fmt.Errorf("invalid name")

	ErrSomethingExists = errors.New("something with this name already exists")
	ErrEmptyObjectName = errors.New("object name is empty")

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

	ErrForbidden = gost.NewError(0, 0, "forbidden")
)
