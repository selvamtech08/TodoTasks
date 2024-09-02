package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Remarks   string             `json:"remarks" bson:"remarks"`
	Completed bool               `json:"completed" bson:"completed"`
	Progress  int                `json:"progress" bson:"progress"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at"`
	DeadLine  time.Time          `json:"dead_line" bson:"dead_line"`
}

type UpdateTask struct {
	Title     *string `json:"title" bson:"title,omitempty"`
	NewTitle  *string `json:"new_title" bson:"title,omitempty"`
	Remarks   *string `json:"remarks" bson:"remarks,omitempty"`
	Completed *bool   `json:"completed" bson:"completed,omitempty"`
	Progress  *int    `json:"progress" bson:"progress,omitempty"`
	DeadLine  *int    `json:"dead_line" bson:"dead_line,omitempty"`
}
