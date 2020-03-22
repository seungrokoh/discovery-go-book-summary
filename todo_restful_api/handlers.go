package main

import (
	"encoding/json"
	"errors"
	"github.com/jaeyeom/gogo/task"
	"log"
	"net/http"
	"text/template"
)

var m = NewMemoryDataAccess()

const pathPrefix = "/api/v1/task/"

func apiHandler(w http.ResponseWriter, r *http.Request) {
	getID := func() (ID, error) {
		id := task.ID(r.URL.Path[len(pathPrefix):])
		if id == "" {
			return ID(id), errors.New("apiHandler: ID is empty")
		}
		return ID(id), nil
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
		id, err := getID()
		if err != nil {
			log.Println(err)
			return
		}
		t, err := m.Get(id)
		err = json.NewEncoder(w).Encode(Response{
			ID:    id,
			Task:  t,
			Error: ResponseError{err},
		})
		if err != nil {
			log.Println(err)
		}
	case "PUT":
		id, err := getID()
		if err != nil {
			log.Println(err)
			return
		}
		tasks, err := getTasks()
		if err != nil {
			log.Println(err)
			return
		}
		for _, t := range tasks {
			err = m.Put(id, t)
			err = json.NewEncoder(w).Encode(Response{
				ID:    id,
				Task:  t,
				Error: ResponseError{err},
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
	case "POST":
		tasks, err := getTasks()
		if err != nil {
			log.Println(err)
			return
		}
		for _, t := range tasks {
			id, err := m.Post(t)
			err = json.NewEncoder(w).Encode(Response{
				ID:    id,
				Task:  t,
				Error: ResponseError{err},
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
	case "DELETE":
		id, err := getID()
		if err != nil {
			log.Println(err)
			return
		}
		err = m.Delete(id)
		err = json.NewEncoder(w).Encode(Response{
			ID:    id,
			Error: ResponseError{err},
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// htmlHandler
const htmlPrefix = "/task/"

var tmpl = template.Must(template.ParseGlob("html/*html"))

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println(r.Method, "method is not supported")
		return
	}
	getID := func() (ID, error) {
		id := task.ID(r.URL.Path[len(htmlPrefix):])
		if id == "" {
			return ID(id), errors.New("htmlHandler: ID is empty")
		}
		return ID(id), nil
	}
	id, err := getID()

	if err != nil {
		log.Println(err)
		return
	}
	t, err := m.Get(id)
	err = tmpl.ExecuteTemplate(w, "task.html", &Response {
		ID: id,
		Task: t,
		Error: ResponseError{err},
	})
	if err != nil {
		log.Println(err)
		return
	}
}
