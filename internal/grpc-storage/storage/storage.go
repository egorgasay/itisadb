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
	transactionlogger "grpc-storage/internal/grpc-storage/transaction-logger"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
	"grpc-storage/internal/schema"
	"grpc-storage/pkg/logger"
	"log"
	"sync"
)

var ErrNotFound = errors.New("the value does not exist")

type Storage struct {
	DBStore *mongo.Database
	sync.RWMutex
	RAMStorage map[string]string
	TLogger    transactionlogger.ITransactionLogger
	logger     logger.ILogger
}

func New(cfg *config.DBConfig, logger logger.ILogger) (*Storage, error) {
	if cfg == nil {
		return nil, errors.New("empty configuration")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DataSourceCred))
	if err != nil {
		return nil, err
	}

	st := &Storage{
		DBStore:    client.Database("grpc-server"),
		RAMStorage: make(map[string]string, 10),
		logger:     logger,
	}

	err = st.InitTLogger()
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Storage) InitTLogger() error {
	var err error

	s.TLogger, err = transactionlogger.NewTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errs := s.TLogger.ReadEvents()
	e, ok := service.Event{}, true
	s.TLogger.Run()
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

	return err
}

func (s *Storage) Set(key string, val string) {
	s.Lock()
	defer s.Unlock()
	s.RAMStorage[key] = val
}

func (s *Storage) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.RAMStorage[key]
	if !ok {
		res, err := s.get(key)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrNotFound
		} else if err != nil {
			return "", err
		}
		s.RAMStorage[key] = res

		return res, nil
	}

	return val, nil
}

// get gets value by key from db.
func (s *Storage) get(key string) (string, error) {
	c := s.DBStore.Collection("map")
	filter := bson.D{{"Key", key}}

	var kv schema.KeyValue
	ctx := context.Background()
	if err := c.FindOne(ctx, filter).Decode(&kv); err != nil {
		log.Println(err)
		return "", err
	}

	return kv.Value, nil
}

func (s *Storage) Save() error {
	c := s.DBStore.Collection("map")
	s.Lock()
	defer s.Unlock()

	ctx := context.Background()
	opts := options.Update().SetUpsert(true)
	for key, value := range s.RAMStorage {
		filter := bson.D{{"Key", key}}
		update := bson.D{{"$set", bson.D{{"Key", key}, {"Value", value}}}}
		_, err := c.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			s.logger.Warn(err.Error())
		}
	}

	return s.TLogger.Clear()
}
