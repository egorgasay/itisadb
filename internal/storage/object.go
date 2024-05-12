package storage

import (
	"encoding/json"
	"sync"

	"github.com/dolthub/swiss"
	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/models"
	"itisadb/pkg"
)

type object struct {
	name       string
	values     *swiss.Map[string, Something]
	attachedTo []string
	level      models.Level
	*sync.RWMutex
}

func NewObject(name string, attachedTo []string, level models.Level) *object {
	if attachedTo == nil {
		attachedTo = []string{name}
	}

	return &object{
		values:     swiss.NewMap[string, Something](10),
		attachedTo: attachedTo,
		name:       name,
		level:      level,
		RWMutex:    &sync.RWMutex{},
	}
}

func (v *object) Object() gost.Option[*object] {
	return gost.Some(v)
}

func (v *object) Value() (opt gost.Option[value]) {
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

func (v *object) Get(name string) (r gost.Option[string]) {
	v.RLock()
	defer v.RUnlock()

	val, ok := v.values.Get(name)
	if !ok {
		return r.None()
	}

	switch some := val.Value(); some.IsSome() {
	case true:
		return r.Some(some.Unwrap().value)
	default:
		return r.None()
	}
}

func (v *object) Size() int {
	v.RLock()
	defer v.RUnlock()

	return v.values.Count()
}

func (v *object) IsAttached(name string) bool {
	cloned := func() []string {
		v.RLock()
		defer v.RUnlock()
		return pkg.Clone(v.attachedTo)
	}()

	for _, n := range cloned {
		if n == name {
			return true
		}
	}

	return false
}

func (v *object) setAttached(attachedTo []string) {
	v.Lock()
	defer v.Unlock()

	v.attachedTo = append(v.attachedTo, attachedTo...)
}

func (v *object) AttachObject(src *object) (r gost.ResultN) {
	v.RLock()
	defer func() {
		v.RUnlock()
		if r.IsOk() {
			src.setAttached(v.attachedTo)
		}
	}()

	if v.values == nil {
		v.values = swiss.NewMap[string, Something](10)
		v.values.Put(src.Name(), src)
		return r.Ok()
	}

	if v.values.Has(src.Name()) {
		return r.Ok()
	}

	v.values.Put(src.Name(), src)

	return r.Ok()
}

func (v *object) Iter(f func(k string, v Something) bool) {
	v.RLock()
	defer v.RUnlock()

	v.values.Iter(f)
}

func (v *object) NextOrCreate(name string, level models.Level) Something {
	v.RLock()
	defer v.RUnlock()

	val, ok := v.values.Get(name)
	if !ok {
		blank := NewObject(name, v.attachedTo, max(level, v.level))
		v.values.Put(name, blank)
		return blank
	}

	return val
}

func (v *object) GetValue(key string) (Something, bool) {
	v.RLock()
	defer v.RUnlock()

	val, ok := v.values.Get(key)
	return val, ok
}

func (v *object) Delete(key string) (r gost.ResultN) {
	v.Lock()
	defer v.Unlock()

	if !v.values.Has(key) {
		return r.Err(constants.ErrNotFound)
	}

	v.values.Delete(key)
	return r.Ok()
}

func (v *object) RecreateObject() {
	v.Lock()
	defer v.Unlock()

	v.values = swiss.NewMap[string, Something](10)
}

func (v *object) Set(key string, val string) {
	v.Lock()
	defer v.Unlock()

	v.values.Put(key, &value{
		value: val,
	})
}

func (v *object) Has(key string) bool {
	v.RLock()
	defer v.RUnlock()

	return v.values.Has(key)
}

//type ValueJSON struct {
//	Key string `json:"key"`
//	value
//}

func (v *object) MarshalJSON() (r gost.Result[[]byte]) {
	v.RLock()
	defer v.RUnlock()

	arr := make([]any, 0, 100)
	var data map[string]interface{}

	v.values.Iter(func(k string, v Something) bool {
		if v != nil {
			if v.IsValue() {
				val := v.Value().Unwrap()
				arr = append(arr, map[string]interface{}{
					"key":       k,
					"value":     val.value,
					"read_only": val.readOnly,
				})
			} else {
				arr = append(arr, v)
			}
		}

		return false
	})

	data = map[string]interface{}{
		"name":        v.name,
		"level":       v.level.String(),
		"attached_to": v.attachedTo,
		"values":      arr,
	}

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return r.Err(constants.ErrInternal.Extend(0, err.Error()))
	}

	return r.Ok(b)
}

func (v *object) setLevel(level models.Level) {
	v.Lock()
	defer v.Unlock()

	v.level = level
}

func (v *object) Level() models.Level {
	v.RLock()
	defer v.RUnlock()

	return v.level
}
