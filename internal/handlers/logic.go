package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newTask.ID = nextID
	nextID++
	newTask.Status = false
	Tasks = append(Tasks, newTask)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Tasks)
}

func ReworkTask(w http.ResponseWriter, r *http.Request) {
	var newStatus TaskStatus
	err := json.NewDecoder(r.Body).Decode(&newStatus)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i := range Tasks {
		if Tasks[i].ID == id {
			Tasks[i].Status = newStatus.Status
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func DelTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i := range Tasks {
		if Tasks[i].ID == id {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)

}
