package storage

import (
	"github.com/dolthub/swiss"
	"sync"
)

type value struct {
	value      string
	next       *swiss.Map[string, ivalue]
	mutex      *sync.RWMutex
	isIndex    bool
	isAttached bool
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
	Next(name string) (val ivalue, ok bool)
	Size() int
	Iter(func(k string, v ivalue) bool)
	AttachIndex(name string, val ivalue) error
	DeleteIndex()
	Delete(key string) error
	Has(key string) bool
	IsAttached() bool
	setAttached()
}

func NewIndex() *value {
	return &value{mutex: &sync.RWMutex{}, next: swiss.NewMap[string, ivalue](100), isIndex: true}
}

func (v *value) GetValue() (val string) {
	v.mutex.RLock()
	val = v.value
	v.mutex.RUnlock()
	return val
}

func (v *value) IsEmpty() bool {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	return v.next == nil
}

func (v *value) Get(name string) (string, error) {
	v.mutex.RLock()
	val, ok := v.next.Get(name)
	v.mutex.RUnlock()
	if !ok {
		return "", ErrNotFound
	}

	return val.GetValue(), nil
}

func (v *value) Size() int {
	v.mutex.RLock()
	size := v.next.Count()
	v.mutex.RUnlock()
	return size
}

func (v *value) IsIndex() bool {
	return v.isIndex
}

func (v *value) IsAttached() bool {
	return v.isAttached
}

func (v *value) setAttached() {
	v.isAttached = true
}

func (v *value) AttachIndex(name string, val ivalue) (err error) {
	v.mutex.Lock()
	defer func() {
		v.mutex.Unlock()
		if err == nil {
			val.setAttached()
		}
	}()
	if v.next == nil {
		v.next = swiss.NewMap[string, ivalue](100)
		v.next.Put(name, val)
		return nil
	}
	if v.next.Has(name) {
		return nil
	}
	v.isIndex = true
	v.next.Put(name, val)
	return nil
}

func (v *value) Iter(f func(k string, v ivalue) bool) {
	v.next.Iter(f)
}

func (v *value) DeleteIndex() {
	v.mutex.Lock()
	v.next = nil
	v.mutex.Unlock()
}

func (v *value) NextOrCreate(name string) ivalue {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	val, ok := v.next.Get(name)
	if !ok {
		blank := NewIndex()
		v.next.Put(name, blank)
		return blank
	}

	return val
}

func (v *value) Next(name string) (ivalue, bool) {
	v.mutex.RLock()
	val, ok := v.next.Get(name)
	v.mutex.RUnlock()
	if !ok {
		return nil, false
	}

	return val, true
}

func (v *value) Set(key, val string) {
	v.mutex.Lock()
	v.next.Put(key, &value{value: val, mutex: &sync.RWMutex{}})
	v.mutex.Unlock()
}

func (v *value) SetValueUnique(val string) error {
	v.mutex.Lock()
	if v.next.Has(val) {
		v.mutex.Unlock()
		return ErrAlreadyExists
	}
	v.value = val
	v.mutex.Unlock()
	return nil
}

func (v *value) CreateIndex(name string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	if v.next.Has(name) {
		return
	}

	v.next.Put(name, &value{next: swiss.NewMap[string, ivalue](100), mutex: &sync.RWMutex{}})
}

func (v *value) RecreateIndex() {
	v.mutex.Lock()
	v.next = swiss.NewMap[string, ivalue](10000)
	v.mutex.Unlock()
}

func (v *value) Delete(key string) error {
	v.mutex.Lock()
	ok := v.next.Has(key)
	v.mutex.Unlock()

	if !ok {
		return ErrNotFound
	}

	v.next.Delete(key)
	return nil
}

func (v *value) Has(key string) bool {
	v.mutex.RLock()
	ok := v.next.Has(key)
	v.mutex.RUnlock()
	return ok
}

//func (v *value) Lock() {
//	v.mutex.Lock()
//}
//
//func (v *value) Unlock() {
//	v.mutex.Unlock()
//}
