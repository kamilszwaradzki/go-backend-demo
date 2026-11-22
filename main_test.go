package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// reset czyści stan globalny przed każdym testem
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
		strings.NewReader(`{"title":"nowy","status":"done"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, chcę 200", rec.Code)
	}

	// sprawdź, że zmiana NAPRAWDĘ się zapisała
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
		strings.NewReader(`{"title":"x"}`)))
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
