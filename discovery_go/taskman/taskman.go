package main

import (
	"github.com/jaeyeom/gogo/task"
)

type ID string

type DataAccess interface {
	Get(id ID) (task.Task, error)
	Put(id ID, t task.Task) error
	Post(t task.Task) (ID, error)
	Delete(id ID) error
}

type MemoryDataAccess struct {
	tasks map[ID]task.Task
	nextID int64
}

func NewMemoryDataAccess() DataAccess {
	return &MemoryDataAccess{
		tasks: map[ID]task.Task{},
		nextID: int64(1),
	}
}


