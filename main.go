package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var (
	mu     sync.RWMutex
	tasks  []Task
	nextID = 1
)

func main() {
	r := newRouter()
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/tasks", GetTasks).Methods("GET")
	r.HandleFunc("/api/tasks", CreateTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", GetTask).Methods("GET")
	r.HandleFunc("/api/tasks/{id}", UpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", DeleteTask).Methods("DELETE")
	r.Use(corsMiddleware)
	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ustaw nagłówki CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Obsługa preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	snapshot := make([]Task, len(tasks))
	copy(snapshot, tasks)
	mu.RUnlock()

	sendJSON(w, http.StatusOK, Response{Success: true, Data: snapshot})
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid task ID",
		})
		return
	}

	mu.RLock()
	var found *Task
	for _, t := range tasks {
		if t.ID == id {
			t := t
			found = &t
			break
		}
	}
	mu.RUnlock()

	if found == nil {
		sendJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Task not found",
		})
		return
	}
	sendJSON(w, http.StatusOK, Response{Success: true, Data: *found})
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if task.Title == "" {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Title is required",
		})
		return
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	if task.Status == "" {
		task.Status = "pending"
	}

	mu.Lock()
	task.ID = nextID
	nextID++
	tasks = append(tasks, task)
	mu.Unlock()

	sendJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: "Task created successfully",
		Data:    task,
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid task ID",
		})
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	mu.Lock()
	var found *Task
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = updatedTask.Title
			tasks[i].Description = updatedTask.Description
			tasks[i].Status = updatedTask.Status
			tasks[i].UpdatedAt = time.Now()
			t := tasks[i] // kopia do wysłania po zdjęciu blokady
			found = &t
			break
		}
	}
	mu.Unlock()

	if found == nil {
		sendJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Task not found",
		})
		return
	}
	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Task updated successfully",
		Data:    *found,
	})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid task ID",
		})
		return
	}

	mu.Lock()
	var found *Task
	for i, task := range tasks {
		if task.ID == id {
			task := tasks[i]
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = &task
			break
		}
	}
	mu.Unlock()
	if found == nil {
		sendJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Task not found",
		})
		return
	}
	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Task deleted successfully",
	})
}
