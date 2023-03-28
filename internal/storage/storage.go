package storage

import (
	"context"
	"errors"
	_ "github.com/egorgasay/dockerdb/v2"
	"github.com/egorgasay/grpc-storage/config"
	"github.com/egorgasay/grpc-storage/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var NotFoundErr = errors.New("the value does not exist")

type Storage struct {
	DBStore *mongo.Database
	sync.RWMutex
	RAMStorage map[string]string
}

func New(cfg *config.DBConfig) (*Storage, error) {
	if cfg == nil {
		return nil, errors.New("empty configuration")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DataSourceCred))
	if err != nil {
		return nil, err
	}

	return &Storage{
		DBStore:    client.Database("grpc-server"),
		RAMStorage: make(map[string]string, 10),
	}, nil
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
			return "", NotFoundErr
		} else if err != nil {
			return "", nil
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

func (s *Storage) Save() {
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
			log.Println(err)
		}
	}
}
