package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

var tasks []Task
var nextID = 1

func main() {
	r := mux.NewRouter()

	// Task routes
	r.HandleFunc("/api/tasks", GetTasks).Methods("GET")
	r.HandleFunc("/api/tasks", CreateTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", GetTask).Methods("GET")
	r.HandleFunc("/api/tasks/{id}", UpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", DeleteTask).Methods("DELETE")

	// CORS middleware
	r.Use(corsMiddleware)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
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
	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    tasks,
	})
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

	for _, task := range tasks {
		if task.ID == id {
			sendJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    task,
			})
			return
		}
	}

	sendJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Task not found",
	})
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

	task.ID = nextID
	nextID++
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	
	if task.Status == "" {
		task.Status = "pending"
	}

	tasks = append(tasks, task)

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

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = updatedTask.Title
			tasks[i].Description = updatedTask.Description
			tasks[i].Status = updatedTask.Status
			tasks[i].UpdatedAt = time.Now()

			sendJSON(w, http.StatusOK, Response{
				Success: true,
				Message: "Task updated successfully",
				Data:    tasks[i],
			})
			return
		}
	}

	sendJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Task not found",
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

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			sendJSON(w, http.StatusOK, Response{
				Success: true,
				Message: "Task deleted successfully",
			})
			return
		}
	}

	sendJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Task not found",
	})
}