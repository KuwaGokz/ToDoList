package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func doRequest(handler http.HandlerFunc, method, path, body, id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if id != "" {
		req.SetPathValue("id", id)
	}
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

func resetTData(t *testing.T) {
	t.Helper()
	Tasks = []Task{}
	nextID = 1
	dataFile = filepath.Join(t.TempDir(), "tasks.json")
}

func TestCreateTask(t *testing.T) {
	resetTData(t)
	rec := doRequest(CreateTask, "POST", "/tasks", `{"title":"ASW!@#2"}`, "")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var createdTask Task
	err := json.NewDecoder(rec.Body).Decode(&createdTask)
	if err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	if createdTask.ID != 1 {
		t.Fatalf("expected ID 1, got %d", createdTask.ID)
	}

	if createdTask.Title != "ASW!@#2" {
		t.Fatalf("expected title %q, got %q", "ASW!@#2", createdTask.Title)
	}
}

func TestCreateTask_WithMessage(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"title":"a123","message":"b456"}`, "")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var createdTask Task
	err := json.NewDecoder(rec.Body).Decode(&createdTask)
	if err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	if createdTask.ID != 1 {
		t.Fatalf("expected ID 1, got %d", createdTask.ID)
	}

	if createdTask.Title != "a123" {
		t.Fatalf("expected title %q, got %q", "a123", createdTask.Title)
	}

	if createdTask.Message != "b456" {
		t.Fatalf("expected message %q, got %q", "b456", createdTask.Message)
	}
}

func TestCreateTask_TrimsSpaces(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"title":"  a  "}`, "")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var createdTask Task
	err := json.NewDecoder(rec.Body).Decode(&createdTask)
	if err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	if createdTask.ID != 1 {
		t.Fatalf("expected ID 1, got %d", createdTask.ID)
	}

	if createdTask.Title != "a" {
		t.Fatalf("expected title %q, got %q", "a", createdTask.Title)
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"title":""}`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(Tasks))
	}
}

func TestCreateTask_SpacesOnlyTitle(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"title":"     "}`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(Tasks))
	}
}

func TestCreateTask_NoTitleField(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"message":"wdcwww1324"}`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(Tasks))
	}
}

func TestCreateTask_BadJSON(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"title":"wdcwww1`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(Tasks))
	}
}

func TestCreateTask_TitleTooLong(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"title":"`+strings.Repeat("w", MaxTitleLen+1)+`"}`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(Tasks))
	}
}

func TestCreateTask_MessageTooLong(t *testing.T) {
	resetTData(t)

	rec := doRequest(CreateTask, "POST", "/tasks", `{"message":`+strings.Repeat("w", MaxMessageLen), "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(Tasks))
	}
}
