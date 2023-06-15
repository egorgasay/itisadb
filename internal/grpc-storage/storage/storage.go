package storage

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	_ "github.com/egorgasay/dockerdb/v2"
	"itisadb/internal/grpc-storage/config"
	tlogger "itisadb/internal/grpc-storage/transaction-logger"
	"itisadb/pkg/logger"
	"os"
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
	InitTLogger(Type string, dir string) error
	Set(key string, val string, unique bool) error
	WriteSet(key string, val string)
	Get(key string) (string, error)
	GetFromIndex(name string, key string) (string, error)
	SetToIndex(name string, key string, value string, uniques bool) error
	AttachToIndex(dst string, src string) error
	DeleteIndex(name string) error
	CreateIndex(name string) (err error)
	GetIndex(name string, prefix string) (map[string]string, error)
	Size(name string) (uint64, error)
	IsIndex(name string) bool
	Save() error
	DeleteIfExists(key string)
	Delete(key string) error
	DeleteAttr(name string, key string) error
	WriteDelete(key string)
	NoTLogger() bool
	GetFromDisk(key string) (string, error)
	GetFromDiskIndex(name, key string) (string, error)
}

type Storage struct {
	ramStorage ramStorage
	indexes    indexes
	tLogger    *tlogger.TransactionLogger
	logger     logger.ILogger
	noTLogger  bool
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

func New(cfg *config.Config, logger logger.ILogger, enableTLogger bool) (*Storage, error) {
	if cfg == nil {
		return nil, errors.New("empty configuration")
	}

	st := &Storage{
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](100000), RWMutex: &sync.RWMutex{}, path: "C:\\tmp"},
		indexes:    indexes{Map: swiss.NewMap[string, ivalue](100000), RWMutex: &sync.RWMutex{}, path: "C:\\tmp"},
		logger:     logger,
		noTLogger:  true,
	}

	if enableTLogger {
		tl, err := tlogger.New(cfg.TLoggerDir)
		if err != nil {
			return nil, fmt.Errorf("can't create transaction logger: %w", err)
		}

		st.tLogger = tl
		st.noTLogger = false
	}

	return st, nil
}

func (s *Storage) NoTLogger() bool {
	return s.noTLogger
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

func (s *Storage) WriteSet(key, val string) {
	s.tLogger.WriteSet(key, val)
}

func (s *Storage) WriteDelete(key string) {
	s.tLogger.WriteDelete(key)
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
	path := strings.Split(name, "/")
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
			m, err := s.GetIndex(name+"/"+key, prefix)
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
	path := strings.Split(name, "/")

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

func (is *indexes) save() (err error) {
	path := is.path
	if err = os.MkdirAll(path, 0777); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	is.Lock()
	is.Iter(func(key string, value ivalue) bool {
		if value.IsIndex() {
			err = value.save(path + "/" + key)
			if err != nil {
				return true
			}
		}

		return false
	})
	is.Unlock()

	return err
}

func (rs *ramStorage) save() (err error) {
	path := rs.path

	if err = os.MkdirAll(path, 0777); err != nil {
		return err
	}

	rs.Lock()

	var f *os.File
	rs.Iter(func(key string, value string) bool {
		f, err = os.OpenFile(path+"/"+key, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil && !os.IsExist(err) {
			return true
		}
		defer f.Close()

		_, err = f.WriteString(value)
		if err != nil {
			return true
		}

		return false
	})
	rs.Unlock()

	return err
}

func (s *Storage) Save() error {
	fmt.Println("Saving pairs to disk...")
	err := s.ramStorage.save()
	if err != nil {
		return err
	}

	fmt.Println("Saving indexes to disk...")
	err = s.indexes.save()
	if err != nil {
		return err
	}

	if s.tLogger == nil {
		return nil
	}
	return s.tLogger.Clear()
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

func (s *Storage) GetFromDisk(key string) (string, error) {
	path := s.ramStorage.path + "/" + key
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotFound
		}
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var buf bytes.Buffer
	for scanner.Scan() {
		buf.Write(scanner.Bytes())
	}

	return buf.String(), nil
}

func (s *Storage) GetFromDiskIndex(name, key string) (string, error) {
	path := s.indexes.path + "/" + name + "/" + key
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrIndexNotFound
		}
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var buf bytes.Buffer
	for scanner.Scan() {
		buf.Write(scanner.Bytes())
	}

	return buf.String(), nil
}
