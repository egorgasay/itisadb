package models

import "itisadb/internal/constants"

type GetOptions struct {
	Server *int32
}

type SetOptions struct {
	Server   *int32
	ReadOnly bool
}

type DeleteOptions struct {
	Server *int32
}

type Level byte

func (l Level) Higher() Level {
	return max(l+1, constants.MinLevel)
}

func (l Level) Lower() Level {
	return min(l-1, constants.MaxLevel)
}

func (l Level) IsHighest() bool {
	return l == constants.MaxLevel
}

func (l Level) IsLowest() bool {
	return l == constants.MinLevel
}

type ObjectOptions struct {
	Server *int32
	Level  Level
}

type ObjectToJSONOptions struct {
	Server *int32
}

type DeleteObjectOptions struct {
	Server *int32
}

type IsObjectOptions struct {
	Server *int32
}

type SizeOptions struct {
	Server *int32
}

type AttachToObjectOptions struct {
	Server *int32
}

type SetToObjectOptions struct {
	Server   *int32
	Uniques  bool
	ReadOnly bool
}

type GetFromObjectOptions struct {
	Server *int32
}

type ConnectOptions struct {
	Server *int32
}

type DeleteAttrOptions struct {
	Server *int32
}

type CreateUserOptions struct {
	Level Level
}
