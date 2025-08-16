// handlers/qa.go
package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kweusuf/novel-qa-go/models"
	"github.com/kweusuf/novel-qa-go/services"

	"github.com/gin-gonic/gin"
)

type QAHandler struct {
	novelService  *services.NovelService
	chromaService *services.ChromaService
	ollamaService *services.OllamaService
}

func NewQAHandler(ns *services.NovelService, cs *services.ChromaService, os *services.OllamaService) *QAHandler {
	return &QAHandler{
		novelService:  ns,
		chromaService: cs,
		ollamaService: os,
	}
}

func (qh *QAHandler) ShowIndex(c *gin.Context) {
	models := []string{"phi3", "llama3", "mistral", "gemma"}
	c.HTML(http.StatusOK, "index.html", gin.H{"models": models})
}

func (qh *QAHandler) UploadNovel(c *gin.Context) {
	// Use MultipartForm to get all files associated with the 'files' field
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to parse multipart form: %v", err)
		return
	}

	files := form.File["files"] // Get the slice of *multipart.FileHeader
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "No files provided")
		return
	}

	var results []string // To store results for each file
	processedCount := 0

	for _, fileHeader := range files {
		// Validate file type (same as before)
		if len(fileHeader.Filename) < 4 || fileHeader.Filename[len(fileHeader.Filename)-4:] != ".txt" {
			results = append(results, fmt.Sprintf("Skipped '%s': Only .txt files allowed", fileHeader.Filename))
			continue
		}

		// Save file (same as before, but using fileHeader)
		dst := "novels/" + fileHeader.Filename
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			results = append(results, fmt.Sprintf("Failed to save '%s': %v", fileHeader.Filename, err))
			continue // Continue with next file
		}

		// Read content (same as before)
		content, err := qh.novelService.ReadNovel(dst)
		if err != nil {
			results = append(results, fmt.Sprintf("Failed to read '%s': %v", fileHeader.Filename, err))
			continue // Continue with next file
		}

		// Process and add to ChromaDB (same core logic)
		chunks := qh.novelService.ProcessNovel(fileHeader.Filename, content)
		err = qh.chromaService.AddDocuments(chunks)
		if err != nil {
			results = append(results, fmt.Sprintf("Failed to add '%s' to DB: %v", fileHeader.Filename, err))
			continue // Continue with next file
		}

		results = append(results, fmt.Sprintf("Successfully uploaded '%s' (%d chunks added)", fileHeader.Filename, len(chunks)))
		processedCount++
	}

	if processedCount == 0 {
		// If no files were successfully processed
		c.String(http.StatusInternalServerError, "No files were successfully uploaded. Details:\n%s", strings.Join(results, "\n"))
		return
	}

	// Return a summary of results
	c.String(http.StatusOK, "Upload Summary:\n%s", strings.Join(results, "\n"))
}

func (qh *QAHandler) AskQuestion(c *gin.Context) {
	var req models.QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Validate model
	validModels := map[string]bool{
		"phi3": true, "llama3": true, "mistral": true, "gemma": true,
	}
	if !validModels[req.Model] {
		req.Model = "phi3" // Default to phi3 if invalid
	}

	// Get context from ChromaDB
	context, err := qh.chromaService.Query(req.Question, 2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve context: " + err.Error()})
		return
	}

	// Use custom endpoint if provided, otherwise use default service
	var answer string
	if req.OllamaEndpoint != "" {
		// Create a temporary Ollama service with custom endpoint
		customOllamaService := services.NewOllamaService(req.OllamaEndpoint)
		answer, err = customOllamaService.Ask(req.Question, req.Model, context)
	} else {
		// Use the default Ollama service
		answer, err = qh.ollamaService.Ask(req.Question, req.Model, context)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get answer from model: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}

func (qh *QAHandler) GetModels(c *gin.Context) {
	// Get Ollama endpoint from query parameter or use default
	ollamaEndpoint := c.Query("endpoint")
	if ollamaEndpoint == "" {
		// Use the default Ollama service endpoint from the handler's service
		// This will use the same endpoint that was configured at startup
		models, err := qh.ollamaService.GetModels()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get models: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"models": models})
		return
	}

	// Create a temporary Ollama service to get models from custom endpoint
	ollamaService := services.NewOllamaService(ollamaEndpoint)

	// Get available models
	models, err := ollamaService.GetModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get models: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}
