package store

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/selvamtech08/todogo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewTaskStore(collection *mongo.Collection, ctx context.Context) TaskStoreager {
	return &Task{
		collection: collection,
		ctx:        ctx,
	}
}

func (ts *Task) Create(task model.Task) error {
	task.CreatedAt = time.Now()
	task.DeadLine = time.Now().AddDate(0, 0, 3)
	result, err := ts.collection.InsertOne(ts.ctx, task)
	if err != nil {
		return err
	}

	log.Println("new task created, id:", result.InsertedID)
	return nil
}

func (ts *Task) Get(taskName string) (*model.Task, error) {
	filter := bson.D{{Key: "title", Value: taskName}}
	result := ts.collection.FindOne(ts.ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var task model.Task
	if err := result.Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (ts *Task) GetPending() ([]*model.Task, error) {
	filter := bson.D{{Key: "completed", Value: false}}
	cursor, err := ts.collection.Find(ts.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ts.ctx)

	var tasks []*model.Task

	for cursor.Next(ts.ctx) {
		var task model.Task
		err := cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if cursor.Err() != nil {
		return nil, err
	}
	return tasks, nil

}

func (ts *Task) GetAll() ([]*model.Task, error) {
	var tasks []*model.Task
	cursor, err := ts.collection.Find(ts.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ts.ctx)

	for cursor.Next(ts.ctx) {
		var task model.Task
		err := cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if cursor.Err() != nil {
		return nil, err
	}

	return tasks, nil
}

func (ts *Task) Update(task model.UpdateTask) error {
	filter := bson.M{"title": task.Title}
	changes := bson.D{}
	if task.NewTitle != nil {
		changes = append(changes, bson.E{Key: "title", Value: task.NewTitle})
	}
	if task.Remarks != nil {
		changes = append(changes, bson.E{Key: "remarks", Value: task.Remarks})
	}
	if task.DeadLine != nil {
		deadLine := time.Now().AddDate(0, 0, *task.DeadLine)
		changes = append(changes, bson.E{Key: "dead_line", Value: deadLine})
	}
	if task.Completed != nil {
		changes = append(changes, bson.E{Key: "completed", Value: task.Completed})
	}
	if task.Progress != nil {
		changes = append(changes, bson.E{Key: "progress", Value: task.Progress})
	}
	changes = append(changes, bson.E{Key: "updated_at", Value: time.Now()})
	result, err := ts.collection.UpdateOne(ts.ctx, filter, bson.M{"$set": changes})
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("no matched document found for given title")
	}

	log.Println("task updated, count:", result.ModifiedCount)
	return nil
}

func (ts *Task) Remove(taskName string) error {
	filter := bson.M{"title": taskName}
	result, err := ts.collection.DeleteOne(ts.ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no record found for given title")
	}
	log.Println("task deleted, count:", result.DeletedCount)
	return nil
}
