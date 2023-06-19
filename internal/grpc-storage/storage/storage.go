package storage

import (
	"errors"
	_ "github.com/egorgasay/dockerdb/v2"
	"itisadb/pkg/logger"
	"strings"
	"sync"

	"github.com/dolthub/swiss"
)

var ErrNotFound = errors.New("the value does not exist")
var ErrAlreadyExists = errors.New("the value already exists")
var ErrIndexNotFound = errors.New("index not found")
var ErrSomethingExists = errors.New("something with this name already exists")
var ErrEmptyIndexName = errors.New("index name is empty")
var ErrCircularAttachment = errors.New("circular attachment not allowed")

type IStorage interface {
	Set(key string, val string, unique bool) error
	Get(key string) (string, error)
	DeleteIfExists(key string)
	Delete(key string) error

	SetToIndex(name string, key string, value string, uniques bool) error
	GetFromIndex(name string, key string) (string, error)
	AttachToIndex(dst string, src string) error
	DeleteIndex(name string) error
	CreateIndex(name string) (err error)
	GetIndex(name string, prefix string) (map[string]string, error)
	Size(name string) (uint64, error)
	IsIndex(name string) bool
	DeleteAttr(name string, key string) error
}

type Storage struct {
	ramStorage ramStorage
	indexes    indexes
	logger     logger.ILogger
}

type ramStorage struct {
	*swiss.Map[string, string]
	path string
	*sync.RWMutex
}

type indexes struct {
	*swiss.Map[string, ivalue]
	*sync.RWMutex
	path string
}

func New(logger logger.ILogger) (*Storage, error) {
	st := &Storage{
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](10000000), RWMutex: &sync.RWMutex{}, path: "C:\\tmp"},
		indexes:    indexes{Map: swiss.NewMap[string, ivalue](100000), RWMutex: &sync.RWMutex{}, path: "C:\\tmp"},
		logger:     logger,
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

func (s *Storage) GetFromIndex(name, key string) (string, error) {
	v, err := s.findIndex(name)
	if err != nil {
		return "", err
	}

	return v.Get(key)
}

func (s *Storage) SetToIndex(name, key, value string, uniques bool) error {
	index, err := s.findIndex(name)
	if err != nil {
		return err
	}

	if uniques && index.Has(key) {
		return ErrAlreadyExists
	}

	index.Set(key, value)
	return nil
}

func (s *Storage) AttachToIndex(dst, src string) error {
	index1, err := s.findIndex(dst)
	if err != nil {
		return err
	}

	index2, err := s.findIndex(src)
	if err != nil {
		return err
	}

	if index2.IsAttached(index1.Name()) {
		return ErrCircularAttachment
	}

	err = index1.AttachIndex(index2)
	return err
}

func (s *Storage) DeleteIndex(name string) error {
	val, err := s.findIndex(name)
	if err != nil {
		return err
	}

	val.DeleteIndex()

	return nil
}

func (s *Storage) CreateIndex(name string) (err error) {
	path := strings.Split(name, ".")
	if name == "" || len(path) == 0 {
		return ErrEmptyIndexName
	}

	val, ok := s.indexes.Get(path[0])
	if !ok || val.IsEmpty() {
		s.indexes.Lock()
		val = NewIndex(path[0], nil)
		s.indexes.Put(path[0], val)
		s.indexes.Unlock()
	}

	path = path[1:]

	for _, indexName := range path {
		val = val.NextOrCreate(indexName)
		if !val.IsIndex() {
			return ErrSomethingExists
		} else if val.IsEmpty() {
			val.RecreateIndex()
		}
	}
	return nil
}

func (s *Storage) GetIndex(name string, prefix string) (map[string]string, error) {
	index, err := s.findIndex(name)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	index.Iter(func(key string, value ivalue) bool {
		if value.IsIndex() {
			prefix = prefix + "\t"
			m, err := s.GetIndex(name+"."+key, prefix)
			if err != nil {
				result[key] = err.Error()
			} else {
				result[key] = mapToString(m, prefix)
			}
		} else {
			result[key] = value.GetValue()
		}
		return false
	})
	return result, nil
}

func mapToString(m map[string]string, prefix string) string {
	b := strings.Builder{}
	for k, v := range m {
		b.WriteString("\n")
		b.WriteString(prefix)
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v)
	}
	return b.String()
}

func (s *Storage) findIndex(name string) (ivalue, error) {
	path := strings.Split(name, ".")

	if len(path) == 0 {
		return nil, ErrIndexNotFound
	}

	val, ok := s.indexes.Get(path[0])
	if !ok {
		return nil, ErrIndexNotFound
	}

	path = path[1:]

	for _, indexName := range path {
		switch val.IsIndex() {
		case true:
			val, ok = val.Next(indexName)
			if !ok {
				return nil, ErrIndexNotFound
			}
		default:
			return nil, ErrSomethingExists
		}
	}

	if !val.IsIndex() || val.IsEmpty() {
		return nil, ErrIndexNotFound
	}

	return val, nil
}

// Size returns the size of the index
func (s *Storage) Size(name string) (uint64, error) {
	index, err := s.findIndex(name)
	if err != nil {
		return 0, err
	}
	return uint64(index.Size()), nil
}

func (s *Storage) IsIndex(name string) bool {
	if val, err := s.findIndex(name); err != nil {
		return false
	} else {
		return val.IsIndex()
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
	index, err := s.findIndex(name)
	if err != nil {
		return err
	}

	return index.Delete(key)
}
