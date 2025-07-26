package handlers

import (
	"fmt"
	"net/http"

	"github.com/kweusuf/novel-qa-go/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	novelService  *services.NovelService
	chromaService *services.ChromaService
}

func NewUploadHandler(ns *services.NovelService, cs *services.ChromaService) *UploadHandler {
	return &UploadHandler{
		novelService:  ns,
		chromaService: cs,
	}
}

func (uh *UploadHandler) UploadNovel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to get file: %v", err)
		return
	}

	if len(file.Filename) < 4 || file.Filename[len(file.Filename)-4:] != ".txt" {
		c.String(http.StatusBadRequest, "Only .txt files allowed")
		return
	}

	// Save file temporarily
	dst := "novels/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusInternalServerError, "Failed to save file: %v", err)
		return
	}

	// Read content
	content, err := uh.novelService.ReadNovel(dst)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %v", err)
		return
	}

	// Process and add to ChromaDB
	chunks := uh.novelService.ProcessNovel(file.Filename, content)
	err = uh.chromaService.AddDocuments(chunks)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to add to database: %v", err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Successfully uploaded '%s' (%d chunks added)", file.Filename, len(chunks)))
}
