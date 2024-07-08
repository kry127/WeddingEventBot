package mongo

import (
	"context"
	"time"

	"kry127.ru/weddingbot/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClientFromConfig(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// opts := new(options.ClientOptions).ApplyURI("mongodb://158.160.141.41:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.2.10&tls=false")
	opts := new(options.ClientOptions).ApplyURI(cfg.MongoDBConnString).SetAuth(options.Credential{
		Username: cfg.User,
		Password: cfg.Password,
	})
	return mongo.Connect(ctx, opts)
}
