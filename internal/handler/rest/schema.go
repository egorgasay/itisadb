package rest

type SetRequest struct {
	Server  *int32 `json:"server,omitempty"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Uniques bool   `json:"uniques"`
}

type GetRequest struct {
	Key    string `json:"key"`
	Server *int32 `json:"server,omitempty"`
}

type DelRequest struct {
	Key    string `json:"key"`
	Server *int32 `json:"server"`
}

type GetFromObjectRequest struct {
	Key    string `json:"key"`
	Object string `json:"object"`
	Server *int32 `json:"server,omitempty"`
}

type SetToObjectRequest struct {
	Key     string `json:"key"`
	Object  string `json:"object"`
	Value   string `json:"value"`
	Server  *int32 `json:"server,omitempty"`
	Uniques bool   `json:"uniques"`
}

type DelFromObjectRequest struct {
	Key    string `json:"key"`
	Object string `json:"object"`
	Server *int32 `json:"server,omitempty"`
}

type ObjectToJSONRequest struct {
	Server *int32 `json:"server,omitempty"`
	Object string `json:"object"`
}

type DelObjectRequest struct {
	Server *int32 `json:"server,omitempty"`
	Object string `json:"object"`
}

type SizeObjectRequest struct {
	Server *int32 `json:"server,omitempty"`
	Object string `json:"object"`
}

type IsObjectRequest struct {
	Server *int32 `json:"server,omitempty"`
	Name   string `json:"name"`
}

type AttachRequest struct {
	Server *int32 `json:"server,omitempty"`
	Dst    string `json:"dst"`
	Src    string `json:"src"`
}

type ConnectRequest struct {
	Address   string `json:"address"`
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Server    int32  `json:"server"`
}

type DisconnectRequest struct {
	Server int32 `json:"server"`
}
