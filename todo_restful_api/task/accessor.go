package task

type ID string

type Accessor interface {
	Get(id ID) (Task, error)
	GetAll() (List, error)
	Put(id ID, t Task) error
	Post(t *Task) (ID, error)
	Delete(id ID) error
}
