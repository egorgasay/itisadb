package storage

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/egorgasay/dockerdb/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"itisadb/internal/grpc-storage/config"
	tlogger "itisadb/internal/grpc-storage/transaction-logger"
	"itisadb/internal/grpc-storage/transaction-logger/service"
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

type Storage struct {
	dbStore    *mongo.Database
	ramStorage ramStorage
	indexes    indexes
	tLogger    tlogger.ITransactionLogger
	logger     logger.ILogger
}

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
	GetIndex(name string) (map[string]string, error)
	Size(name string) (uint64, error)
	IsIndex(name string) bool
	Save() error
	DeleteIfExists(key string)
	Delete(key string) error
	DeleteAttr(name string, key string) error
}

type ramStorage struct {
	*swiss.Map[string, string]
	*sync.RWMutex
}

type indexes struct {
	*swiss.Map[string, ivalue]
	*sync.RWMutex
}

func NewWithTLogger(cfg *config.Config, logger logger.ILogger) (*Storage, error) {
	st, err := New(cfg, logger)
	if err != nil {
		return nil, err
	}

	err = st.InitTLogger(cfg.TLoggerType, cfg.TLoggerDir)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func New(cfg *config.Config, logger logger.ILogger) (*Storage, error) {
	if cfg == nil {
		return nil, errors.New("empty configuration")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBConfig.DataSourceCred))
	if err != nil {
		return nil, err
	}
	st := &Storage{
		dbStore:    client.Database("grpc-server"),
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](100000), RWMutex: &sync.RWMutex{}},
		indexes:    indexes{Map: swiss.NewMap[string, ivalue](100000), RWMutex: &sync.RWMutex{}},
		logger:     logger,
	}

	return st, nil
}

func (s *Storage) InitTLogger(Type string, dir string) error {
	var err error
	s.tLogger, err = tlogger.NewTransactionLogger(Type, dir)
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errs := s.tLogger.ReadEvents()
	e, ok := service.Event{}, true
	s.tLogger.Run()
	for ok && err == nil {
		select {
		case err, ok = <-errs:
		case e, ok = <-events:
			switch e.EventType {
			case service.Set:
				s.Set(e.Key, e.Value, false)
			case service.Delete:
			}
		}
	}

	return nil
}

func (s *Storage) Set(key, val string, unique bool) error {
	s.ramStorage.Lock()
	if unique && s.ramStorage.Has(key) {
		return ErrAlreadyExists
	}
	s.ramStorage.Put(key, val)

	s.ramStorage.Unlock()
	return nil
}

func (s *Storage) WriteSet(key, val string) {
	s.tLogger.WriteSet(key, val)
}

func (s *Storage) Get(key string) (string, error) {
	s.ramStorage.RLock()

	val, ok := s.ramStorage.Get(key)
	s.ramStorage.RUnlock()
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

var ErrWrongIndexName = errors.New("wrong index name provided")

func (s *Storage) AttachToIndex(dst, src string) error {
	index1, err := s.findIndex(dst)
	if err != nil {
		return err
	}

	index2, err := s.findIndex(src)
	if err != nil {
		return err
	}

	source := strings.Split(src, "/")
	if len(source) == 0 {
		return ErrWrongIndexName // TODO: catch
	}

	err = index1.AttachIndex(source[len(source)-1], index2)
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
		val = NewIndex()
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
		val.CreateIndex(indexName)
	}
	return nil
}

func (s *Storage) GetIndex(name string) (map[string]string, error) {
	index, err := s.findIndex(name)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	index.Iter(func(key string, value ivalue) bool {
		k := ""
		if value.IsIndex() {
			k = "index"
		} else {
			k = value.GetValue()
		}
		result[key] = k
		return false
	})
	return result, nil
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

func (s *Storage) Save() error {
	c := s.dbStore.Collection("map")
	s.ramStorage.Lock()

	ctx := context.Background()
	opts := options.Update().SetUpsert(true)
	s.ramStorage.Iter(func(key string, value string) bool {
		filter := bson.D{{"Key", key}}
		update := bson.D{{"$set", bson.D{{"Key", key}, {"Value", value}}}}
		_, err := c.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			s.logger.Warn(err.Error())
		}
		return true
	})

	s.ramStorage.Unlock()
	return s.tLogger.Clear()
}

func (s *Storage) DeleteIfExists(key string) {
	s.ramStorage.Lock()
	s.ramStorage.Delete(key)
	s.ramStorage.Unlock()
}

func (s *Storage) Delete(key string) error {
	s.ramStorage.Lock()
	if _, ok := s.ramStorage.Get(key); !ok {
		s.ramStorage.Unlock()
		return ErrNotFound
	}

	s.ramStorage.Delete(key)
	s.ramStorage.Unlock()

	return nil
}

func (s *Storage) DeleteAttr(name, key string) error {
	index, err := s.findIndex(name)
	if err != nil {
		return err
	}

	return index.Delete(key)
}
