package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/models"

	"github.com/dolthub/swiss"
)

type Storage struct {
	ramStorage  ramStorage
	objects     objects
	users       users
	objectsInfo objectsInfo
}

type ramStorage struct {
	*swiss.Map[string, models.Value]
	*sync.RWMutex
}

type objects struct {
	*swiss.Map[string, Something]
	*sync.RWMutex
}

type users struct {
	*swiss.Map[int, models.User]
	*sync.RWMutex
}

type Something interface {
	Object() gost.Option[*object]
	IsObject() bool

	Value() gost.Option[value]
	IsValue() bool
}

type objectsInfo struct {
	*swiss.Map[string, models.ObjectInfo]
	*sync.RWMutex
}

func New() (*Storage, error) {
	st := &Storage{
		objectsInfo: objectsInfo{Map: swiss.NewMap[string, models.ObjectInfo](10_000), RWMutex: &sync.RWMutex{}},
		ramStorage:  ramStorage{Map: swiss.NewMap[string, models.Value](10_000_000), RWMutex: &sync.RWMutex{}},
		objects:     objects{Map: swiss.NewMap[string, Something](100_000), RWMutex: &sync.RWMutex{}},
		users:       users{Map: swiss.NewMap[int, models.User](100), RWMutex: &sync.RWMutex{}},
	}

	return st, nil
}

func (s *Storage) Set(key, val string, opts models.SetOptions) error {
	s.ramStorage.Lock()
	defer s.ramStorage.Unlock()

	s.ramStorage.Put(key, models.Value{ReadOnly: opts.ReadOnly, Level: opts.Level, Value: val})

	return nil
}

func (s *Storage) Get(key string) (models.Value, error) {
	s.ramStorage.RLock()
	defer s.ramStorage.RUnlock()

	val, ok := s.ramStorage.Get(key)
	if !ok {
		return models.Value{}, constants.ErrNotFound
	}

	return val, nil
}

func (s *Storage) GetFromObject(name, key string) (string, error) {
	s.objects.RLock()
	defer s.objects.RUnlock()

	v, err := s.findObject(name)
	if err != nil {
		return "", err
	}

	return v.Get(key)
}

func (s *Storage) SetToObject(name, key, value string, opts models.SetToObjectOptions) error {
	s.objects.Lock()
	defer s.objects.Unlock()

	obj, err := s.findObject(name)
	if err != nil {
		return err
	}

	if opts.ReadOnly && obj.Has(key) {
		return constants.ErrAlreadyExists
	}

	obj.Set(key, value)
	return nil
}

func (s *Storage) AttachToObject(dst, src string) error {
	s.objects.Lock()
	defer s.objects.Unlock()

	object1, err := s.findObject(dst)
	if err != nil {
		return err
	}

	object2, err := s.findObject(src)
	if err != nil {
		return err
	}

	if object1.IsAttached(object2.Name()) {
		return constants.ErrCircularAttachment
	}

	err = object1.AttachObject(object2)
	if err != nil {
		return err
	}

	info, err := s.GetObjectInfo(dst)
	if err != nil {
		return err
	}

	s.AddObjectInfo(fmt.Sprintf("%s.%s", dst, src), info)

	return err
}

func (s *Storage) DeleteObject(name string) error {
	s.objects.Lock()
	defer s.objects.Unlock()

	split := strings.Split(name, ".")
	if name == "" || len(split) == 0 {
		return constants.ErrEmptyObjectName
	}
	objName := split[len(split)-1]

	parent := name
	if len(split) > 1 {
		parent = strings.Join(split[:len(split)-1], ".")
	} else {
		if !s.objects.Has(name) {
			return constants.ErrNotFound
		}
		s.objects.Delete(name)
		return nil
	}

	par, err := s.findObject(parent)
	if err != nil {
		return err
	}

	return par.Delete(objName)
}

