package storage

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"grpc-storage/internal/memory-balancer/config"
	"grpc-storage/internal/schema"
)

type Storage struct {
	DBStore *mongo.Database
}

func New(cfg *config.Config) (*Storage, error) {
	if cfg == nil {
		return nil, errors.New("empty configuration")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	return &Storage{
		DBStore: client.Database("grpc-server"),
	}, nil
}

// Set adds key:value pair to db.
func (s *Storage) Set(key string, val string) error {
	c := s.DBStore.Collection("map")
	ctx := context.Background()
	opts := options.Update().SetUpsert(true)

	filter := bson.D{{"Key", key}}
	update := bson.D{{"$set", bson.D{{"Key", key}, {"Value", val}}}}

	_, err := c.UpdateOne(ctx, filter, update, opts)
	return err
}

// Get gets value by key from db.
func (s *Storage) Get(key string) (string, error) {
	c := s.DBStore.Collection("map")
	filter := bson.D{{"Key", key}}

	var kv schema.KeyValue
	return kv.Value, c.FindOne(context.Background(), filter).Decode(&kv)
}
