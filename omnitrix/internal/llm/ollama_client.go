package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OllamaClient struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaClient creates a new client.
func NewOllamaClient(baseURL, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "phi3" // Default model
	}

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 120 * time.Second, // I've set this to 120s because generating code was taking much more time in localhost than writing poems.
		},
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"` // We want a single JSON response, not a stream for now
	// Stream can be used in future if integrating any UI
}

type ollamaResponse struct {
	Response string `json:"response"`
}

// Classify takes the user prompt and returns the intent label
func (o *OllamaClient) Classify(ctx context.Context, userPrompt string) (ClassifierResponse, error) {
	// System Prompt to respond with only "label"
	fullPrompt := fmt.Sprintf(`
		You are an Intent Classifier. 
		Analyze the user's input and categorize it into exactly one of these labels:
		- coding (writing code, debugging, explanation of code)
		- creative (stories, poems, jokes, songs, movies)
		- math_logic (calculations, logic puzzles, reasoning)
		- sql_db (SQL queries, database design, schema)
		- general (history, facts, general knowledge, greetings that missed the reflex layer)

		Input: "%s"

		Instructions:
		1. Output ONLY the label name.
		2. Do not add punctuation or explanations.
		3. If unsure, default to "general".
	`, userPrompt)

	reqBody := ollamaRequest{
		Model:  o.model,
		Prompt: fullPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return ClassifierResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return ClassifierResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return ClassifierResponse{}, fmt.Errorf("ollama connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ClassifierResponse{}, fmt.Errorf("ollama returned status: %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ClassifierResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	cleanedLabel := strings.TrimSpace(strings.ToLower(result.Response))

	// TODO: Validate if cleanedLabel is one of our supported IntentLabels or not

	return ClassifierResponse{
		Label:      IntentLabel(cleanedLabel),
		Confidence: 1.0, // Mocked it for now
	}, nil
}

func (o *OllamaClient) Generate(ctx context.Context, modelID, prompt string) (string, error) {

	reqBody := ollamaRequest{
		Model:  modelID,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama status: %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	return result.Response, nil
}
