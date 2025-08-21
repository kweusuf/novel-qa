package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewChromaService(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	if service.dbPath != dbPath {
		t.Errorf("Expected dbPath %s, got %s", dbPath, service.dbPath)
	}

	// Check if directory was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Expected database directory to be created")
	}
}

func TestChromaService_getCollectionPath(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	expectedPath := filepath.Join(dbPath, "documents.json")
	actualPath := service.getCollectionPath()

	if actualPath != expectedPath {
		t.Errorf("Expected collection path %s, got %s", expectedPath, actualPath)
	}
}

func TestChromaService_Initialize_NewCollection(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Remove any existing collection to test new creation
	collectionPath := service.getCollectionPath()
	os.Remove(collectionPath)

	service.Initialize()

	// Check if collection file was created
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		t.Error("Expected collection file to be created")
	}

	// Check if collection file contains empty array
	data, err := os.ReadFile(collectionPath)
	if err != nil {
		t.Errorf("Failed to read collection file: %v", err)
	}

	expectedContent := "[]"
	if string(data) != expectedContent {
		t.Errorf("Expected empty array %s, got %s", expectedContent, string(data))
	}
}

func TestChromaService_Initialize_ExistingCollection(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Create a collection file first
	collectionPath := service.getCollectionPath()
	initialData := `[{"id":"test","text":"test content","embed":[0.1,0.2]}]`
	err := os.WriteFile(collectionPath, []byte(initialData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test collection file: %v", err)
	}

	service.Initialize()

	// Check if existing collection is preserved
	data, err := os.ReadFile(collectionPath)
	if err != nil {
		t.Errorf("Failed to read collection file: %v", err)
	}

	if string(data) != initialData {
		t.Errorf("Expected existing data to be preserved, got %s", string(data))
	}
}

func TestChromaService_AddDocuments(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Initialize empty collection
	service.Initialize()

	// Create test chunks
	chunks := []NovelChunk{
		{ID: "test1", Text: "This is test content 1"},
		{ID: "test2", Text: "This is test content 2"},
	}

	err := service.AddDocuments(chunks)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify documents were added
	collectionPath := service.getCollectionPath()
	data, err := os.ReadFile(collectionPath)
	if err != nil {
		t.Errorf("Failed to read collection file: %v", err)
	}

	var docs []ChromaDocument
	if err := json.Unmarshal(data, &docs); err != nil {
		t.Errorf("Failed to unmarshal documents: %v", err)
	}

	if len(docs) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(docs))
	}

	// Check document properties
	for i, doc := range docs {
		if doc.ID != chunks[i].ID {
			t.Errorf("Expected ID %s, got %s", chunks[i].ID, doc.ID)
		}
		if doc.Text != chunks[i].Text {
			t.Errorf("Expected text %s, got %s", chunks[i].Text, doc.Text)
		}
		if len(doc.Embed) != 384 {
			t.Errorf("Expected embedding length 384, got %d", len(doc.Embed))
		}
	}
}

func TestChromaService_AddDocuments_AppendToExisting(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Initialize with existing document
	service.Initialize()
	existingChunks := []NovelChunk{
		{ID: "existing", Text: "Existing content"},
	}
	service.AddDocuments(existingChunks)

	// Add new chunks
	newChunks := []NovelChunk{
		{ID: "new1", Text: "New content 1"},
		{ID: "new2", Text: "New content 2"},
	}
	err := service.AddDocuments(newChunks)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all documents exist
	data, err := os.ReadFile(service.getCollectionPath())
	if err != nil {
		t.Errorf("Failed to read collection file: %v", err)
	}

	var docs []ChromaDocument
	if err := json.Unmarshal(data, &docs); err != nil {
		t.Errorf("Failed to unmarshal documents: %v", err)
	}

	if len(docs) != 3 {
		t.Errorf("Expected 3 documents, got %d", len(docs))
	}
}

func TestChromaService_Query_WithMatches(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Initialize and add test documents
	service.Initialize()
	chunks := []NovelChunk{
		{ID: "doc1", Text: "The quick brown fox jumps over the lazy dog"},
		{ID: "doc2", Text: "A brown bear is a large mammal"},
		{ID: "doc3", Text: "The lazy dog sleeps all day"},
	}
	service.AddDocuments(chunks)

	// Query for "brown"
	result, err := service.Query("brown", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Check that result contains matching documents
	lines := strings.Split(result, "\n\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 results, got %d", len(lines))
	}

	// Verify both brown-related documents are included
	if !strings.Contains(result, "fox") || !strings.Contains(result, "bear") {
		t.Error("Expected result to contain both brown-related documents")
	}
}

func TestChromaService_Query_NoMatches(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Initialize and add test documents
	service.Initialize()
	chunks := []NovelChunk{
		{ID: "doc1", Text: "The quick brown fox jumps over the lazy dog"},
		{ID: "doc2", Text: "A brown bear is a large mammal"},
	}
	service.AddDocuments(chunks)

	// Query for non-existent term
	result, err := service.Query("purple", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected fallback result when no matches found")
	}

	// Should return first few documents as fallback
	lines := strings.Split(result, "\n\n")
	if len(lines) < 1 {
		t.Error("Expected at least one fallback document")
	}
}

func TestChromaService_Query_EmptyCollection(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Initialize empty collection
	service.Initialize()

	result, err := service.Query("test", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result for empty collection, got %s", result)
	}
}

func TestChromaService_Query_NonExistentCollection(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Don't initialize collection

	_, err := service.Query("test", 2)
	if err == nil {
		t.Error("Expected error for non-existent collection")
	}
}

func TestChromaService_generateDummyEmbedding(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	embedding := service.generateDummyEmbedding()

	if len(embedding) != 384 {
		t.Errorf("Expected embedding length 384, got %d", len(embedding))
	}

	// Check that all values are between 0 and 1
	for i, val := range embedding {
		if val < 0 || val > 1 {
			t.Errorf("Embedding value at index %d is out of range: %f", i, val)
		}
	}
}

func TestChromaService_Query_LimitResults(t *testing.T) {
	dbPath := "test_chroma_db"
	service := NewChromaService(dbPath)
	defer os.RemoveAll(dbPath)

	// Initialize and add many test documents
	service.Initialize()
	chunks := []NovelChunk{
		{ID: "doc1", Text: "Test content one"},
		{ID: "doc2", Text: "Test content two"},
		{ID: "doc3", Text: "Test content three"},
		{ID: "doc4", Text: "Test content four"},
		{ID: "doc5", Text: "Test content five"},
	}
	service.AddDocuments(chunks)

	// Query with limit of 3
	result, err := service.Query("test", 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	lines := strings.Split(result, "\n\n")
	if len(lines) != 3 {
		t.Errorf("Expected exactly 3 results, got %d", len(lines))
	}
}
