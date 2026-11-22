package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func reset() {
	mu.Lock()
	tasks, nextID = nil, 1
	mu.Unlock()
}

func TestConcurrentCreate(t *testing.T) {
	reset()
	r := newRouter()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/api/tasks",
				strings.NewReader(`{"title":"x"}`))
			r.ServeHTTP(httptest.NewRecorder(), req)
		}()
	}
	wg.Wait()

	if len(tasks) != 50 {
		t.Fatalf("oczekiwano 50 zadań, jest %d", len(tasks))
	}
	seen := map[int]bool{}
	for _, task := range tasks {
		if seen[task.ID] {
			t.Fatalf("zduplikowane ID: %d", task.ID)
		}
		seen[task.ID] = true
	}
}
func TestCreateTask(t *testing.T) {
	reset()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks",
		strings.NewReader(`{"title":"Kup mleko"}`))
	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, chcę 201", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"pending"`) {
		t.Errorf("brak domyślnego statusu: %s", rec.Body)
	}
}

func TestCreateTaskWithoutTitle(t *testing.T) {
	reset()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{}`))
	newRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, chcę 400", rec.Code)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	reset()
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/tasks/999", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, chcę 404", rec.Code)
	}
}

func TestGetTaskInvalidID(t *testing.T) {
	reset()
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/tasks/abc", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, chcę 400", rec.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	reset()
	r := newRouter()
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/tasks",
		strings.NewReader(`{"title":"stary"}`)))

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/tasks/1",
		strings.NewReader(`{"title":"nowy","status":"completed"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, chcę 200", rec.Code)
	}

	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil))
	if !strings.Contains(rec.Body.String(), `"title":"nowy"`) {
		t.Fatalf("update się nie zapisał: %s", rec.Body)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	reset()
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/tasks/999",
		strings.NewReader(`{"title":"x","status":"pending"}`)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, chcę 404", rec.Code)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	reset()
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/tasks/999", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, chcę 404", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "deleted") {
		t.Fatalf("podwójna odpowiedź: %s", rec.Body)
	}
}

func TestUpdateTaskValidation(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"pusty tytuł", `{"title":"","status":"pending"}`},
		{"same spacje", `{"title":"   ","status":"pending"}`},
		{"brak statusu", `{"title":"x"}`},
		{"zły status", `{"title":"x","status":"done"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			r := newRouter()
			r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/tasks",
				strings.NewReader(`{"title":"oryginał"}`)))

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/tasks/1",
				strings.NewReader(tc.body)))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, chcę 400", rec.Code)
			}

			rec = httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil))
			if !strings.Contains(rec.Body.String(), `"title":"oryginał"`) {
				t.Fatalf("zadanie zostało zmienione mimo błędu: %s", rec.Body)
			}
		})
	}
}

func TestCreateTaskInvalidStatus(t *testing.T) {
	reset()
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/tasks",
		strings.NewReader(`{"title":"x","status":"done"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, chcę 400", rec.Code)
	}
}
