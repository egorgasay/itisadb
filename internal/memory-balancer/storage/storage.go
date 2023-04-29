package storage

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"itisadb/internal/memory-balancer/config"
	"itisadb/internal/schema"
)

type Storage struct {
	dbStore *mongo.Database
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
		dbStore: client.Database("grpc-server"),
	}, nil
}

// Set adds key:value pair to db.
func (s *Storage) Set(key string, val string) error {
	c := s.dbStore.Collection("map")
	ctx := context.Background()
	opts := options.Update().SetUpsert(true)

	filter := bson.D{{"Key", key}}
	update := bson.D{{"$set", bson.D{{"Key", key}, {"Value", val}}}}

	_, err := c.UpdateOne(ctx, filter, update, opts)
	return err
}

// SetUnique adds key:value pair to db and returns an error if it already exists.
func (s *Storage) SetUnique(key string, val string) error {
	c := s.dbStore.Collection("map")
	ctx := context.Background()
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"Key", key}}
	update := bson.D{{"$set", bson.D{{"Key", key}, {"Value", val}}}}
	_, err := c.UpdateOne(ctx, filter, update, opts)
	return err
}

// Get gets value by key from db.
func (s *Storage) Get(key string) (string, error) {
	c := s.dbStore.Collection("map")
	filter := bson.D{{"Key", key}}

	var kv schema.KeyValue
	return kv.Value, c.FindOne(context.Background(), filter).Decode(&kv)
}
