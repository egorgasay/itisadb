package storage

import (
	"github.com/dolthub/swiss"
	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/pkg"
)

type object struct {
	name       string
	values     *swiss.Map[string, something]
	attachedTo []string
}

func NewObject(name string, attachedTo []string) *object {
	if attachedTo == nil {
		attachedTo = []string{name}
	}
	return &object{
		values:     swiss.NewMap[string, something](10),
		attachedTo: attachedTo,
		name:       name,
	}
}

func (v *object) Object() gost.Option[*object] {
	return gost.Some(v)
}

func (v *object) Value() (opt gost.Option[string]) {
	return opt.None()
}

func (v *object) IsValue() bool {
	return false
}

func (v *object) IsObject() bool {
	return true
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

	switch some := val.Value(); some.IsSome() {
	case true:
		return some.Unwrap(), nil
	default:
		return "", constants.ErrSomethingExists
	}
}

func (v *object) Size() int {
	return v.values.Count()
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
	v.attachedTo = append(v.attachedTo, attachedTo...)
}

func (v *object) AttachObject(src *object) (err error) {
	defer func() {
		if err == nil {
			src.setAttached(v.attachedTo)
		}
	}()

	if v.values == nil {
		v.values = swiss.NewMap[string, something](10)
		v.values.Put(src.Name(), src)
		return nil
	}

	if v.values.Has(src.Name()) {
		return nil
	}

	v.values.Put(src.Name(), src)

	return nil
}

func (v *object) Iter(f func(k string, v something) bool) {
	v.values.Iter(f)
}

func (v *object) NextOrCreate(name string) something {
	val, ok := v.values.Get(name)
	if !ok {
		blank := NewObject(name, v.attachedTo)
		v.values.Put(name, blank)
		return blank
	}

	return val
}

func (v *object) GetValue(key string) (something, bool) {
	val, ok := v.values.Get(key)
	return val, ok
}

func (v *object) Delete(key string) error {
	if !v.values.Has(key) {
		return constants.ErrNotFound
	}

	v.values.Delete(key)
	return nil
}

func (v *object) RecreateObject() {
	v.values = swiss.NewMap[string, something](10)
}

func (v *object) Set(key string, val string) {
	v.values.Put(key, &value{
		value: val,
	})
}

func (v *object) Has(key string) bool {
	return v.values.Has(key)
}
