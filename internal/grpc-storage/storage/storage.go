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

type Storage struct {
	dbStore *mongo.Database
	sync.RWMutex
	ramStorage *swiss.Map[string, any]
	indexes    sync.Map
	tLogger    tlogger.ITransactionLogger
	logger     logger.ILogger
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
		ramStorage: swiss.NewMap[string, any](100000),
		logger:     logger,
	}

	err = st.InitTLogger(cfg.TLoggerType, cfg.TLoggerDir)
	if err != nil {
		return nil, err
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
	s.Lock()
	if unique {
		if _, ok := s.ramStorage.Get(key); ok {
			return ErrAlreadyExists
		}
	}
	s.ramStorage.Put(key, val)

	s.Unlock()
	return nil
}

func (s *Storage) WriteSet(key, val string) {
	s.tLogger.WriteSet(key, val)
}

func (s *Storage) Get(key string) (string, error) {
	s.RLock()

	val, ok := s.ramStorage.Get(key)
	if !ok {
		return "", ErrNotFound
	}

	s.RUnlock()

	switch val.(type) {
	case string:
		return val.(string), nil
	default:
		return "", errors.New("wrong type")
	}
}

var ErrIndexNotFound = errors.New("index not found")

func (s *Storage) GetFromIndex(name, key string) (string, error) {
	path := strings.Split(name, "/")

	if len(path) == 0 {
		return "", ErrIndexNotFound
	}

	var index any
	var ok bool

	index, ok = s.indexes.Load(path[0])
	if !ok {
		return "", ErrIndexNotFound
	}

	for _, indexName := range path {
		switch index.(type) {
		case *swiss.Map[string, any]:
			ind := index.(*swiss.Map[string, any])
			index, ok = ind.Get(indexName)
			if !ok {
				return "", ErrIndexNotFound
			}
		default:
			return "", ErrIndexNotFound
		}
	}

	final, ok := index.(*swiss.Map[string, any])
	if !ok {
		return "", ErrIndexNotFound
	}

	value, ok := final.Get(key)
	if !ok {
		return "", ErrNotFound
	}

	val, ok := value.(string)
	if !ok {
		return "", ErrNotFound
	}

	return val, nil
}

func (s *Storage) SetToIndex(name, key, value string) error {
	path := strings.Split(name, "/")

	if len(path) == 0 {
		return ErrIndexNotFound
	}

	var index any
	var ok bool

	index, ok = s.indexes.Load(path[0])
	if !ok {
		return ErrIndexNotFound
	}

	for _, indexName := range path {
		switch index.(type) {
		case *swiss.Map[string, any]:
			ind := index.(*swiss.Map[string, any])
			index, ok = ind.Get(indexName)
			if !ok {
				return ErrIndexNotFound
			}
		default:
			return ErrIndexNotFound
		}
	}

	final, ok := index.(*swiss.Map[string, any])
	if !ok {
		return ErrIndexNotFound
	}

	final.Put(key, value)
	return nil
}

func (s *Storage) IsIndex(name string) (bool, error) {
	path := strings.Split(name, "/")

	if len(path) == 0 {
		return false, ErrNotFound
	}

	var index any
	var ok bool

	index, ok = s.indexes.Load(path[0])
	if !ok {
		return false, ErrNotFound
	}

	for _, indexName := range path {
		switch index.(type) {
		case *swiss.Map[string, any]:
			ind := index.(*swiss.Map[string, any])
			index, ok = ind.Get(indexName)
			if !ok {
				return false, ErrNotFound
			}
		default:
			return false, ErrNotFound
		}
	}

	_, ok = index.(*swiss.Map[string, any])
	return ok, nil
}

func (s *Storage) Save() error {
	c := s.dbStore.Collection("map")
	s.Lock()

	ctx := context.Background()
	opts := options.Update().SetUpsert(true)
	s.ramStorage.Iter(func(key string, value any) bool {
		filter := bson.D{{"Key", key}}
		update := bson.D{{"$set", bson.D{{"Key", key}, {"Value", value}}}}
		_, err := c.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			s.logger.Warn(err.Error())
		}
		return true
	})

	s.Unlock()
	return s.tLogger.Clear()
}
