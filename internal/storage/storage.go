package storage

import (
	"fmt"
	"strings"
	"sync"

	"itisadb/internal/constants"
	"itisadb/internal/models"

	"github.com/egorgasay/gost"

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
	*swiss.Map[string, models.User]
	*sync.RWMutex
	changeID uint64
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
		users:       users{Map: swiss.NewMap[string, models.User](100), RWMutex: &sync.RWMutex{}},
	}

	return st, nil
}

func (s *Storage) Set(key, val string, opts models.SetOptions) (r gost.ResultN) {
	s.ramStorage.Lock()
	defer s.ramStorage.Unlock()

	s.ramStorage.Put(key, models.Value{ReadOnly: opts.ReadOnly, Level: opts.Level, Value: val})

	return r.Ok()
}

func (s *Storage) Get(key string) (r gost.Option[models.Value]) {
	s.ramStorage.RLock()
	defer s.ramStorage.RUnlock()

	val, ok := s.ramStorage.Get(key)
	if !ok {
		return r.None()
	}

	return r.Some(val)
}

func (s *Storage) GetFromObject(name, key string) (r gost.Option[string]) {
	s.objects.RLock()
	defer s.objects.RUnlock()

	v := s.findObject(name)
	switch v.IsSome() {
	case true:
		return v.Unwrap().Get(key)
	default:
		return r.None()
	}
}

func (s *Storage) SetToObject(name, key, value string, opts models.SetToObjectOptions) (r gost.ResultN) {
	s.objects.Lock()
	defer s.objects.Unlock()

	obj := s.findObject(name)
	switch obj.IsSome() {
	case false:
		return r.Err(constants.ErrObjectNotFound)
	}

	object := obj.Unwrap()

	if opts.ReadOnly && object.Has(key) {
		return r.Err(constants.ErrAlreadyExists)
	}

	object.Set(key, value)
	return r.Ok()
}

func (s *Storage) AttachToObject(dst, src string) (r gost.ResultN) {
	s.objects.Lock()
	defer s.objects.Unlock()

	var obj1, obj2 *object

	object1 := s.findObject(dst)
	switch object1.IsSome() {
	case true:
		obj1 = object1.Unwrap()
	default:
		return r.Err(constants.ErrObjectNotFound)
	}

	object2 := s.findObject(src)
	switch object2.IsSome() {
	case true:
		obj2 = object2.Unwrap()
	default:
		return r.Err(constants.ErrObjectNotFound)
	}

	if obj1.IsAttached(obj2.Name()) {
		return r.Err(constants.ErrCircularAttachment)
	}

	rAttach := obj1.AttachObject(obj2)
	if rAttach.IsErr() {
		return r.Err(rAttach.Error())
	}

	infoR := s.GetObjectInfo(dst)
	if infoR.IsNone() {
		return r.Err(constants.ErrObjectNotFound)
	}

	info := infoR.Unwrap()
	s.AddObjectInfo(fmt.Sprintf("%s.%s", dst, src), info)

	return r.Ok()
}

func (s *Storage) DeleteObject(name string) (r gost.ResultN) {
	s.objects.Lock()
	defer s.objects.Unlock()

	split := strings.Split(name, ".")
	if name == "" || len(split) == 0 {
		return r.Err(constants.ErrEmptyObjectName)
	}
	objName := split[len(split)-1]

	parent := name
	if len(split) > 1 {
		parent = strings.Join(split[:len(split)-1], ".")
	} else {
		if !s.objects.Has(name) {
			return r.Err(constants.ErrObjectNotFound)
		}
		s.objects.Delete(name)
		return r
	}

	par := s.findObject(parent)
	switch par.IsSome() {
	case true:
		return par.Unwrap().Delete(objName)
	default:
		return r.Err(constants.ErrObjectNotFound)
	}
}

