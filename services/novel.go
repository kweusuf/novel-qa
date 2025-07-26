// services/novel.go
package services

import (
	"bufio"
	"fmt"
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
		if filepath.Ext(file.Name()) == ".txt" {
			filePath := filepath.Join(ns.novelsDir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			words := strings.Fields(string(content))
			for i := 0; i < len(words); i += 400 {
				end := i + 400
				if end > len(words) {
					end = len(words)
				}
				chunk := strings.Join(words[i:end], " ")
				chunks = append(chunks, NovelChunk{
					ID:   file.Name() + "-" + string(rune(i/400)),
					Text: chunk,
				})
			}
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
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
