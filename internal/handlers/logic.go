package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newTask.Title = strings.TrimSpace(newTask.Title)
	newTask.Message = strings.TrimSpace(newTask.Message)

	if newTask.Title == "" {
		http.Error(w, "Title is empty", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(newTask.Title) > MaxTitleLen {
		http.Error(w, "Title is too long", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(newTask.Message) > MaxMessageLen {
		http.Error(w, "Message is too long", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newTask.ID = nextID
	nextID++
	newTask.Status = false
	Tasks = append(Tasks, newTask)

	if err := saveTasks(); err != nil {
		mu.Unlock()
		http.Error(w, "Task saving failed", http.StatusInternalServerError)
		return
	}
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
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
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
			if err := saveTasks(); err != nil {
				http.Error(w, "task saving failed", http.StatusInternalServerError)
				return
			}
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
			if err := saveTasks(); err != nil {
				http.Error(w, "Task saving failed", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)

}
