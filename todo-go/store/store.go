package store

import "github.com/selvamtech08/todogo/model"

type TaskStoreager interface {
	Create(model.Task) error
	GetPending() ([]*model.Task, error)
	Get(string) (*model.Task, error)
	GetAll() ([]*model.Task, error)
	Update(model.UpdateTask) error
	Remove(string) error
}
