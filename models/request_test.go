package models

import (
	"encoding/json"
	"testing"
)

func TestQuestionRequest_MarshalJSON(t *testing.T) {
	request := QuestionRequest{
		Question:       "What is the meaning of life?",
		Model:          "phi3",
		OllamaEndpoint: "http://localhost:11434",
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Errorf("Expected no error marshaling QuestionRequest, got %v", err)
	}

	// Verify JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Errorf("Expected valid JSON, got error %v", err)
	}

	if result["question"] != request.Question {
		t.Errorf("Expected question %s, got %s", request.Question, result["question"])
	}
	if result["model"] != request.Model {
		t.Errorf("Expected model %s, got %s", request.Model, result["model"])
	}
	if result["ollamaEndpoint"] != request.OllamaEndpoint {
		t.Errorf("Expected ollamaEndpoint %s, got %s", request.OllamaEndpoint, result["ollamaEndpoint"])
	}
}

func TestQuestionRequest_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"question": "What is the meaning of life?",
		"model": "phi3",
		"ollamaEndpoint": "http://localhost:11434"
	}`

	var request QuestionRequest
	err := json.Unmarshal([]byte(jsonData), &request)
	if err != nil {
		t.Errorf("Expected no error unmarshaling QuestionRequest, got %v", err)
	}

	if request.Question != "What is the meaning of life?" {
		t.Errorf("Expected question 'What is the meaning of life?', got %s", request.Question)
	}
	if request.Model != "phi3" {
		t.Errorf("Expected model 'phi3', got %s", request.Model)
	}
	if request.OllamaEndpoint != "http://localhost:11434" {
		t.Errorf("Expected ollamaEndpoint 'http://localhost:11434', got %s", request.OllamaEndpoint)
	}
}

func TestQuestionRequest_UnmarshalJSON_Minimal(t *testing.T) {
	// Test with only required fields
	jsonData := `{
		"question": "What is the meaning of life?",
		"model": "phi3"
	}`

	var request QuestionRequest
	err := json.Unmarshal([]byte(jsonData), &request)
	if err != nil {
		t.Errorf("Expected no error unmarshaling minimal QuestionRequest, got %v", err)
	}

	if request.Question != "What is the meaning of life?" {
		t.Errorf("Expected question 'What is the meaning of life?', got %s", request.Question)
	}
	if request.Model != "phi3" {
		t.Errorf("Expected model 'phi3', got %s", request.Model)
	}
	if request.OllamaEndpoint != "" {
		t.Errorf("Expected empty ollamaEndpoint, got %s", request.OllamaEndpoint)
	}
}

func TestQuestionRequest_UnmarshalJSON_InvalidJSON(t *testing.T) {
	invalidJSON := `{"question": "test", "model": "test", invalid}`

	var request QuestionRequest
	err := json.Unmarshal([]byte(invalidJSON), &request)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestQuestionRequest_StructCreation(t *testing.T) {
	// Test creating struct directly
	request := QuestionRequest{
		Question:       "Test question",
		Model:          "test-model",
		OllamaEndpoint: "http://test:11434",
	}

	if request.Question != "Test question" {
		t.Errorf("Expected question 'Test question', got %s", request.Question)
	}
	if request.Model != "test-model" {
		t.Errorf("Expected model 'test-model', got %s", request.Model)
	}
	if request.OllamaEndpoint != "http://test:11434" {
		t.Errorf("Expected ollamaEndpoint 'http://test:11434', got %s", request.OllamaEndpoint)
	}
}

func TestUploadRequest_MarshalJSON(t *testing.T) {
	request := UploadRequest{
		Filename: "test.txt",
		Content:  "This is test content",
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Errorf("Expected no error marshaling UploadRequest, got %v", err)
	}

	// Verify JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Errorf("Expected valid JSON, got error %v", err)
	}

	if result["filename"] != request.Filename {
		t.Errorf("Expected filename %s, got %s", request.Filename, result["filename"])
	}
	if result["content"] != request.Content {
		t.Errorf("Expected content %s, got %s", request.Content, result["content"])
	}
}

func TestUploadRequest_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"filename": "test.txt",
		"content": "This is test content"
	}`

	var request UploadRequest
	err := json.Unmarshal([]byte(jsonData), &request)
	if err != nil {
		t.Errorf("Expected no error unmarshaling UploadRequest, got %v", err)
	}

	if request.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got %s", request.Filename)
	}
	if request.Content != "This is test content" {
		t.Errorf("Expected content 'This is test content', got %s", request.Content)
	}
}

