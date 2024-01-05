package storage

import (
	"encoding/json"

	"github.com/egorgasay/gost"
)

type value struct {
	value    string
	readOnly bool
}

func (v value) MarshalJSON() ([]byte, error) {
	var data = make(map[string]interface{})

	data["value"] = v.value
	data["read_only"] = v.readOnly

	return json.Marshal(data)
}

func NewValue(val string, readOnly bool) *value {
	return &value{value: val, readOnly: readOnly}
}

func (v *value) Object() (opt gost.Option[*object]) {
	return opt.None()
}

func (v *value) IsObject() bool {
	return false
}

func (v *value) Value() gost.Option[string] {
	return gost.Some(v.value)
}

func (v *value) IsValue() bool {
	return true
}
