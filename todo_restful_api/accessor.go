package main

import (
	"errors"
	"github.com/jaeyeom/gogo/task"
)

type ID string

type DataAccess interface {
	Get(id ID) (task.Task, error)
	Put(id ID, t task.Task) error
	Post(t task.Task) (ID, error)
	Delete(id ID) error
}

var ErrTaskNotExist = errors.New("task does not exist")
