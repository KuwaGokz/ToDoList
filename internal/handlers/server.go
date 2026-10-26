package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
	ID      int    `json:"id"`
}

type TaskStatus struct {
	Status bool `json:"status"`
}

const (
	MaxTitleLen   = 120
	MaxMessageLen = 640
)

var Tasks = []Task{}
var mu sync.Mutex
var nextID = 1

func Server() {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler:           mux,
		Addr:              "127.0.0.1:8080",
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	mux.HandleFunc("/", BasicPage)
	mux.HandleFunc("GET /tasks", GetTask)
	mux.HandleFunc("POST /tasks", CreateTask)
	mux.HandleFunc("PATCH /tasks/{id}", ReworkTask)
	mux.HandleFunc("DELETE /tasks/{id}", DelTask)

	log.Fatal(server.ListenAndServe())
}
