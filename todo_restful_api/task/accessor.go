package task

import (
	"errors"
)

var ErrTaskNotExist = errors.New("task does not exist")

type ID string

type Accessor interface {
	Get(id ID) (Task, error)
	Put(id ID, t Task) error
	Post(t Task) (ID, error)
	Delete(id ID) error
}
