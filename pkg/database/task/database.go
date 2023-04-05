package task

import (
	"context"

	"github.com/denislavpetkov/task-manager/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx = context.Background()
)

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
		return &mongodbInstance{}, nil
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return &mongodbInstance{}, nil
	}

	mongodbDatabase := client.Database(config.GetDatabaseName())

	return &mongodbInstance{
		mongodbDatabase: mongodbDatabase,
	}, nil
}
