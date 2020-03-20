package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jaeyeom/gogo/task"
	"log"
	"net/http"
)

type ID int64

type DataAccess interface {
	Get(id ID) (task.Task, error)
	Put(id ID, t task.Task) error
	Post(t task.Task) (ID, error)
	Delete(id ID) error
}

var ErrTaskNotExist = errors.New("task does not exist")

type MemoryDataAccess struct {
	tasks  map[ID]task.Task
	nextID int64
}

func (m *MemoryDataAccess) Get(id ID) (task.Task, error) {
	t, exists := m.tasks[id]
	if !exists {
		return task.Task{}, ErrTaskNotExist
	}
	return t, nil
}

func (m *MemoryDataAccess) Put(id ID, t task.Task) error {
	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotExist
	}
	m.tasks[id] = t
	return nil
}

func (m *MemoryDataAccess) Post(t task.Task) (ID, error) {
	id := ID(fmt.Sprint(m.nextID))
	m.nextID++
	m.tasks[id] = t
	return id, nil
}

func (m *MemoryDataAccess) Delete(id ID) error {
	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotExist
	}
	delete(m.tasks, id)
	return nil
}

func NewMemoryDataAccess() DataAccess {
	return &MemoryDataAccess{
		tasks:  map[ID]task.Task{},
		nextID: int64(1),
	}
}

type ResponseError struct {
	Err error
}

func (err ResponseError) MarshalJSON() ([]byte, error) {
	if err.Err == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%v\"", err.Err)), nil
}

func (err *ResponseError) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v == nil {
		err.Err = nil
		return nil
	}
	switch tv := v.(type) {
	case string:
		if tv == ErrTaskNotExist.Error() {
			err.Err = ErrTaskNotExist
			return nil
		}
		err.Err = errors.New(tv)
		return nil
	default:
		return errors.New("ResponseError unmarshal failed")
	}
}

type Response struct {
	ID    ID            `json:"id,omitempty"`
	Task  task.Task     `json:"task"`
	Error ResponseError `json:"error"`
}

var m = NewMemoryDataAccess()

const pathPrefix = "/api/v1/task/"

func apiHandler(w http.ResponseWriter, r *http.Request) {
	getID := func() (ID, error) {
		id := task.ID(r.URL.Path[len(pathPrefix):])
		if id == "" {
			return id, errors.New("apiHandler: ID is empty")
		}
		return id, nil
	}
	getTasks := func() ([]task.Task, error) {
		var result []task.Task
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		encodedTasks, ok := r.PostForm["task"]
		if !ok {
			return nil, errors.New("task parameter expected")
		}
		for _, encodedTasks := range encodedTasks {
			var t task.Task
			if err := json.Unmarshal([]byte(encodedTasks), &t); err != nil {
				return nil, err
			}
			result = append(result, t)
		}
		return result, nil
	}

	switch r.Method {
	case "GET":
		panic("unimplemented")
	case "PUT":
		panic("unimplemented")
	case "POST":
		panic("unimplemented")
	case "DELETE":
		panic("unimplemented")
	}
}

func main() {
	http.HandleFunc(pathPrefix, apiHandler)
	log.Fatal(http.ListenAndServe(":8887", nil))
}
