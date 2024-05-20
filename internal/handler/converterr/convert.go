package converterr

import (
	"context"
	"fmt"
	"itisadb/internal/constants"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConvertErr struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) ConvertErr {
	return ConvertErr{logger: logger}
}

func (c ConvertErr) ToGRPC(err error) error {
	c.logger.Info(err.Error())

	baseError, _ := Unwrap(err)
	switch baseError {
	case constants.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case constants.ErrObjectNotFound:
		return status.Error(codes.ResourceExhausted, err.Error())
	case constants.ErrUnavailable:
		return status.Error(codes.Unavailable, err.Error())
	case constants.ErrInvalidName:
		return status.Error(codes.InvalidArgument, err.Error())
	case constants.ErrAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case constants.ErrCircularAttachment:
		return status.Error(codes.FailedPrecondition, err.Error())
	case constants.ErrWrongCredentials:
		return status.Error(codes.Unauthenticated, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	case constants.ErrForbidden:
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return err
	}
}

func FromGRPC(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return constants.ErrNotFound
	case codes.ResourceExhausted:
		return constants.ErrObjectNotFound
	case codes.Unavailable:
		return constants.ErrUnavailable
	case codes.InvalidArgument:
		return constants.ErrInvalidName
	case codes.AlreadyExists:
		return constants.ErrAlreadyExists
	case codes.FailedPrecondition:
		return constants.ErrCircularAttachment
	case codes.Unauthenticated:
		return constants.ErrWrongCredentials
	default:
		return err
	}
}

func Unwrap(MultiError error) (baseError error, insideError error) {
	switch err := MultiError.(type) {
	case interface{ Unwrap() error }:
		return Unwrap(err.Unwrap())
	case interface{ Unwrap() []error }:
		errStack := append([]error(nil), err.Unwrap()...)
		if len(errStack) < 2 {
			return nil, fmt.Errorf("invalid error stack: %v", errStack)
		}
		return errStack[0], errStack[1]
	default:
		return MultiError, MultiError
	}
}
