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

func TestGetTask_Empty(t *testing.T) {
	resetTData(t)

	rec := doRequest(GetTask, "GET", "/tasks", "", "")

	var got []Task
	err := json.NewDecoder(rec.Body).Decode(&got)
	if err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}
	if got == nil {
		t.Fatal("expected empty array, got null")
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(got))
	}
}

func TestGetTask_List(t *testing.T) {
	resetTData(t)

	Tasks = []Task{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}, {ID: 3, Title: "c"}}

	rec := doRequest(GetTask, "GET", "/tasks", "", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got []Task
	err := json.NewDecoder(rec.Body).Decode(&got)
	if err != nil {
		t.Fatalf("cannot decode respone %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(got))
	}

	if got[0].Title != "a" {
		t.Fatalf("expected a")
	}

	if got[1].Title != "b" {
		t.Fatalf("expected b")
	}

	if got[2].Title != "c" {
		t.Fatalf("expected c")
	}
}

func TestReworkTask_SetDone(t *testing.T) {
	resetTData(t)

	Tasks = []Task{{ID: 1, Title: "a"}}

	rec := doRequest(ReworkTask, "PATCH", "/tasks/1", `{"status":true}`, "1")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if Tasks[0].Status != true {
		t.Fatalf("expected status true, got %v", Tasks[0].Status)
	}
}

func TestReworkTask_KeepsOtherFields(t *testing.T) {
	resetTData(t)

	Tasks = []Task{{ID: 1, Title: "a", Message: "b", Status: false}}

	rec := doRequest(ReworkTask, "PATCH", "/tasks/1", `{"status":true}`, "1")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if Tasks[0].Status != true {
		t.Fatalf("expected true, got %v", Tasks[0].Status)
	}

	if Tasks[0].Title != "a" {
		t.Fatalf("expected a, got %v", Tasks[0].Title)
	}

	if Tasks[0].Message != "b" {
		t.Fatalf("expected b, got %v", Tasks[0].Message)
	}
}

func TestReworkTask_NotFound(t *testing.T) {
	resetTData(t)

	Tasks = []Task{{ID: 1, Title: "a", Message: "b", Status: false}}

	rec := doRequest(ReworkTask, "PATCH", "/tasks/1", `{"status":true}`, "943")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	if Tasks[0].Status == true {
		t.Fatalf("expected false, got %v", Tasks[0].Status)
	}
}

func TestDelTask_RemovesCorrectOne(t *testing.T) {
	resetTData(t)

	Tasks = []Task{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}, {ID: 3, Title: "c"}}

	rec := doRequest(DelTask, "DELETE", "/tasks/2", "", "2")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	if len(Tasks) != 2 {
		t.Fatalf("expected >2, got %v", Tasks)
	}

	if Tasks[0].ID != 1 {
		t.Fatalf("expected Tasks[0].ID 1, got %d", Tasks[0].ID)
	}

	if Tasks[1].ID != 3 {
		t.Fatalf("expected Tasks[1].ID 3, got %d", Tasks[1].ID)
	}
}
