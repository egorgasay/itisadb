package models

import "github.com/egorgasay/itisadb-go-sdk"

type GetOptions struct {
	Server *int32
}

func (o GetOptions) ToSDK() itisadb.GetOptions {
	return itisadb.GetOptions{
		Server: o.Server,
	}
}

type SetOptions struct {
	Server   *int32
	ReadOnly bool
	Unique   bool
	Level    *int8 // TODO: handle?
}

func (o SetOptions) ToSDK() itisadb.SetOptions {
	return itisadb.SetOptions{
		Server:   o.Server,
		ReadOnly: o.ReadOnly,
	}
}

type DeleteOptions struct {
	Server *int32
}

func (o DeleteOptions) ToSDK() itisadb.DeleteOptions {
	return itisadb.DeleteOptions{
		Server: o.Server,
	}
}

type Level byte

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
	Server *int32
	Level  Level
}

func (o ObjectOptions) ToSDK() itisadb.ObjectOptions {
	return itisadb.ObjectOptions{
		Server: o.Server,
		Level:  itisadb.Level(o.Level),
	}
}

type ObjectToJSONOptions struct {
	Server *int32
}

func (o ObjectToJSONOptions) ToSDK() itisadb.ObjectToJSONOptions {
	return itisadb.ObjectToJSONOptions{
		Server: o.Server,
	}
}

type DeleteObjectOptions struct {
	Server *int32
}

func (o DeleteObjectOptions) ToSDK() itisadb.DeleteObjectOptions {
	return itisadb.DeleteObjectOptions{
		Server: o.Server,
	}
}

type IsObjectOptions struct {
	Server *int32
}

func (o IsObjectOptions) ToSDK() itisadb.IsObjectOptions {
	return itisadb.IsObjectOptions{
		Server: o.Server,
	}
}

type SizeOptions struct {
	Server *int32
}

func (o SizeOptions) ToSDK() itisadb.SizeOptions {
	return itisadb.SizeOptions{
		Server: o.Server,
	}
}

type AttachToObjectOptions struct {
	Server *int32
}

func (o AttachToObjectOptions) ToSDK() itisadb.AttachToObjectOptions {
	return itisadb.AttachToObjectOptions{
		Server: o.Server,
	}
}

type SetToObjectOptions struct {
	Server   *int32
	ReadOnly bool
}

func (o SetToObjectOptions) ToSDK() itisadb.SetToObjectOptions {
	return itisadb.SetToObjectOptions{
		Server:   o.Server,
		ReadOnly: o.ReadOnly,
	}
}

type GetFromObjectOptions struct {
	Server *int32
	Level  *Level //TODO:???
}

func (o GetFromObjectOptions) ToSDK() itisadb.GetFromObjectOptions {
	return itisadb.GetFromObjectOptions{
		Server: o.Server,
	}
}

type ConnectOptions struct {
	Server *int32
}

func (o ConnectOptions) ToSDK() itisadb.ConnectOptions {
	return itisadb.ConnectOptions{
		Server: o.Server,
	}
}

type DeleteAttrOptions struct {
	Server *int32
}

//func (o DeleteAttrOptions) ToSDK() itisadb.De {
//	return itisadb.DeleteAttrOptions{
//		Server: o.Server,
//	}
//}

type CreateUserOptions struct {
	Level Level
}

func (o CreateUserOptions) ToSDK() itisadb.CreateUserOptions {
	return itisadb.CreateUserOptions{
		Level: itisadb.Level(o.Level),
	}
}
