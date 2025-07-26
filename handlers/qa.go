// handlers/qa.go
package handlers

import (
	"fmt"
	"net/http"

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
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to get file: %v", err)
		return
	}

	if len(file.Filename) < 4 || file.Filename[len(file.Filename)-4:] != ".txt" {
		c.String(http.StatusBadRequest, "Only .txt files allowed")
		return
	}

	// Save file
	dst := "novels/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusInternalServerError, "Failed to save file: %v", err)
		return
	}

	// Read content
	content, err := qh.novelService.ReadNovel(dst)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %v", err)
		return
	}

	// Process and add to ChromaDB
	chunks := qh.novelService.ProcessNovel(file.Filename, content)
	err = qh.chromaService.AddDocuments(chunks)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to add to database: %v", err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Successfully uploaded '%s' (%d chunks added)", file.Filename, len(chunks)))
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

	// Ask Ollama
	answer, err := qh.ollamaService.Ask(req.Question, req.Model, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get answer from model: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
