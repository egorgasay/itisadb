package storage

import (
	"context"
	"errors"
	"github.com/egorgasay/grpc-storage/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type Storage struct {
	DBStore    *mongo.Database
	Mu         sync.RWMutex
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
