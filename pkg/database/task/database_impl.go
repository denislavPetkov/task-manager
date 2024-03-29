package task

import (
	"fmt"
	"regexp"

	"github.com/denislavpetkov/task-manager/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// used as a primary key
	taskTitlePrimaryKey = "title"
)

func (m *mongodbInstance) CreateCollection(name string) error {
	// The name of the collection to create. See Naming Restrictions.
	err := m.mongodbDatabase.CreateCollection(ctx, name)
	if err != nil {
		if match, _ := regexp.MatchString(`Collection .* already exists`, err.Error()); match {
			return nil
		}
		logger.Error("Failed to create a new collection")
		return err
	}

	// set task title as a unique index
	_, err = m.mongodbDatabase.Collection(name).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: taskTitlePrimaryKey, Value: -1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Error("Failed to set the title as a unique index to the new collection")
		return err
	}

	logger.Info("Created a new collection successfully")

	return err
}

func (m *mongodbInstance) AddTask(collection string, task model.Task) error {
	dbCollection := m.mongodbDatabase.Collection(collection)

	_, err := dbCollection.InsertOne(ctx, task)
	if err != nil {
		logger.Error("Failed to insert task to the collection")
		return fmt.Errorf("failed to add task %v to the %s collection, error: %w", task, collection, err)
	}

	logger.Info("Inserted a new task to the collection successfully")

	return nil
}

func (m *mongodbInstance) UpdateTask(collection, taskTitle string, updatedtask model.Task) error {
	dbCollection := m.mongodbDatabase.Collection(collection)

	result, err := dbCollection.ReplaceOne(ctx, bson.D{
		{
			Key:   taskTitlePrimaryKey,
			Value: taskTitle,
		},
	},
		updatedtask,
	)
	if err != nil {
		logger.Error("Failed to update task in the collection")
		return fmt.Errorf("failed to update task '%s' in %s collection, error: %w", taskTitle, collection, err)
	}

	if result.MatchedCount == 0 {
		logger.Error("Task not found in the collection")
		return mongo.ErrNoDocuments
	}

	logger.Info("Updated a task in the collection successfully")

	return nil
}

func (m *mongodbInstance) GetTask(collection, taskTitle string) (model.Task, error) {
	dbCollection := m.mongodbDatabase.Collection(collection)

	result := dbCollection.FindOne(ctx, bson.D{
		{
			Key:   taskTitlePrimaryKey,
			Value: taskTitle,
		},
	},
	)

	task := model.Task{}

	err := result.Decode(&task)
	if err != nil {
		logger.Error("Failed to decode db task into object")
		return model.Task{}, fmt.Errorf("failed to decode task into object, error: %w", err)
	}

	logger.Info("Got a task from the collection successfully")

	return task, nil
}

func (m *mongodbInstance) GetTasks(collection string) ([]model.Task, error) {
	dbCollection := m.mongodbDatabase.Collection(collection)

	cursor, err := dbCollection.Find(ctx, bson.D{})
	if err != nil {
		logger.Error("Failed to get tasks from the collection")
		return []model.Task{}, fmt.Errorf("failed to get tasks from %s collection, error: %w", collection, err)
	}

	var tasks []model.Task

	err = cursor.All(ctx, &tasks)
	if err != nil {
		logger.Error("Failed to get tasks from the collection")
		return []model.Task{}, fmt.Errorf("failed to get tasks from %s collection, error: %w", collection, err)
	}

	logger.Info("Got tasks from the collection successfully")

	return tasks, nil
}

func (m *mongodbInstance) DeleteTask(collection, taskTitle string) error {
	dbCollection := m.mongodbDatabase.Collection(collection)

	_, err := dbCollection.DeleteOne(ctx,
		bson.D{
			{
				Key:   taskTitlePrimaryKey,
				Value: taskTitle,
			},
		},
	)
	if err != nil {
		logger.Error("Failed to get delete task from the collection")
		return fmt.Errorf("failed to delete task '%s' from %s collection, error: %w", taskTitle, collection, err)
	}

	logger.Info("Delete task from the collection successfully")

	return nil
}
