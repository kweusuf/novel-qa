package services

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ChromaService struct {
	dbPath string
}

type ChromaDocument struct {
	ID    string    `json:"id"`
	Text  string    `json:"text"`
	Embed []float64 `json:"embed"`
}

func NewChromaService(dbPath string) *ChromaService {
	os.MkdirAll(dbPath, 0755)
	return &ChromaService{dbPath: dbPath}
}

func (cs *ChromaService) getCollectionPath() string {
	return filepath.Join(cs.dbPath, "documents.json")
}

func (cs *ChromaService) AddDocuments(chunks []NovelChunk) error {
	docs := []ChromaDocument{}

	// Load existing documents
	if data, err := os.ReadFile(cs.getCollectionPath()); err == nil {
		json.Unmarshal(data, &docs)
	}

	// Add new chunks
	for _, chunk := range chunks {
		docs = append(docs, ChromaDocument{
			ID:    chunk.ID,
			Text:  chunk.Text,
			Embed: cs.generateDummyEmbedding(), // In real app, use actual embedding
		})
	}

	data, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	return os.WriteFile(cs.getCollectionPath(), data, 0644)
}

func (cs *ChromaService) Query(question string, nResults int) (string, error) {
	data, err := os.ReadFile(cs.getCollectionPath())
	if err != nil {
		return "", err
	}

	var docs []ChromaDocument
	if err := json.Unmarshal(data, &docs); err != nil {
		return "", err
	}

	// Simple keyword matching (in real app, use vector similarity)
	var results []string
	questionLower := strings.ToLower(question)

	for _, doc := range docs {
		if len(results) >= nResults {
			break
		}
		if strings.Contains(strings.ToLower(doc.Text), questionLower) {
			results = append(results, doc.Text)
		}
	}

	// If no matches found, return first few documents
	if len(results) == 0 && len(docs) > 0 {
		for i := 0; i < nResults && i < len(docs); i++ {
			results = append(results, docs[i].Text)
		}
	}

	return strings.Join(results, "\n\n"), nil
}

// Dummy embedding generation (replace with real sentence transformer)
func (cs *ChromaService) generateDummyEmbedding() []float64 {
	rand.Seed(time.Now().UnixNano())
	embed := make([]float64, 384) // MiniLM-L6-v2 dimension
	for i := range embed {
		embed[i] = rand.Float64()
	}
	return embed
}

func (cs *ChromaService) Initialize() {
	// Create initial collection if it doesn't exist
	if _, err := os.Stat(cs.getCollectionPath()); os.IsNotExist(err) {
		docs := []ChromaDocument{}
		data, _ := json.Marshal(docs)
		os.WriteFile(cs.getCollectionPath(), data, 0644)
		log.Println("📁 Created new ChromaDB collection")
	} else {
		log.Println("🔁 Using existing ChromaDB collection")
	}
}
