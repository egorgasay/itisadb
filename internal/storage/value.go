package storage

import (
	"github.com/dolthub/swiss"
	"itisadb/internal/constants"
	"itisadb/pkg"
)

type object struct {
	value      string
	name       string
	values     *swiss.Map[string, *object]
	attachedTo []string
}

func NewObject(name string, attachedTo []string) *object {
	if attachedTo == nil {
		attachedTo = []string{name}
	}
	return &object{
		values:     swiss.NewMap[string, *object](10),
		attachedTo: attachedTo,
		name:       name,
	}
}

func (v *object) GetValue() (val string) {
	val = v.value
	return val
}

func (v *object) Name() string {
	return v.name
}

func (v *object) IsEmpty() bool {
	return v.values == nil
}

func (v *object) Get(name string) (string, error) {
	val, ok := v.values.Get(name)
	if !ok {
		return "", constants.ErrNotFound
	}

	return val.GetValue(), nil
}

func (v *object) Size() int {
	return v.values.Count()
}

func (v *object) IsObject() bool {
	return v.values != nil
}

func (v *object) IsAttached(name string) bool {
	for _, n := range pkg.Clone(v.attachedTo) {
		if n == name {
			return true
		}
	}
	return false
}

func (v *object) setAttached(attachedTo []string) {
	v.attachedTo = attachedTo
}

func (v *object) AttachObject(src *object) (err error) {
	defer func() {
		if err == nil {
			src.setAttached(v.attachedTo)
		}
	}()
	if v.values == nil {
		v.values = swiss.NewMap[string, *object](10)
		v.values.Put(src.Name(), src)
		return nil
	}
	if v.values.Has(src.Name()) {
		return nil
	}
	v.values.Put(src.Name(), src)
	return nil
}

func (v *object) Iter(f func(k string, v *object) bool) {
	v.values.Iter(f)
}

func (v *object) NextOrCreate(name string) *object {
	val, ok := v.values.Get(name)
	if !ok {
		blank := NewObject(name, v.attachedTo)
		v.values.Put(name, blank)
		return blank
	}

	return val
}

func (v *object) Object(name string) (*object, bool) {
	val, ok := v.values.Get(name)
	if !ok {
		return nil, false
	}

	return val, true
}

func (v *object) SetValueUnique(val string) error {
	if v.values.Has(val) {
		return constants.ErrAlreadyExists
	}

	v.value = val
	return nil
}

func (v *object) Delete(key string) error {
	if !v.values.Has(key) {
		return constants.ErrNotFound
	}

	v.values.Delete(key)
	return nil
}

func (v *object) RecreateObject() {
	v.values = swiss.NewMap[string, *object](10)
}

func (v *object) Set(key string, value string) {
	v.values.Put(key, &object{
		value: value,
		name:  key,
	})
}

func (v *object) Has(key string) bool {
	return v.values.Has(key)
}
