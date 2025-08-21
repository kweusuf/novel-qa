package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewOllamaService(t *testing.T) {
	baseURL := "http://localhost:11434"
	service := NewOllamaService(baseURL)

	if service.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, service.baseURL)
	}

	if service.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	if service.client.Timeout != 300*time.Second {
		t.Errorf("Expected timeout %v, got %v", 300*time.Second, service.client.Timeout)
	}
}

func TestOllamaService_Ask_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/chat" {
			response := `{"model":"test","created_at":"2023-01-01T00:00:00Z","message":{"role":"assistant","content":"Test response"},"done":true}`
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
		}
	}))
	defer server.Close()

	service := NewOllamaService(server.URL)
	_, err := service.Ask("test question", "test-model", "test context")

	// Should succeed with mock server
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestOllamaService_Ask_HTTPError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	service := NewOllamaService(server.URL)
	_, err := service.Ask("test question", "test-model", "test context")

	if err == nil {
		t.Error("Expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("Expected error to contain 'status 500', got %v", err)
	}
}

func TestOllamaService_Ask_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	service := NewOllamaService(server.URL)
	_, err := service.Ask("test question", "test-model", "test context")

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestOllamaService_GetModels_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			response := `{"models":[{"name":"phi3"},{"name":"llama3"},{"name":"mistral"}]}`
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
		}
	}))
	defer server.Close()

	service := NewOllamaService(server.URL)
	models, err := service.GetModels()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedModels := []string{"phi3", "llama3", "mistral"}
	if len(models) != len(expectedModels) {
		t.Errorf("Expected %d models, got %d", len(expectedModels), len(models))
	}

	for i, model := range models {
		if model != expectedModels[i] {
			t.Errorf("Expected model %s, got %s", expectedModels[i], model)
		}
	}
}

func TestOllamaService_GetModels_HTTPError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	service := NewOllamaService(server.URL)
	_, err := service.GetModels()

	if err == nil {
		t.Error("Expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "status 404") {
		t.Errorf("Expected error to contain 'status 404', got %v", err)
	}
}

func TestOllamaService_GetModels_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	service := NewOllamaService(server.URL)
	_, err := service.GetModels()

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestOllamaService_GetModels_EmptyResponse(t *testing.T) {
	// Create a mock server that returns empty models array
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"models":[]}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer server.Close()

	service := NewOllamaService(server.URL)
	models, err := service.GetModels()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 0 {
		t.Errorf("Expected 0 models, got %d", len(models))
	}
}
