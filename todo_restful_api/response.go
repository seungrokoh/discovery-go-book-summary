package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"todo/task"
)

var ErrTaskNotExist = errors.New("task does not exist")

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
	ID    task.ID       `json:"id,omitempty"`
	Task  task.Task     `json:"task"`
	Error ResponseError `json:"error"`
}

type ResponseList []Response

func (r ResponseList) Len() int {
	return len(r)
}

func (r ResponseList) Less(i, j int) bool {
	return string(r[i].ID) < string(r[j].ID)
}

func (r ResponseList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
