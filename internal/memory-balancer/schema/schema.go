package schema

type SetRequest struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Server int32  `json:"server,omitempty"`
}

type GetRequest struct {
	Key    string `json:"key"`
	Server int32  `json:"server,omitempty"`
}
