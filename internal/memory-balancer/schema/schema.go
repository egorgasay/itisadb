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

type GetIndexRequest struct {
	Index string `json:"index"`
}
