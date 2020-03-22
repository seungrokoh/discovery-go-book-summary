package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"text/template"
)

var m = NewInMemoryAccessor()


func apiHandler(w http.ResponseWriter, r *http.Request) {
	getID := func() (ID, error) {
		id := ID(r.URL.Path[len(apiPathPrefix):])
		if id == "" {
			return ID(id), errors.New("apiHandler: ID is empty")
		}
		return ID(id), nil
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

var tmpl = template.Must(template.ParseGlob("html/*html"))

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println(r.Method, "method is not supported")
		return
	}
	getID := func() (ID, error) {
		id := ID(r.URL.Path[len(htmlPathPrefix):])
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
	err = tmpl.ExecuteTemplate(w, "html", &Response {
		ID: id,
		Task: t,
		Error: ResponseError{err},
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func getTasks(r *http.Request) ([]Task, error) {
	var result []Task
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	encodedTasks, ok := r.PostForm["task"]
	if !ok {
		return nil, errors.New("task parameter expected")
	}
	for _, encodedTasks := range encodedTasks {
		var t Task
		if err := json.Unmarshal([]byte(encodedTasks), &t); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}