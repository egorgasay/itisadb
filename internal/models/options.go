package models

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
