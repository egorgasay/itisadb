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

type FromIndexRequest struct {
	Key    string `json:"key"`
	Server int32  `json:"server,omitempty"`
}
