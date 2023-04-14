package storage

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/egorgasay/dockerdb/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grpc-storage/internal/grpc-storage/config"
	tlogger "grpc-storage/internal/grpc-storage/transaction-logger"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
	"grpc-storage/pkg/logger"
	"sync"
)

var ErrNotFound = errors.New("the value does not exist")

type Storage struct {
	dbStore *mongo.Database
	sync.RWMutex
	ramStorage map[string]string
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
		ramStorage: make(map[string]string, 10),
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
				s.Set(e.Key, e.Value)
			}
		}
	}

	return nil
}

func (s *Storage) Set(key string, val string) {
	s.Lock()
	defer s.Unlock()
	s.ramStorage[key] = val
}

func (s *Storage) WriteSet(key string, val string) {
	s.tLogger.WriteSet(key, val)
}

func (s *Storage) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.ramStorage[key]
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
	for key, value := range s.ramStorage {
		filter := bson.D{{"Key", key}}
		update := bson.D{{"$set", bson.D{{"Key", key}, {"Value", value}}}}
		_, err := c.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			s.logger.Warn(err.Error())
		}
	}

	return s.tLogger.Clear()
}
