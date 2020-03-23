package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
	"todo/task"
)

var m = task.NewInMemoryAccessor()

var tmpl = template.Must(template.ParseGlob("html/*.html"))

func getResponseList(f func() (task.List, error)) ([]Response, error) {
	tasks, err := f()
	if err != nil {
		log.Println(err)
		return []Response{}, err
	}

	var responseList []Response
	for _, t := range tasks {
		responseList = append(responseList, createResponse(t.ID, t, err))
	}
	return responseList, nil
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	getID := func() (task.ID, error) {
		id := task.ID(r.URL.Path[len(htmlPathPrefix):])
		if id == "" {
			return id, errors.New("htmlHandler: ID is empty")
		}
		return id, nil
	}
	id, err := getID()
	if err != nil {
		// id를 url에 넣지 않았을 경우
		responseList, err := getResponseList(m.GetAll)
		if err != nil {
			err = tmpl.ExecuteTemplate(w, "task.html", &Response{
				ID:    "",
				Task:  task.Task{},
				Error: ResponseError{err},
			})
			log.Println(err)
			return
		}
		for _, r := range responseList {
			err = tmpl.ExecuteTemplate(w, "task.html", r)
			if err != nil {
				log.Println(err)
				return
			}
		}
	} else {
		t, err := m.Get(id)
		err = tmpl.ExecuteTemplate(w, "task.html", &Response{
			ID:    id,
			Task:  t,
			Error: ResponseError{err},
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func getTasks(r *http.Request) (task.List, error) {
	var result task.List
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	// parameter로 넘어온 task를 가져온다.
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
	log.Println(result)
	return result, nil
}

func createResponse(id task.ID, t task.Task, err error) Response {
	return Response{
		ID:    id,
		Task:  t,
		Error: ResponseError{err},
	}
}

func apiGetHandler(w http.ResponseWriter, r *http.Request) {
	id := task.ID(mux.Vars(r)["id"])
	t, err := m.Get(id)
	err = json.NewEncoder(w).Encode(createResponse(id, t, err))
	if err != nil {
		log.Println(err)
	}
}

func apiAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	responseList, err := getResponseList(m.GetAll)
	if err != nil {
		err = json.NewEncoder(w).Encode(err)
		return
	}
	err = json.NewEncoder(w).Encode(responseList)
	fmt.Println(responseList)
	if err != nil {
		log.Println(err)
	}
}

func apiPutHandler(w http.ResponseWriter, r *http.Request) {
	id := task.ID(mux.Vars(r)["id"])
	tasks, err := getTasks(r)
	if err != nil {
		log.Println(err)
		return
	}
	for _, t := range tasks {
		err = m.Put(id, t)
		err = json.NewEncoder(w).Encode(createResponse(id, t, err))

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func apiPostHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := getTasks(r)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(len(tasks))

	for _, t := range tasks {
		id, err := m.Post(&t)
		err = json.NewEncoder(w).Encode(createResponse(id, t, err))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func apiDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := task.ID(mux.Vars(r)["id"])
	err := m.Delete(id)
	err = json.NewEncoder(w).Encode(createResponse(id, task.Task{}, err))
	if err != nil {
		log.Println(err)
		return
	}
}
