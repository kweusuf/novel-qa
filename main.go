// main.go
package main

import (
	"log"

	"github.com/kweusuf/novel-qa-go/handlers"
	"github.com/kweusuf/novel-qa-go/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize services
	novelService := services.NewNovelService("novels")
	chromaService := services.NewChromaService("chroma_db")
	chromaService.Initialize() // Initialize the ChromaDB
	ollamaService := services.NewOllamaService("http://localhost:11434")

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

	log.Println("🚀 Starting server at http://localhost:8080")
	r.Run(":8080")
}
