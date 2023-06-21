package converterr

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotFound = errors.New("not found")
var ErrUnavailable = errors.New("service unavailable")
var ErrIndexNotFound = errors.New("index not found")
var ErrExists = errors.New("already exists")

func Get(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}
