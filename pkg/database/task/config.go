package task

import (
	"fmt"
	"os"
)

const (
	MONGODB_CONNECTION_URI = "MONGODB_CONNECTION_URI"
	MONGODB_DATABASE_NAME  = "MONGODB_DATABASE_NAME"
)

type MongodbConfig interface {
	GetConnectionUri() string
	GetDatabaseName() string
}

type mongodbConfig struct {
	connectionUri string
	databaseName  string
}

func NewMongodbConfig() (MongodbConfig, error) {
	connectionUri, ok := os.LookupEnv(MONGODB_CONNECTION_URI)
	if !ok {
		logger.Error(fmt.Sprintf("Env var %s missing", MONGODB_CONNECTION_URI))
		return mongodbConfig{}, fmt.Errorf("env var %s missing", MONGODB_CONNECTION_URI)
	}

	databaseName, ok := os.LookupEnv(MONGODB_DATABASE_NAME)
	if !ok {
		logger.Error(fmt.Sprintf("Env var %s missing", MONGODB_DATABASE_NAME))
		return mongodbConfig{}, fmt.Errorf("env var %s missing", MONGODB_DATABASE_NAME)
	}

	return mongodbConfig{
		connectionUri: connectionUri,
		databaseName:  databaseName,
	}, nil
}

func (c mongodbConfig) GetConnectionUri() string {
	return c.connectionUri
}

func (c mongodbConfig) GetDatabaseName() string {
	return c.databaseName
}
