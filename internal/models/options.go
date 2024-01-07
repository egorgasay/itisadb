package models

import (
	"github.com/egorgasay/itisadb-go-sdk"
)

type GetOptions struct {
	Server int32
}

func (o GetOptions) ToSDK() itisadb.GetOptions {
	return itisadb.GetOptions{}
}

type SetOptions struct {
	Server   int32
	ReadOnly bool
	Unique   bool
	Level    Level
}

func (o SetOptions) ToSDK() itisadb.SetOptions {
	return itisadb.SetOptions{
		ReadOnly: o.ReadOnly,
		Unique:   o.Unique,
		Level:    o.Level.ToSDK(),
	}
}

type DeleteOptions struct {
	Server int32
}

func (o DeleteOptions) ToSDK() itisadb.DeleteOptions {
	return itisadb.DeleteOptions{}
}

type Level byte

func (l Level) ToSDK() itisadb.Level {
	return itisadb.Level(l)
}

func (l Level) String() string {
	switch itisadb.Level(l) {
	case itisadb.DefaultLevel:
		return "Default"
	case itisadb.RestrictedLevel:
		return "Restricted"
	case itisadb.SecretLevel:
		return "Secret"
	}

	return "Unknown"
}

//func (l Level) Higher() Level {
//	return max(l+1, constants.MinLevel)
//}
//
//func (l Level) Lower() Level {
//	return min(l-1, constants.MaxLevel)
//}
//
//func (l Level) IsHighest() bool {
//	return l == constants.MaxLevel
//}
//
//func (l Level) IsLowest() bool {
//	return l == constants.MinLevel
//}

type ObjectOptions struct {
	Server int32
	Level  Level
}

func (o ObjectOptions) ToSDK() itisadb.ObjectOptions {
	return itisadb.ObjectOptions{
		Level: itisadb.Level(o.Level),
	}
}

type ObjectToJSONOptions struct {
	Server int32
}

func (o ObjectToJSONOptions) ToSDK() itisadb.ObjectToJSONOptions {
	return itisadb.ObjectToJSONOptions{}
}

type DeleteObjectOptions struct {
	Server int32
}

func (o DeleteObjectOptions) ToSDK() itisadb.DeleteObjectOptions {
	return itisadb.DeleteObjectOptions{}
}

type IsObjectOptions struct {
	Server int32
}

func (o IsObjectOptions) ToSDK() itisadb.IsObjectOptions {
	return itisadb.IsObjectOptions{}
}

type SizeOptions struct {
	Server int32
}

func (o SizeOptions) ToSDK() itisadb.SizeOptions {
	return itisadb.SizeOptions{}
}

type AttachToObjectOptions struct {
	Server int32
}

func (o AttachToObjectOptions) ToSDK() itisadb.AttachToObjectOptions {
	return itisadb.AttachToObjectOptions{}
}

type SetToObjectOptions struct {
	Server   int32
	ReadOnly bool
}

func (o SetToObjectOptions) ToSDK() itisadb.SetToObjectOptions {
	return itisadb.SetToObjectOptions{
		ReadOnly: o.ReadOnly,
	}
}

type GetFromObjectOptions struct {
	Server int32
}

func (o GetFromObjectOptions) ToSDK() itisadb.GetFromObjectOptions {
	return itisadb.GetFromObjectOptions{}
}

type ConnectOptions struct {
	Server int32
}

func (o ConnectOptions) ToSDK() itisadb.ConnectOptions {
	return itisadb.ConnectOptions{}
}

type DeleteAttrOptions struct {
	Server int32
}

func (o DeleteAttrOptions) ToSDK() itisadb.DeleteKeyOptions {
	return itisadb.DeleteKeyOptions{}
}

type CreateUserOptions struct {
	Level Level
}

func (o CreateUserOptions) ToSDK() itisadb.NewUserOptions {
	return itisadb.NewUserOptions{
		Level: itisadb.Level(o.Level),
	}
}