// CreateObject ..
func (s *Storage) CreateObject(name string, opts models.ObjectOptions) (r gost.ResultN) {
	s.objects.Lock()
	defer s.objects.Unlock()

	obj := s.findObject(name)
	switch obj.IsSome() {
	case true:
		o := obj.Unwrap()
		o.setLevel(opts.Level)
		s.objects.Put(name, o)
		return r
	}

	path := strings.Split(name, constants.ObjectSeparator)
	if name == "" || len(path) == 0 {
		return r.Err(constants.ErrEmptyObjectName)
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
			return r.Err(constants.ErrSomethingExists)
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
			return r.Err(constants.ErrSomethingExists)
		}
	}

	return r
}

// TODO: JSONToObject
func (s *Storage) ObjectToJSON(name string) (r gost.Result[string]) {
	s.objects.RLock()
	defer s.objects.RUnlock()

	switch obj := s.findObject(name); obj.IsSome() {
	case true:
		rMarshal := obj.Unwrap().MarshalJSON()
		if rMarshal.IsErr() {
			return r.Err(rMarshal.Error())
		}
		return r.Ok(string(rMarshal.Unwrap()))
	default:
		return r.Err(constants.ErrObjectNotFound)
	}
}

func (s *Storage) findObject(name string) (r gost.Option[*object]) {
	path := strings.Split(name, ".")

	if len(path) == 0 {
		return r.None()
	}

	var (
		val Something
		ok  bool
	)

	val, ok = s.objects.Get(path[0])
	if !ok {
		return r.None()
	}

	path = path[1:]

	for _, objectName := range path {
		switch o := val.Object(); o.IsSome() {
		case true:
			val, ok = o.Unwrap().GetValue(objectName)
			if !ok {
				return r.None()
			}
		default:
			return r.None()
		}
	}

	// TODO: || val.IsEmpty()  ???

	switch o := val.Object(); o.IsSome() {
	case true:
		return r.Some(o.Unwrap())
	default:
		return r.None()
	}
}

// Size returns the size of the object
func (s *Storage) Size(name string) (r gost.Result[uint64]) {
	s.objects.RLock()
	defer s.objects.RUnlock()

	object := s.findObject(name)
	switch object.IsSome() {
	case true:
		return r.Ok(uint64(object.Unwrap().Size()))
	default:
		return r.Err(constants.ErrObjectNotFound)
	}
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

func (s *Storage) Delete(key string) (r gost.ResultN) {
	s.ramStorage.Lock()
	defer s.ramStorage.Unlock()

	if _, ok := s.ramStorage.Get(key); !ok {
		return r.Err(constants.ErrNotFound)
	}

	s.ramStorage.Delete(key)

	return r
}

func (s *Storage) DeleteAttr(name, key string) (r gost.ResultN) {
	s.objects.Lock()
	defer s.objects.Unlock()

	switch object := s.findObject(name); object.IsSome() {
	case true:
		return object.Unwrap().Delete(key)
	default:
		return r.Err(constants.ErrObjectNotFound)
	}
}

func (s *Storage) AddObjectInfo(name string, info models.ObjectInfo) {
	s.objectsInfo.Lock()
	defer s.objectsInfo.Unlock()

	s.objectsInfo.Put(name, info)
}

func (s *Storage) GetObjectInfo(name string) (r gost.Option[models.ObjectInfo]) {
	s.objectsInfo.RLock()
	defer s.objectsInfo.RUnlock()

	val, ok := s.objectsInfo.Get(name)
	if !ok {
		return r.None()
	}

	return r.Some(val)
}

func (s *Storage) DeleteObjectInfo(name string) {
	s.objectsInfo.Lock()
	defer s.objectsInfo.Unlock()

	s.objectsInfo.Delete(name)
}

func (s *Storage) GetUsersFromChangeID(id uint64) gost.Result[[]models.User] {
	s.users.RLock()
	defer s.users.RUnlock()

	var find []models.User

	s.users.Iter(func(k string, v models.User) (stop bool) {
		if v.GetChangeID() < id {
			find = append(find, v)
		}
		return false
	})

	return gost.Ok(find)
}

func (s *Storage) GetUserChangeID() uint64 {
	s.users.RLock()
	defer s.users.RUnlock()

	return s.users.changeID
}
