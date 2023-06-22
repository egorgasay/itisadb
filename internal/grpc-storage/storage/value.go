package storage

import (
	"github.com/dolthub/swiss"
	"sync"
)

type value struct {
	value      string
	name       string
	Next       *swiss.Map[string, ivalue]
	mutex      *sync.RWMutex
	isIndex    bool
	attachedTo map[string]bool
}

type ivalue interface {
	GetValue() (val string)
	Get(name string) (val string, err error)
	Set(key, val string)
	SetValueUnique(val string) error
	CreateIndex(name string)
	RecreateIndex()
	IsIndex() bool
	IsEmpty() bool
	NextOrCreate(name string) (val ivalue)
	GetIValue(name string) (val ivalue, ok bool)
	Size() int
	Iter(func(k string, v ivalue) bool)
	AttachIndex(src ivalue) error
	DeleteIndex()
	Delete(key string) error
	Has(key string) bool
	IsAttached(name string) bool
	setAttached(attachedTo map[string]bool)
	Name() string
}

func NewIndex(name string, attachedTo map[string]bool) *value {
	if attachedTo == nil {
		attachedTo = make(map[string]bool)
	}
	attachedTo[name] = true
	return &value{
		mutex:      &sync.RWMutex{},
		Next:       swiss.NewMap[string, ivalue](10),
		attachedTo: attachedTo,
		isIndex:    true,
		name:       name,
	}
}

func (v *value) GetValue() (val string) {
	v.mutex.RLock()
	val = v.value
	v.mutex.RUnlock()
	return val
}

func (v *value) Name() string {
	return v.name
}

func (v *value) IsEmpty() bool {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	return v.Next == nil
}

func (v *value) Get(name string) (string, error) {
	v.mutex.RLock()
	val, ok := v.Next.Get(name)
	v.mutex.RUnlock()
	if !ok {
		return "", ErrNotFound
	}

	return val.GetValue(), nil
}

func (v *value) Size() int {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.Next.Count()
}

func (v *value) IsIndex() bool {
	return v.isIndex
}

func (v *value) IsAttached(name string) bool {
	return v.attachedTo[name]
}

func (v *value) setAttached(attachedTo map[string]bool) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	for k := range attachedTo {
		v.attachedTo[k] = true
	}
}

func (v *value) AttachIndex(src ivalue) (err error) {
	v.mutex.Lock()
	defer func() {
		v.mutex.Unlock()
		if err == nil {
			src.setAttached(v.attachedTo)
		}
	}()
	if v.Next == nil {
		v.Next = swiss.NewMap[string, ivalue](10)
		v.Next.Put(src.Name(), src)
		return nil
	}
	if v.Next.Has(src.Name()) {
		return nil
	}
	v.isIndex = true
	v.Next.Put(src.Name(), src)
	return nil
}

func (v *value) Iter(f func(k string, v ivalue) bool) {
	v.Next.Iter(f)
}

func (v *value) DeleteIndex() {
	v.mutex.Lock()
	v.Next = nil
	v.mutex.Unlock()
}

func (v *value) NextOrCreate(name string) ivalue {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	val, ok := v.Next.Get(name)
	if !ok {
		blank := NewIndex(name, v.attachedTo)
		v.Next.Put(name, blank)
		return blank
	}

	return val
}

func (v *value) GetIValue(name string) (ivalue, bool) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	val, ok := v.Next.Get(name)
	if !ok {
		return nil, false
	}

	return val, true
}

func (v *value) Set(key, val string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.Next.Put(key, &value{value: val, mutex: &sync.RWMutex{}, name: key})
}

func (v *value) SetValueUnique(val string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if v.Next.Has(val) {
		return ErrAlreadyExists
	}

	v.value = val
	return nil
}

func (v *value) CreateIndex(name string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	if v.Next.Has(name) {
		return
	}

	v.Next.Put(name, &value{Next: swiss.NewMap[string, ivalue](100), mutex: &sync.RWMutex{}})
}

func (v *value) RecreateIndex() {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.Next = swiss.NewMap[string, ivalue](10000)
}

func (v *value) Delete(key string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if !v.Next.Has(key) {
		return ErrNotFound
	}

	v.Next.Delete(key)
	return nil
}

func (v *value) Has(key string) bool {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.Next.Has(key)
}
