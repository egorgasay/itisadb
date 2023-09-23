package storage

import (
	"encoding/json"
	"errors"
	_ "github.com/egorgasay/dockerdb/v2"
	"strings"
	"sync"

	"github.com/dolthub/swiss"
)

var ErrNotFound = errors.New("the value does not exist")
var ErrAlreadyExists = errors.New("the value already exists")
var ErrObjectNotFound = errors.New("object not found")
var ErrSomethingExists = errors.New("something with this name already exists")
var ErrEmptyObjectName = errors.New("object name is empty")
var ErrCircularAttachment = errors.New("circular attachment not allowed")

type IStorage interface {
	Set(key string, val string, unique bool) error
	Get(key string) (string, error)
	DeleteIfExists(key string)
	Delete(key string) error

	SetToObject(name string, key string, value string, uniques bool) error
	GetFromObject(name string, key string) (string, error)
	AttachToObject(dst string, src string) error
	DeleteObject(name string) error
	CreateObject(name string) (err error)
	ObjectToJSON(name string) (string, error)
	Size(name string) (uint64, error)
	IsObject(name string) bool
	DeleteAttr(name string, key string) error
}

type Storage struct {
	ramStorage ramStorage
	objects    objects
}

type ramStorage struct {
	*swiss.Map[string, string]
	path string
	*sync.RWMutex
}

type objects struct {
	*swiss.Map[string, ivalue]
	*sync.RWMutex
	path string
}

func New() (*Storage, error) {
	st := &Storage{
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](10000000), RWMutex: &sync.RWMutex{}, path: "C:\\tmp"},
		objects:    objects{Map: swiss.NewMap[string, ivalue](100000), RWMutex: &sync.RWMutex{}, path: "C:\\tmp"},
	}

	return st, nil
}

func (s *Storage) Set(key, val string, unique bool) error {
	s.ramStorage.Lock()
	defer s.ramStorage.Unlock()
	if unique && s.ramStorage.Has(key) {
		return ErrAlreadyExists
	}
	s.ramStorage.Put(key, val)

	return nil
}

func (s *Storage) Get(key string) (string, error) {
	s.ramStorage.RLock()
	defer s.ramStorage.RUnlock()

	val, ok := s.ramStorage.Get(key)
	if !ok {
		return "", ErrNotFound
	}

	return val, nil
}

func (s *Storage) GetFromObject(name, key string) (string, error) {
	v, err := s.findObject(name)
	if err != nil {
		return "", err
	}

	return v.Get(key)
}

func (s *Storage) SetToObject(name, key, value string, uniques bool) error {
	object, err := s.findObject(name)
	if err != nil {
		return err
	}

	if uniques && object.Has(key) {
		return ErrAlreadyExists
	}

	object.Set(key, value)
	return nil
}

func (s *Storage) AttachToObject(dst, src string) error {
	object1, err := s.findObject(dst)
	if err != nil {
		return err
	}

	object2, err := s.findObject(src)
	if err != nil {
		return err
	}

	if object2.IsAttached(object1.Name()) {
		return ErrCircularAttachment
	}

	err = object1.AttachObject(object2)
	return err
}

func (s *Storage) DeleteObject(name string) error {
	val, err := s.findObject(name)
	if err != nil {
		return err
	}

	val.DeleteObject()

	return nil
}

func (s *Storage) CreateObject(name string) (err error) {
	path := strings.Split(name, ".")
	if name == "" || len(path) == 0 {
		return ErrEmptyObjectName
	}

	val, ok := s.objects.Get(path[0])
	if !ok || val.IsEmpty() {
		s.objects.Lock()
		val = NewObject(path[0], nil)
		s.objects.Put(path[0], val)
		s.objects.Unlock()
	}

	path = path[1:]

	for _, objectName := range path {
		val = val.NextOrCreate(objectName)
		if !val.IsObject() {
			return ErrSomethingExists
		} else if val.IsEmpty() {
			val.RecreateObject()
		}
	}
	return nil
}

// TODO: JSONToObject
func (s *Storage) ObjectToJSON(name string) (string, error) {
	object, err := s.findObject(name)
	if err != nil {
		return "", err
	}

	en, err := json.MarshalIndent(object, "", "\t")
	return string(en), err
}

func (v *value) MarshalJSON() ([]byte, error) {
	arr := make([]any, 0, 100)
	var data map[string]interface{}

	if v.values != nil {
		v.values.Iter(func(k string, v ivalue) bool {
			val := v.(*value)
			if val != nil {
				arr = append(arr, v.(*value))
			}

			return false
		})

		data = map[string]interface{}{
			"name":   v.Name(),
			"values": arr,
		}
	} else {
		data = map[string]interface{}{
			"name":  v.Name(),
			"value": v.value,
		}
	}

	return json.MarshalIndent(data, "", "\t")
}

func (s *Storage) findObject(name string) (ivalue, error) {
	path := strings.Split(name, ".")

	if len(path) == 0 {
		return nil, ErrObjectNotFound
	}

	val, ok := s.objects.Get(path[0])
	if !ok {
		return nil, ErrObjectNotFound
	}

	path = path[1:]

	for _, objectName := range path {
		switch val.IsObject() {
		case true:
			val, ok = val.GetIValue(objectName)
			if !ok {
				return nil, ErrObjectNotFound
			}
		default:
			return nil, ErrSomethingExists
		}
	}

	if !val.IsObject() || val.IsEmpty() {
		return nil, ErrObjectNotFound
	}

	return val, nil
}

// Size returns the size of the object
func (s *Storage) Size(name string) (uint64, error) {
	object, err := s.findObject(name)
	if err != nil {
		return 0, err
	}
	return uint64(object.Size()), nil
}

func (s *Storage) IsObject(name string) bool {
	if val, err := s.findObject(name); err != nil {
		return false
	} else {
		return val.IsObject()
	}
}

func (s *Storage) DeleteIfExists(key string) {
	s.ramStorage.Lock()
	s.ramStorage.Delete(key)
	s.ramStorage.Unlock()
}

func (s *Storage) Delete(key string) error {
	s.ramStorage.Lock()
	defer s.ramStorage.Unlock()
	if _, ok := s.ramStorage.Get(key); !ok {
		return ErrNotFound
	}

	s.ramStorage.Delete(key)

	return nil
}

func (s *Storage) DeleteAttr(name, key string) error {
	object, err := s.findObject(name)
	if err != nil {
		return err
	}

	return object.Delete(key)
}
