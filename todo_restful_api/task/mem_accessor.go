package task

import (
	"errors"
	"fmt"
)

var ErrTaskNotExist = errors.New("task does not exist")

type InMemoryAccessor struct {
	tasks  map[ID]Task
	nextID int64
}

func NewInMemoryAccessor() Accessor {
	return &InMemoryAccessor{
		tasks:  map[ID]Task{},
		nextID: int64(1),
	}
}

func (m *InMemoryAccessor) Get(id ID) (Task, error) {
	t, exists := m.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotExist
	}
	return t, nil
}

func (m InMemoryAccessor) GetAll() (map[ID]Task, error) {
	if len(m.tasks) <= 0 {
		return map[ID]Task{}, ErrTaskNotExist
	}
	var temp = make(map[ID]Task, len(m.tasks))
	for key, value := range m.tasks {
		temp[key] = value
	}
	return temp, nil
}

func (m *InMemoryAccessor) Put(id ID, t Task) error {
	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotExist
	}
	temp := m.tasks[id]

	if title := t.Title; title != "" {
		temp.Title = title
	}
	temp.Status = t.Status
	m.tasks[id] = temp
	return nil
}

func (m *InMemoryAccessor) Post(t Task) (ID, error) {
	id := ID(fmt.Sprint(m.nextID))
	m.nextID++
	m.tasks[id] = t
	return id, nil
}

func (m *InMemoryAccessor) Delete(id ID) error {
	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotExist
	}
	delete(m.tasks, id)
	return nil
}
