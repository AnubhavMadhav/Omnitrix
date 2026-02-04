package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GroqClient struct {
	apiKey string
	client *http.Client
}

func NewGroqClient(apiKey string) *GroqClient {
	return &GroqClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Groq Request Structures (OpenAI Compatible)
type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (g *GroqClient) Generate(ctx context.Context, modelID, prompt string) (string, error) {

	reqBody := groqRequest{
		Model: modelID,
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("groq connection error: %w", err)
	}
	defer resp.Body.Close()

	// bodyBytes, _ := io.ReadAll(resp.Body)
	// fmt.Printf("API Error Body: %s\n", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq api error: status %d", resp.StatusCode)
	}

	// fmt.Println("Status Code : ", resp.StatusCode)

	var result groqResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("groq json decode error: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("groq returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}
