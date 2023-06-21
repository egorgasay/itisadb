package converterr

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotFound = errors.New("not found")
var ErrUnavailable = errors.New("service unavailable")
var ErrInvalidName = errors.New("invalid index name")
var ErrIndexNotFound = errors.New("index not found")
var ErrExists = errors.New("already exists")
var ErrCircularAttachment = errors.New("circular attachment not allowed")

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

func Del(err error) error {
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

func Set(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return ErrExists
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func GetFromIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func DelFromIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func SetToIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return ErrExists
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func GetIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func Index(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return ErrInvalidName
	case codes.AlreadyExists:
		return ErrExists
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func DelIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func SizeIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func IsIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func AttachIndex(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.PermissionDenied:
		return ErrCircularAttachment
	case codes.ResourceExhausted:
		return ErrIndexNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}
