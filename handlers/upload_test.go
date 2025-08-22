package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kweusuf/novel-qa-go/services"

	"github.com/gin-gonic/gin"
)

func TestNewUploadHandler(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")

	handler := NewUploadHandler(novelService, chromaService)

	if handler.novelService == nil {
		t.Error("Expected novelService to be set")
	}
	if handler.chromaService == nil {
		t.Error("Expected chromaService to be set")
	}

	// Clean up
	os.RemoveAll("test_novels")
	os.RemoveAll("test_chroma_db")
}

func TestUploadHandler_UploadNovel_ValidFile(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	handler := NewUploadHandler(novelService, chromaService)

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	// Create test file content
	testContent := "This is test novel content for testing file upload."
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.txt")
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

func TestUploadHandler_UploadNovel_InvalidFileType(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	handler := NewUploadHandler(novelService, chromaService)

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	// Create test file content with invalid extension
	testContent := "This is test content."
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.pdf")
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

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	response := w.Body.String()
	if !bytes.Contains([]byte(response), []byte("Only .txt and .epub files are supported")) {
		t.Error("Expected response to contain file type error message")
	}
}

func TestUploadHandler_UploadNovel_NoFile(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	handler := NewUploadHandler(novelService, chromaService)

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	req := httptest.NewRequest("POST", "/upload", nil)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/upload", handler.UploadNovel)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	response := w.Body.String()
	if !bytes.Contains([]byte(response), []byte("Failed to get file")) {
		t.Error("Expected response to contain file error message")
	}
}

func TestUploadHandler_UploadNovel_LargeFile(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	handler := NewUploadHandler(novelService, chromaService)

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	// Create a larger test file
	testContent := bytes.Repeat([]byte("This is test content for a larger file. "), 1000)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "large.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write(testContent)
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

func TestUploadHandler_UploadNovel_EmptyFile(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	handler := NewUploadHandler(novelService, chromaService)

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	// Create empty test file
	testContent := ""
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "empty.txt")
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

	// Should succeed even with empty file
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUploadHandler_UploadNovel_SpecialCharactersInFilename(t *testing.T) {
	novelService := services.NewNovelService("test_novels")
	chromaService := services.NewChromaService("test_chroma_db")
	handler := NewUploadHandler(novelService, chromaService)

	defer func() {
		os.RemoveAll("test_novels")
		os.RemoveAll("test_chroma_db")
	}()

	// Create test file with special characters in filename
	testContent := "Test content with special filename"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test-file_123.txt")
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
