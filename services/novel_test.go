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

	// Should only load .txt files (EPUB files would be tested separately with real EPUB files)
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

func TestNovelService_ReadNovel_EPUB(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	// Create a mock EPUB file (this will fail to process but tests the detection)
	testContent := "This is not a real EPUB file but tests extension detection"
	filePath := filepath.Join(dir, "test.epub")
	err := os.WriteFile(filePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test EPUB file: %v", err)
	}

	// This should attempt to process as EPUB and fail gracefully
	content, err := service.ReadNovel(filePath)

	// We expect an error because our mock content isn't a valid EPUB
	if err == nil {
		t.Error("Expected error when reading invalid EPUB file, got nil")
	}

	if content != "" {
		t.Errorf("Expected empty content for invalid EPUB, got: %s", content)
	}
}

func TestNovelService_ReadNovel_TXT(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	testContent := "This is test content for TXT file."
	filePath := filepath.Join(dir, "test.txt")
	err := os.WriteFile(filePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test TXT file: %v", err)
	}

	content, err := service.ReadNovel(filePath)

	if err != nil {
		t.Errorf("Expected no error reading TXT file, got %v", err)
	}

	if content != testContent {
		t.Errorf("Expected content %s, got %s", testContent, content)
	}
}

func TestNovelService_ExtractTextFromHTML(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple HTML with tags",
			html:     "<p>This is <strong>bold</strong> text.</p>",
			expected: "This is bold text.",
		},
		{
			name:     "HTML with nested tags",
			html:     "<div><p>Nested <em>emphasized</em> text</p></div>",
			expected: "Nested emphasized text",
		},
		{
			name:     "HTML with extra whitespace",
			html:     "<p>  Text   with   spaces  </p>",
			expected: "Text with spaces",
		},
		{
			name:     "HTML with line breaks",
			html:     "<p>Line one</p>\n<p>Line two</p>",
			expected: "Line one Line two",
		},
		{
			name:     "Empty HTML",
			html:     "<p></p>",
			expected: "",
		},
		{
			name:     "HTML with only tags",
			html:     "<div><span></span></div>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractTextFromHTML(tt.html)
			if result != tt.expected {
				t.Errorf("extractTextFromHTML(%q) = %q, want %q", tt.html, result, tt.expected)
			}
		})
	}
}

func TestNovelService_ProcessNovel_WithHTMLTags(t *testing.T) {
	dir := "test_novels"
	service := NewNovelService(dir)
	defer os.RemoveAll(dir)

	// Test content with HTML tags (simulating EPUB content)
	htmlContent := "<html><body><h1>Title</h1><p>This is paragraph one.</p><p>This is paragraph two.</p></body></html>"

	// First test the HTML extraction function directly
	cleanContent := service.extractTextFromHTML(htmlContent)
	expectedClean := "Title This is paragraph one. This is paragraph two."

	if cleanContent != expectedClean {
		t.Errorf("extractTextFromHTML() = %q, want %q", cleanContent, expectedClean)
	}

	// Then test that ProcessNovel works with clean content
	chunks := service.ProcessNovel("test.epub", cleanContent)

	if len(chunks) == 0 {
		t.Error("Expected at least one chunk, got zero")
	}

	// Check that the first chunk contains the clean text
	if chunks[0].Text != cleanContent {
		t.Errorf("Expected chunk text %q, got %q", cleanContent, chunks[0].Text)
	}
}
