package storage

import (
	"github.com/egorgasay/gost"
)

type value struct {
	value    string
	readOnly bool
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
