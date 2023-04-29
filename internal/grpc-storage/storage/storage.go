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
	"sync"

	"github.com/dolthub/swiss"
)

var ErrNotFound = errors.New("the value does not exist")
var ErrAlreadyExists = errors.New("the value already exists")

type Storage struct {
	dbStore *mongo.Database
	sync.RWMutex
	ramStorage *swiss.Map[string, string]
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
		ramStorage: swiss.NewMap[string, string](100000),
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
			case service.Delete:
			case service.Set:
				s.Set(e.Key, e.Value, false)
			}
		}
	}

	return nil
}

func (s *Storage) Set(key, val string, unique bool) error {
	s.Lock()
	defer s.Unlock()
	if unique {
		if _, ok := s.ramStorage.Get(key); ok {
			return ErrAlreadyExists
		}
	}
	s.ramStorage.Put(key, val)
	return nil
}

func (s *Storage) WriteSet(key, val string) {
	s.tLogger.WriteSet(key, val)
}

func (s *Storage) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.ramStorage.Get(key)
	if !ok {
		return "", ErrNotFound
	}

	return val, nil
}

func (s *Storage) Save() error {
	c := s.dbStore.Collection("map")
	s.Lock()
	defer s.Unlock()

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

	return s.tLogger.Clear()
}
