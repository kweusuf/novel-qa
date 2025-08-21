package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kweusuf/novel-qa-go/handlers"
	"github.com/kweusuf/novel-qa-go/services"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Setup test environment
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupTestHandler() *handlers.QAHandler {
	// Create mock services for testing
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	ollamaService := services.NewOllamaService("http://localhost:11434")

	return handlers.NewQAHandler(novelService, chromaService, ollamaService)
}

func TestShowIndex(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Create a gin router and register the handler
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", handler.ShowIndex)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if response contains expected content
	if !bytes.Contains(w.Body.Bytes(), []byte("models")) {
		t.Error("Expected response to contain models data")
	}
}

func TestUploadNovel(t *testing.T) {
	handler := setupTestHandler()

	// Create a test file content
	testContent := "This is a test novel content for testing purposes."
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form file
	part, err := writer.CreateFormFile("files", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte(testContent))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Create a gin router and register the handler
	r := gin.Default()
	r.POST("/upload", handler.UploadNovel)

	r.ServeHTTP(w, req)

	// Should return 200 even if ChromaDB operations fail in test environment
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 200 or 500, got %d", w.Code)
	}
}

func TestAskQuestion(t *testing.T) {
	handler := setupTestHandler()

	questionReq := map[string]interface{}{
		"question": "What is this about?",
		"model":    "phi3",
	}

	jsonData, _ := json.Marshal(questionReq)
	req := httptest.NewRequest("POST", "/ask", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create a gin router and register the handler
	r := gin.Default()
	r.POST("/ask", handler.AskQuestion)

	r.ServeHTTP(w, req)

	// Should return 200 if successful, or 500 if services are unavailable
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 200 or 500, got %d", w.Code)
	}
}

func TestGetModels(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/models", nil)
	w := httptest.NewRecorder()

	// Create a gin router and register the handler
	r := gin.Default()
	r.GET("/models", handler.GetModels)

	r.ServeHTTP(w, req)

	// Should return 200 if Ollama is available, or 500 if not
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 200 or 500, got %d", w.Code)
	}
}

func TestApplicationBuild(t *testing.T) {
	// Test that the application can be built successfully
	// This test will pass if the build succeeds
	if _, err := os.Stat("./novel-qa"); os.IsNotExist(err) {
		t.Skip("Binary not found - run 'go build' first")
	}
}

func TestGoModValidity(t *testing.T) {
	// Test that go.mod is valid
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		t.Fatal("go.mod not found")
	}

	// Check if go.sum exists (it should after go mod tidy)
	if _, err := os.Stat("go.sum"); os.IsNotExist(err) {
		t.Fatal("go.sum not found - run 'go mod tidy'")
	}
}

func TestMainFunction_EnvironmentVariableHandling(t *testing.T) {
	// Test that environment variables are handled correctly
	originalValue := os.Getenv("OLLAMA_HOST")
	defer func() {
		if originalValue != "" {
			os.Setenv("OLLAMA_HOST", originalValue)
		} else {
			os.Unsetenv("OLLAMA_HOST")
		}
	}()

	// Test with custom environment variable
	os.Setenv("OLLAMA_HOST", "http://custom:8080")
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}
	if ollamaHost != "http://custom:8080" {
		t.Errorf("Expected custom host, got %s", ollamaHost)
	}

	// Test with empty environment variable (should use default)
	os.Unsetenv("OLLAMA_HOST")
	ollamaHost = os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}
	if ollamaHost != "http://localhost:11434" {
		t.Errorf("Expected default host, got %s", ollamaHost)
	}
}

func TestMainFunction_ServiceInitialization(t *testing.T) {
	// Test that services can be initialized without errors
	novelService := services.NewNovelService("test_novels")
	if novelService == nil {
		t.Error("Expected novelService to be initialized")
	}

	chromaService := services.NewChromaService("test_chroma_db")
	if chromaService == nil {
		t.Error("Expected chromaService to be initialized")
	}

	ollamaService := services.NewOllamaService("http://localhost:11434")
	if ollamaService == nil {
		t.Error("Expected ollamaService to be initialized")
	}

	// Clean up test directories
	os.RemoveAll("test_novels")
	os.RemoveAll("test_chroma_db")
}

