package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Task represents a to-do item with an ID, Title, Description, and Status
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // "pending" or "completed"
}

var tasks []Task  // In-memory task storage
var idCounter = 1 // Task ID counter

// CreateTask handles the creation of a new task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	json.NewDecoder(r.Body).Decode(&newTask)
	newTask.ID = idCounter
	idCounter++
	newTask.Status = "pending" // Default status
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

// GetTasks returns all the tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID retrieves a task by its ID
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Get parameters from the URL
	id, _ := strconv.Atoi(params["id"])

	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// UpdateTask updates an existing task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, task := range tasks {
		if task.ID == id {
			var updatedTask Task
			json.NewDecoder(r.Body).Decode(&updatedTask)
			updatedTask.ID = task.ID
			tasks[i] = updatedTask

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// DeleteTask removes a task by its ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...) // Remove task from slice
			w.WriteHeader(http.StatusNoContent)       // Return 204 No Content
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	// Route Handlers
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", GetTaskByID).Methods("GET")
	router.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")

	fmt.Println("Server is running on port 8089")
	log.Fatal(http.ListenAndServe(":8089", router))
}
