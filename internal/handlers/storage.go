package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

var dataFile = "tasksList.json"

func saveTasks() error {
	jsonTaskList, err := json.MarshalIndent(Tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, jsonTaskList, 0644)
}

func loadTasks() {
	tasksData, err := os.ReadFile(dataFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		log.Fatal("Cannot read saved files:", err)
		return
	}

	err = json.Unmarshal(tasksData, &Tasks)
	for _, t := range Tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
}
