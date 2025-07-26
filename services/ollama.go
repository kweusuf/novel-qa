// services/ollama.go
package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OllamaService struct {
	baseURL string
	client  *http.Client
}

type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // Explicitly set this
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaStreamResponse struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

func NewOllamaService(baseURL string) *OllamaService {
	return &OllamaService{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 300 * time.Second, // 5 minute timeout
		},
	}
}

func (os *OllamaService) Ask(question, model, context string) (string, error) {
	prompt := fmt.Sprintf(`
You are a helpful assistant answering questions based on a novel.
Use only the following context to answer. If unsure, say 'I don't know'.

Context:
%s

Question: %s
Answer:
`, context, question)

	reqBody := OllamaRequest{
		Model:  model,
		Stream: false, // Explicitly set to false
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := os.client.Post(os.baseURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Handle streaming response by reading all lines and taking the last complete one
	// Even though Stream=false, Ollama might still stream in some cases
	scanner := bufio.NewScanner(resp.Body)
	var fullContent strings.Builder
	var lastValidResponse OllamaStreamResponse

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var streamResp OllamaStreamResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			// If we can't parse a line, log it but continue
			fmt.Printf("Warning: Could not parse line as JSON: %s\n", line)
			continue
		}

		// Accumulate content if message exists
		if streamResp.Message.Content != "" {
			fullContent.WriteString(streamResp.Message.Content)
		}

		// Keep track of the last valid response (which should have Done=true)
		lastValidResponse = streamResp
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response stream: %w", err)
	}

	// Prefer the content from the last response if it's marked as done
	// Otherwise, use the accumulated content
	if lastValidResponse.Done && lastValidResponse.Message.Content != "" {
		return lastValidResponse.Message.Content, nil
	} else if fullContent.Len() > 0 {
		return fullContent.String(), nil
	} else if lastValidResponse.Message.Content != "" {
		return lastValidResponse.Message.Content, nil
	}

	return "", fmt.Errorf("no valid response content received")
}
