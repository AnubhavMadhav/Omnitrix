package provider

import (
	"fmt"

	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/llm"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/resolver"
)

// Factory holds instances of our providers
type Factory struct {
	ollama *llm.OllamaClient
	groq   *llm.GroqClient
}

func NewFactory(ollama *llm.OllamaClient, groq *llm.GroqClient) *Factory {
	return &Factory{
		ollama: ollama,
		groq:   groq,
	}
}

// GetGenerator returns the correct client for a given model ID
func (f *Factory) GetGenerator(modelID string) (llm.Generator, error) {
	// Logic: If it's a known Groq model, return Groq client.
	// Otherwise, default to Ollama.

	switch modelID {
	case resolver.ModelLlama3_70B, resolver.ModelLlama3_8B:
		// These are our Premium/Cloud models
		if f.groq == nil {
			return nil, fmt.Errorf("premium model requested but groq client is not configured")
		}
		return f.groq, nil

	case resolver.ModelPhi3Mini, resolver.ModelGemma2B, resolver.ModelDeepSeekCoder, resolver.ModelLlama3_8BLocal:
		// These are our Local/Free models
		return f.ollama, nil

	default:
		// Safety fallback
		return f.ollama, nil
	}
}
