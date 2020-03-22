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

func (m *InMemoryAccessor) GetAll() ([]Task, error) {
	if len(m.tasks) <= 0 {
		return []Task{}, ErrTaskNotExist
	}
	var result []Task
	for id := range m.tasks {
		task, _ := m.Get(id)
		result = append(result, task)
	}
	return result, nil
}

func (m *InMemoryAccessor) Put(id ID, t Task) error {
	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotExist
	}
	m.tasks[id] = t
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
