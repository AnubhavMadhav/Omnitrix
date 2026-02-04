package llm

import (
	"context"
)

type IntentLabel string

const (
	LabelCoding   IntentLabel = "coding"
	LabelCreative IntentLabel = "creative"
	LabelMath     IntentLabel = "math_logic"
	LabelSQL      IntentLabel = "sql_db"
	LabelGeneral  IntentLabel = "general"
)

type ClassifierResponse struct {
	Label      IntentLabel `json:"label"`
	Confidence float64     `json:"confidence"`
}

// This allows us to swap Ollama for OpenAI or Mock tests later easily.
type Client interface {
	Classify(ctx context.Context, prompt string) (ClassifierResponse, error)
}
