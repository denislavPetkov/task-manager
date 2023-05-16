package task

import (
	"context"
	"fmt"

	"github.com/denislavpetkov/task-manager/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

var (
	ctx    = context.Background()
	logger *zap.Logger
)

func init() {
	logger = zap.L().Named("task-db")
}

type Mongodb interface {
	CreateCollection(name string) error

	AddTask(collection string, task model.Task) error
	UpdateTask(collection, taskTitle string, updatedtask model.Task) error

	GetTask(collection string, taskTitle string) (model.Task, error)
	GetTasks(collection string) ([]model.Task, error)

	DeleteTask(collection, taskTitle string) error
}

type mongodbInstance struct {
	mongodbDatabase *mongo.Database
}

func NewMongodb(config MongodbConfig) (Mongodb, error) {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().ApplyURI(config.GetConnectionUri()).SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to mongodb, error: %v", err))
		return &mongodbInstance{}, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to ping mongodb instance, error: %v", err))
		return &mongodbInstance{}, err
	}

	mongodbDatabase := client.Database(config.GetDatabaseName())

	logger.Info("Connected to mongodb successfully")

	return &mongodbInstance{
		mongodbDatabase: mongodbDatabase,
	}, nil
}
