package converterr

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotFound = errors.New("not found")
var ErrUnavailable = errors.New("service unavailable")
var ErrInvalidName = errors.New("invalid object name")
var ErrObjectNotFound = errors.New("object not found")
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

func GetFromObject(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func DelFromObject(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func SetToObject(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return ErrExists
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func ObjectToJSON(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func Object(err error) error {
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

func DelObject(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func SizeObject(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func IsObject(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}

func AttachObject(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.PermissionDenied:
		return ErrCircularAttachment
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.Unavailable:
		return ErrUnavailable
	}

	return err
}
