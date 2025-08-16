// main.go
package main

import (
	"log"
	"os"

	"github.com/kweusuf/novel-qa-go/handlers"
	"github.com/kweusuf/novel-qa-go/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get Ollama host from environment variable, fallback to localhost
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		if os.Getenv("IS_DOCKER") == "true" {
			ollamaHost = "http://host.docker.internal:11434"
		} else {
			ollamaHost = "http://localhost:11434"
		}
	}

	// Initialize services
	novelService := services.NewNovelService("novels")
	chromaService := services.NewChromaService("chroma_db")
	chromaService.Initialize() // Initialize the ChromaDB
	ollamaService := services.NewOllamaService(ollamaHost)

	// Initialize handler
	qaHandler := handlers.NewQAHandler(novelService, chromaService, ollamaService)

	// Set up Gin
	r := gin.Default()

	// Register static files handler
	r.Static("/static", "./static")

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Public routes (no authentication)
	r.GET("/", qaHandler.ShowIndex)
	r.POST("/upload", qaHandler.UploadNovel)
	r.POST("/ask", qaHandler.AskQuestion)
	r.GET("/models", qaHandler.GetModels)

	log.Printf("🚀 Starting server at http://localhost:8080")
	log.Printf("🔗 Using Ollama at: %s", ollamaHost)
	r.Run(":8080")
}
