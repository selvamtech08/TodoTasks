package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/selvamtech08/todogo/model"
	"github.com/selvamtech08/todogo/store"
)

type TaskHandler struct {
	store store.TaskStoreager
}

func errResult(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	data := map[string]any{"error": err.Error()}
	_ = json.NewEncoder(w).Encode(&data)
}

func successResult(w http.ResponseWriter, code int, msg any) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(msg)
}

func NewTaskController(store store.TaskStoreager) TaskHandler {
	return TaskHandler{
		store: store,
	}
}

func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("create route...")
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errResult(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	err := th.store.Create(task)
	if err != nil {
		errResult(w, http.StatusInternalServerError, err)
		return
	}
	successResult(w, http.StatusCreated, "task created successfully")
}

func (th *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("get route...")
	taskName := r.PathValue("name")
	taskName = strings.TrimSpace(taskName)
	if taskName == "" {
		errResult(w, http.StatusBadRequest, errors.New("empty task name given"))
		return
	}

	task, err := th.store.Get(taskName)
	if err != nil {
		errResult(w, http.StatusInternalServerError, err)
		return
	}

	successResult(w, http.StatusOK, task)

}

func (th *TaskHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	log.Println("getpending route...")
	tasks, err := th.store.GetPending()
	if err != nil {
		errResult(w, http.StatusBadRequest, err)
		return
	}
	if len(tasks) == 0 {
		successResult(w, http.StatusOK, map[string]string{"info": "no pending task found"})
		return
	}
	successResult(w, http.StatusOK, tasks)
}

func (th *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("getall route...")
	tasks, err := th.store.GetAll()
	// return error if db has error
	if err != nil {
		errResult(w, http.StatusInternalServerError, err)
		return
	}
	// return all the tasks as success response
	successResult(w, http.StatusOK, tasks)
}

func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("update route...")
	var updateTask model.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&updateTask); err != nil {
		errResult(w, http.StatusBadRequest, err)
		return
	}

	if err := updateTask.Title; err == nil {
		errResult(w, http.StatusInternalServerError, errors.New("title should be given to find task"))
		return
	}
	if err := th.store.Update(updateTask); err != nil {
		errResult(w, http.StatusInternalServerError, err)
		return
	}
	successResult(w, http.StatusAccepted, "task updated successfully!")
}

func (th *TaskHandler) Remove(w http.ResponseWriter, r *http.Request) {
	log.Println("remove route...")
	taskName := r.PathValue("name")
	taskName = strings.TrimSpace(taskName)
	if taskName == "" {
		errResult(w, http.StatusBadRequest, errors.New("title should a valid string, not empty"))
		return
	}

	if err := th.store.Remove(taskName); err != nil {
		errResult(w, http.StatusInternalServerError, err)
		return
	}
	successResult(w, http.StatusOK, "task removed")
}
