package schema

type SetRequest struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Server  int32  `json:"server,omitempty"`
	Uniques bool   `json:"uniques"`
}

type GetRequest struct {
	Key    string `json:"key"`
	Server int32  `json:"server,omitempty"`
}

type DelRequest struct {
	Key    string `json:"key"`
	Server int32  `json:"server,omitempty"`
}

type GetFromIndexRequest struct {
	Key    string `json:"key"`
	Index  string `json:"index"`
	Server int32  `json:"server,omitempty"`
}

type SetToIndexRequest struct {
	Key     string `json:"key"`
	Index   string `json:"index"`
	Value   string `json:"value"`
	Server  int32  `json:"server,omitempty"`
	Uniques bool   `json:"uniques"`
}

type DelFromIndexRequest struct {
	Key    string `json:"key"`
	Index  string `json:"index"`
	Server int32  `json:"server,omitempty"`
}

type IndexToJSONRequest struct {
	Index string `json:"index"`
}

type DelIndexRequest struct {
	Index string `json:"index"`
}

type SizeIndexRequest struct {
	Index string `json:"index"`
}

type IsIndexRequest struct {
	Name string `json:"name"`
}

type AttachRequest struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
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