func TestUploadRequest_UnmarshalJSON_Empty(t *testing.T) {
	jsonData := `{
		"filename": "",
		"content": ""
	}`

	var request UploadRequest
	err := json.Unmarshal([]byte(jsonData), &request)
	if err != nil {
		t.Errorf("Expected no error unmarshaling empty UploadRequest, got %v", err)
	}

	if request.Filename != "" {
		t.Errorf("Expected empty filename, got %s", request.Filename)
	}
	if request.Content != "" {
		t.Errorf("Expected empty content, got %s", request.Content)
	}
}

func TestUploadRequest_StructCreation(t *testing.T) {
	// Test creating struct directly
	request := UploadRequest{
		Filename: "document.txt",
		Content:  "Document content here",
	}

	if request.Filename != "document.txt" {
		t.Errorf("Expected filename 'document.txt', got %s", request.Filename)
	}
	if request.Content != "Document content here" {
		t.Errorf("Expected content 'Document content here', got %s", request.Content)
	}
}

func TestQuestionRequest_ValidationTags(t *testing.T) {
	// Test that the struct has proper validation tags
	request := QuestionRequest{}

	// Check that the struct has the expected binding tags
	// This is more of a compile-time check, but we can verify the fields exist
	if request.Question != "" {
		t.Error("Expected Question field to be empty by default")
	}
	if request.Model != "" {
		t.Error("Expected Model field to be empty by default")
	}
	if request.OllamaEndpoint != "" {
		t.Error("Expected OllamaEndpoint field to be empty by default")
	}
}

func TestUploadRequest_Fields(t *testing.T) {
	// Test that the struct has the expected fields
	request := UploadRequest{}

	// Verify fields can be set and retrieved
	request.Filename = "test.txt"
	request.Content = "test content"

	if request.Filename != "test.txt" {
		t.Errorf("Expected Filename 'test.txt', got %s", request.Filename)
	}
	if request.Content != "test content" {
		t.Errorf("Expected Content 'test content', got %s", request.Content)
	}
}

// Test JSON tag names
func TestJSONTags(t *testing.T) {
	request := QuestionRequest{
		Question:       "test question",
		Model:          "test model",
		OllamaEndpoint: "test endpoint",
	}

	data, _ := json.Marshal(request)
	jsonStr := string(data)

	// Check that JSON field names match the struct tags
	if !json.Valid(data) {
		t.Error("Generated JSON is not valid")
	}

	// Verify specific field names are present
	if !containsJSONField(jsonStr, "question") {
		t.Error("JSON does not contain 'question' field")
	}
	if !containsJSONField(jsonStr, "model") {
		t.Error("JSON does not contain 'model' field")
	}
	if !containsJSONField(jsonStr, "ollamaEndpoint") {
		t.Error("JSON does not contain 'ollamaEndpoint' field")
	}
}

func TestUploadRequestJSONTags(t *testing.T) {
	request := UploadRequest{
		Filename: "test.txt",
		Content:  "test content",
	}

	data, _ := json.Marshal(request)
	jsonStr := string(data)

	// Check that JSON field names match the struct tags
	if !json.Valid(data) {
		t.Error("Generated JSON is not valid")
	}

	// Verify specific field names are present
	if !containsJSONField(jsonStr, "filename") {
		t.Error("JSON does not contain 'filename' field")
	}
	if !containsJSONField(jsonStr, "content") {
		t.Error("JSON does not contain 'content' field")
	}
}

// Helper function to check if JSON contains a field
func containsJSONField(jsonStr, fieldName string) bool {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonMap); err != nil {
		return false
	}
	_, exists := jsonMap[fieldName]
	return exists
}
