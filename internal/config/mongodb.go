package config

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(viper *viper.Viper, uriKey string) *mongo.Client {
	uri := viper.GetString(uriKey)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping MongoDB: %v\n", err)
		os.Exit(1)
	}

	return client
}