// CreateObject ..
func (s *Storage) CreateObject(name string, opts models.ObjectOptions) (err error) {
	s.objects.Lock()
	defer s.objects.Unlock()

	path := strings.Split(name, constants.ObjectSeparator)
	if name == "" || len(path) == 0 {
		return constants.ErrEmptyObjectName
	}

	var val *object

	some, ok := s.objects.Get(path[0])
	if !ok { // TODO: || val.IsEmpty() {
		some = NewObject(path[0], nil, opts.Level)
		s.objects.Put(path[0], some)
	} else {
		switch o := some.Object(); o.IsSome() {
		case true:
			val = o.Unwrap()
		default:
			return constants.ErrSomethingExists
		}
	}

	path = path[1:]

	for _, objectName := range path {
		switch o := val.NextOrCreate(objectName, opts.Level).Object(); o.IsSome() {
		case true:
			if val = o.Unwrap(); val.IsEmpty() {
				val.RecreateObject()
			}
		default:
			return constants.ErrSomethingExists
		}
	}

	if val != nil && val.Level() != opts.Level {
		val.setLevel(opts.Level)
	}

	return nil
}

// TODO: JSONToObject
func (s *Storage) ObjectToJSON(name string) (string, error) {
	s.objects.RLock()
	defer s.objects.RUnlock()

	obj, err := s.findObject(name)
	if err != nil {
		return "", err
	}

	en, err := json.MarshalIndent(obj, "", "\t")
	return string(en), err
}

func (s *Storage) findObject(name string) (*object, error) {
	path := strings.Split(name, ".")

	if len(path) == 0 {
		return nil, constants.ErrObjectNotFound
	}

	var (
		val Something
		ok  bool
	)

	val, ok = s.objects.Get(path[0])
	if !ok {
		return nil, constants.ErrObjectNotFound
	}

	path = path[1:]

	for _, objectName := range path {
		switch o := val.Object(); o.IsSome() {
		case true:
			val, ok = o.Unwrap().GetValue(objectName)
			if !ok {
				return nil, constants.ErrObjectNotFound
			}
		default:
			return nil, constants.ErrSomethingExists
		}
	}

	// TODO: || val.IsEmpty()  ???

	switch o := val.Object(); o.IsSome() {
	case true:
		return o.Unwrap(), nil
	default:
		return nil, constants.ErrObjectNotFound
	}
}

// Size returns the size of the object
func (s *Storage) Size(name string) (uint64, error) {
	s.objects.RLock()
	defer s.objects.RUnlock()

	object, err := s.findObject(name)
	if err != nil {
		return 0, err
	}
	return uint64(object.Size()), nil
}

func (s *Storage) IsObject(name string) bool {
	s.objectsInfo.RLock()
	defer s.objectsInfo.RUnlock()

	_, ok := s.objectsInfo.Get(name)
	return ok
}

func (s *Storage) DeleteIfExists(key string) {
	s.ramStorage.Lock()
	defer s.ramStorage.Unlock()

	s.ramStorage.Delete(key)
}

func (s *Storage) Delete(key string) error {
	s.ramStorage.Lock()
	defer s.ramStorage.Unlock()
	if _, ok := s.ramStorage.Get(key); !ok {
		return constants.ErrNotFound
	}

	s.ramStorage.Delete(key)

	return nil
}

func (s *Storage) DeleteAttr(name, key string) error {
	s.objects.Lock()
	defer s.objects.Unlock()

	object, err := s.findObject(name)
	if err != nil {
		return err
	}

	switch o := object.Object(); o.IsSome() {
	case true:
		return o.Unwrap().Delete(key)
	default:
		return constants.ErrNotFound
	}
}

func (s *Storage) AddObjectInfo(name string, info models.ObjectInfo) {
	s.objectsInfo.Lock()
	defer s.objectsInfo.Unlock()

	s.objectsInfo.Put(name, info)
}

func (s *Storage) GetObjectInfo(name string) (models.ObjectInfo, error) {
	s.objectsInfo.RLock()
	defer s.objectsInfo.RUnlock()

	val, ok := s.objectsInfo.Get(name)
	if !ok {
		return models.ObjectInfo{}, constants.ErrNotFound
	}

	return val, nil
}

func (s *Storage) DeleteObjectInfo(name string) {
	s.objectsInfo.Lock()
	defer s.objectsInfo.Unlock()

	s.objectsInfo.Delete(name)
}

func (s *Storage) GetUserIDByName(username string) (r gost.Result[int]) {
	s.users.RLock()
	defer s.users.RUnlock()

	var find *int

	s.users.Iter(func(k int, v models.User) (stop bool) {
		if v.Login == username {
			find = &v.ID
			return true
		}
		return false
	})

	if find == nil {
		return r.Err(constants.ErrNotFound)
	}

	return r.Ok(*find)
}
