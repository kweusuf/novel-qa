// services/novel.go
package services

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type NovelChunk struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type NovelService struct {
	novelsDir string
}

func NewNovelService(dir string) *NovelService {
	os.MkdirAll(dir, 0755)
	return &NovelService{novelsDir: dir}
}

func (ns *NovelService) LoadNovels() ([]NovelChunk, error) {
	var chunks []NovelChunk

	files, err := os.ReadDir(ns.novelsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := filepath.Join(ns.novelsDir, file.Name())

		var content string
		var err error

		if filepath.Ext(file.Name()) == ".txt" {
			contentBytes, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			content = string(contentBytes)
		} else if filepath.Ext(file.Name()) == ".epub" {
			content, err = ns.readEPUB(filePath)
			if err != nil {
				continue
			}
		} else {
			continue
		}

		words := strings.Fields(content)
		for i := 0; i < len(words); i += 400 {
			end := i + 400
			if end > len(words) {
				end = len(words)
			}
			chunk := strings.Join(words[i:end], " ")
			chunks = append(chunks, NovelChunk{
				ID:   file.Name() + "-" + fmt.Sprintf("%d", i/400),
				Text: chunk,
			})
		}
	}

	return chunks, nil
}

func (ns *NovelService) SaveNovel(filename string, content []byte) error {
	filePath := filepath.Join(ns.novelsDir, filename)
	return os.WriteFile(filePath, content, 0644)
}

func (ns *NovelService) ProcessNovel(filename string, content string) []NovelChunk {
	var chunks []NovelChunk
	scanner := bufio.NewScanner(strings.NewReader(content))
	var text strings.Builder

	for scanner.Scan() {
		text.WriteString(scanner.Text() + " ")
	}

	words := strings.Fields(text.String())
	for i := 0; i < len(words); i += 400 {
		end := i + 400
		if end > len(words) {
			end = len(words)
		}
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, NovelChunk{
			ID:   filename + "-" + fmt.Sprintf("%d", i/400),
			Text: chunk,
		})
	}

	return chunks
}

// Add the missing ReadNovel method
func (ns *NovelService) ReadNovel(filepath string) (string, error) {
	if strings.HasSuffix(filepath, ".epub") {
		return ns.readEPUB(filepath)
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// readEPUB reads and extracts text content from an EPUB file
func (ns *NovelService) readEPUB(filepath string) (string, error) {
	// Open the EPUB file as a ZIP archive
	reader, err := zip.OpenReader(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open EPUB file: %v", err)
	}
	defer reader.Close()

	var content strings.Builder

	// Iterate through all files in the EPUB
	for _, file := range reader.File {
		// Look for HTML/XHTML content files
		if strings.HasSuffix(file.Name, ".html") || strings.HasSuffix(file.Name, ".xhtml") {
			// Open the file from the ZIP archive
			rc, err := file.Open()
			if err != nil {
				continue
			}

			// Read the content
			data, err := io.ReadAll(rc)
			if err != nil {
				rc.Close()
				continue
			}
			rc.Close()

			// Simple text extraction (remove HTML tags)
			text := ns.extractTextFromHTML(string(data))
			content.WriteString(text)
			content.WriteString(" ")
		}
	}

	return content.String(), nil
}

// extractTextFromHTML performs basic HTML tag removal to extract text
func (ns *NovelService) extractTextFromHTML(html string) string {
	// Simple HTML tag removal (basic implementation)
	var result strings.Builder
	inTag := false

	for _, char := range html {
		if char == '<' {
			inTag = true
		} else if char == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(char)
		}
	}

	// Clean up extra whitespace
	text := strings.TrimSpace(result.String())
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Remove multiple spaces
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}