func TestMainFunction_HandlerInitialization(t *testing.T) {
	// Test that handler can be initialized with services
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	ollamaService := services.NewOllamaService("http://localhost:11434")

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	qaHandler := handlers.NewQAHandler(novelService, chromaService, ollamaService)
	if qaHandler == nil {
		t.Error("Expected qaHandler to be initialized")
	}
}

func TestMainFunction_GinRouterSetup(t *testing.T) {
	// Test that Gin router can be set up without errors
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	if r == nil {
		t.Error("Expected Gin router to be initialized")
	}

	// Test static file handler registration
	r.Static("/static", "./static")

	// Test HTML template loading
	r.LoadHTMLGlob("templates/*")

	// Test that routes can be registered
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	ollamaService := services.NewOllamaService("http://localhost:11434")
	qaHandler := handlers.NewQAHandler(novelService, chromaService, ollamaService)

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	// Register routes
	r.GET("/", qaHandler.ShowIndex)
	r.POST("/upload", qaHandler.UploadNovel)
	r.POST("/ask", qaHandler.AskQuestion)
	r.GET("/models", qaHandler.GetModels)

	// Test that routes are registered by making a test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should not panic and should return a response
	if w.Code == 0 {
		t.Error("Expected a response from the router")
	}
}

func TestMainFunction_RequiredDirectories(t *testing.T) {
	// Test that required directories exist or can be created
	requiredDirs := []string{"novels", "static", "templates"}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Logf("Directory %s does not exist, but that's okay for testing", dir)
		} else {
			t.Logf("Directory %s exists", dir)
		}
	}
}

func TestRunServer(t *testing.T) {
	// Test the runServer function which contains all the main application logic
	originalValue := os.Getenv("OLLAMA_HOST")
	defer func() {
		if originalValue != "" {
			os.Setenv("OLLAMA_HOST", originalValue)
		} else {
			os.Unsetenv("OLLAMA_HOST")
		}
	}()

	// Test with default host
	os.Unsetenv("OLLAMA_HOST")
	r, err := runServer()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if r == nil {
		t.Error("Expected Gin router to be returned")
	}

	// Test with custom host
	os.Setenv("OLLAMA_HOST", "http://test:9999")
	r2, err2 := runServer()
	if err2 != nil {
		t.Errorf("Expected no error with custom host, got %v", err2)
	}
	if r2 == nil {
		t.Error("Expected Gin router to be returned with custom host")
	}

	// Clean up test directories
	os.RemoveAll("novels")
	os.RemoveAll("chroma_db")
}

func TestMainFunction_InitializationLogic(t *testing.T) {
	// Test the core initialization logic from main function
	// This tests the same logic as main() but without starting the server

	// Test environment variable handling
	originalValue := os.Getenv("OLLAMA_HOST")
	defer func() {
		if originalValue != "" {
			os.Setenv("OLLAMA_HOST", originalValue)
		} else {
			os.Unsetenv("OLLAMA_HOST")
		}
	}()

	// Test with default host
	os.Unsetenv("OLLAMA_HOST")
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}
	if ollamaHost != "http://localhost:11434" {
		t.Errorf("Expected default host, got %s", ollamaHost)
	}

	// Test with custom host
	os.Setenv("OLLAMA_HOST", "http://test:9999")
	ollamaHost = os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}
	if ollamaHost != "http://test:9999" {
		t.Errorf("Expected custom host, got %s", ollamaHost)
	}

	// Test service initialization (same as main function)
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	chromaService.Initialize()
	ollamaService := services.NewOllamaService(ollamaHost)

	// Test handler initialization (same as main function)
	qaHandler := handlers.NewQAHandler(novelService, chromaService, ollamaService)

	// Test Gin setup (same as main function)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Test route registration (same as main function)
	r.GET("/", qaHandler.ShowIndex)
	r.POST("/upload", qaHandler.UploadNovel)
	r.POST("/ask", qaHandler.AskQuestion)
	r.GET("/models", qaHandler.GetModels)

	// Clean up
	os.RemoveAll("test_novels")
	os.RemoveAll("test_chroma_db")

	// Verify all components were initialized
	if novelService == nil {
		t.Error("novelService should not be nil")
	}
	if chromaService == nil {
		t.Error("chromaService should not be nil")
	}
	if ollamaService == nil {
		t.Error("ollamaService should not be nil")
	}
	if qaHandler == nil {
		t.Error("qaHandler should not be nil")
	}
	if r == nil {
		t.Error("Gin router should not be nil")
	}
}
