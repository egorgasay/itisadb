package constants

import (
	"errors"
	"fmt"
)

var (
	ErrObjectNotFound   = fmt.Errorf("object not found")
	ErrServerNotFound   = errors.New("server not found")
	ErrWrongCredentials = errors.New("wrong credentials")

	ErrNoData        = errors.New("the value is not found")
	ErrUnknownServer = errors.New("unknown server")

	ErrAlreadyExists = errors.New("already exists")
	ErrUnavailable   = errors.New("server is unavailable")
	ErrNotFound      = errors.New("not found")

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

	ErrForbidden = errors.New("forbidden")
)
