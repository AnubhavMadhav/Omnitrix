package llm

import "context"

// Generator defines the behavior for ANY AI provider (Ollama, Groq, OpenAI etc.).
type Generator interface {
	Generate(ctx context.Context, modelID, prompt string) (string, error)
}
