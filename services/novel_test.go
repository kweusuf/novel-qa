package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewNovelService(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)

	if service.novelsDir != dir {
		t.Errorf("Expected novelsDir %s, got %s", dir, service.novelsDir)
	}

	// Check if directory was created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("Expected novels directory to be created")
	}

	// Clean up
	os.RemoveAll(dir)
}

func TestNovelService_ReadNovel(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	// Create a test file
	testContent := "This is test novel content."
	filePath := filepath.Join(dir, "test.txt")
	err := os.WriteFile(filePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	content, err := service.ReadNovel(filePath)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if content != testContent {
		t.Errorf("Expected content %s, got %s", testContent, content)
	}
}

func TestNovelService_ReadNovel_FileNotFound(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	_, err := service.ReadNovel("nonexistent.txt")

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestNovelService_ProcessNovel(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	testContent := "This is a test novel with some content that should be processed into chunks for testing purposes."
	chunks := service.ProcessNovel("test.txt", testContent)

	if len(chunks) == 0 {
		t.Error("Expected at least one chunk, got zero")
	}

	// Check that chunks have proper IDs and text
	for i, chunk := range chunks {
		if chunk.ID == "" {
			t.Errorf("Chunk %d has empty ID", i)
		}
		if chunk.Text == "" {
			t.Errorf("Chunk %d has empty text", i)
		}
		if len(chunk.Text) > 400*2 { // Allow some tolerance for word boundaries
			t.Errorf("Chunk %d text is too long: %d characters", i, len(chunk.Text))
		}
	}
}

func TestNovelService_ProcessNovel_EmptyContent(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	chunks := service.ProcessNovel("empty.txt", "")

	if len(chunks) != 0 {
		t.Errorf("Expected 0 chunks for empty content, got %d", len(chunks))
	}
}

func TestNovelService_ProcessNovel_ShortContent(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	testContent := "Short content"
	chunks := service.ProcessNovel("short.txt", testContent)

	if len(chunks) != 1 {
		t.Errorf("Expected 1 chunk for short content, got %d", len(chunks))
	}

	if chunks[0].Text != testContent {
		t.Errorf("Expected chunk text %s, got %s", testContent, chunks[0].Text)
	}
}

func TestNovelService_SaveNovel(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	filename := "test_save.txt"
	content := []byte("Test content to save")

	err := service.SaveNovel(filename, content)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(dir, filename)
	savedContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read saved file: %v", err)
	}

	if string(savedContent) != string(content) {
		t.Errorf("Expected saved content %s, got %s", string(content), string(savedContent))
	}
}

func TestNovelService_LoadNovels(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	// Create test files
	testFiles := map[string]string{
		"novel1.txt": "This is the first test novel content.",
		"novel2.txt": "This is the second test novel content with more text to create chunks.",
		"readme.md":  "This should be ignored",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(dir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	chunks, err := service.LoadNovels()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should only load .txt files
	txtFileCount := 0
	for filename := range testFiles {
		if filepath.Ext(filename) == ".txt" {
			txtFileCount++
		}
	}

	if len(chunks) == 0 {
		t.Error("Expected at least one chunk from .txt files")
	}

	// Check that all chunks have proper structure
	for i, chunk := range chunks {
		if chunk.ID == "" {
			t.Errorf("Chunk %d has empty ID", i)
		}
		if chunk.Text == "" {
			t.Errorf("Chunk %d has empty text", i)
		}
	}
}
