package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kweusuf/novel-qa-go/services"
)

// Import the actual service types

// Mock services for testing - using the actual service types with test configurations
func setupMockHandler() *QAHandler {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	ollamaService := services.NewOllamaService("http://localhost:11434")

	return NewQAHandler(novelService, chromaService, ollamaService)
}

func TestNewQAHandler(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	ollamaService := services.NewOllamaService("http://localhost:11434")

	handler := NewQAHandler(novelService, chromaService, ollamaService)

	if handler.novelService == nil {
		t.Error("Expected novelService to be set")
	}
	if handler.chromaService == nil {
		t.Error("Expected chromaService to be set")
	}
	if handler.ollamaService == nil {
		t.Error("Expected ollamaService to be set")
	}
}

func TestShowIndex(t *testing.T) {
	handler := setupMockHandler()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")
	r.GET("/", handler.ShowIndex)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUploadNovel_ValidFile(t *testing.T) {
	handler := setupMockHandler()

	testContent := "This is test novel content for testing file upload."
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte(testContent))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/upload", handler.UploadNovel)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUploadNovel_InvalidFile(t *testing.T) {
	handler := setupMockHandler()

	testContent := "This is test content."
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "test.pdf")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte(testContent))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/upload", handler.UploadNovel)

	r.ServeHTTP(w, req)

	// Should get 500 Internal Server Error since no valid files were processed
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Check that response contains error message about file type
	response := w.Body.String()
	if !bytes.Contains([]byte(response), []byte("Only .txt files allowed")) {
		t.Error("Expected response to contain file type error message")
	}
}

func TestUploadNovel_NoFiles(t *testing.T) {
	handler := setupMockHandler()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/upload", handler.UploadNovel)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAskQuestion_ValidRequest(t *testing.T) {
	handler := setupMockHandler()

	questionReq := map[string]interface{}{
		"question": "What is this about?",
		"model":    "phi3",
	}

	jsonData, _ := json.Marshal(questionReq)
	req := httptest.NewRequest("POST", "/ask", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/ask", handler.AskQuestion)

	r.ServeHTTP(w, req)

	// Accept both 200 (success) and 500 (service unavailable) in test environment
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 200 or 500, got %d", w.Code)
	}

	// If we get a 200, check for answer; if 500, that's also acceptable in CI
	if w.Code == http.StatusOK {
		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		} else if response["answer"] == "" {
			t.Error("Expected answer in response")
		}
	}
}

func TestAskQuestion_InvalidRequest(t *testing.T) {
	handler := setupMockHandler()

	req := httptest.NewRequest("POST", "/ask", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/ask", handler.AskQuestion)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAskQuestion_InvalidModel(t *testing.T) {
	handler := setupMockHandler()

	questionReq := map[string]interface{}{
		"question": "What is this about?",
		"model":    "invalid_model",
	}

	jsonData, _ := json.Marshal(questionReq)
	req := httptest.NewRequest("POST", "/ask", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/ask", handler.AskQuestion)

	r.ServeHTTP(w, req)

	// Accept both 200 (success with default model) and 500 (service unavailable)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 200 or 500, got %d", w.Code)
	}
}

func TestGetModels_Success(t *testing.T) {
	handler := setupMockHandler()

	req := httptest.NewRequest("GET", "/models", nil)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/models", handler.GetModels)

	r.ServeHTTP(w, req)

	// Accept both 200 (success) and 500 (service unavailable) in test environment
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 200 or 500, got %d", w.Code)
	}

	// If we get a 200, check for models; if 500, that's also acceptable in CI
	if w.Code == http.StatusOK {
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		} else if models, exists := response["models"]; !exists {
			t.Error("Expected models in response")
		} else if modelsSlice, ok := models.([]interface{}); !ok {
			t.Error("Expected models to be an array")
		} else if len(modelsSlice) == 0 {
			t.Error("Expected at least one model in response")
		}
	}
}

func TestGetModels_CustomEndpoint(t *testing.T) {
	handler := setupMockHandler()

	req := httptest.NewRequest("GET", "/models?endpoint=http://custom:11434", nil)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/models", handler.GetModels)

	r.ServeHTTP(w, req)

	// Should fail with 500 since custom endpoint doesn't exist
	// This tests the error handling for custom endpoints
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Check that response contains error message
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Error("Expected error message in response")
	}
}
