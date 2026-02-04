package resolver

import (
	"testing"

	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/llm"
)

func TestResolveModel(t *testing.T) {
	tests := []struct {
		name     string
		tier     string
		intent   llm.IntentLabel
		expected string
	}{
		// Free Tier Cases
		{"Free - Coding", "free", llm.LabelCoding, ModelDeepSeekCoder},
		{"Free - General", "free", llm.LabelGeneral, ModelLlama3_8B},
		{"Free - Creative", "free", llm.LabelCreative, ModelGemma2B},

		// Premium Tier Cases
		{"Premium - Coding", "premium", llm.LabelCoding, ModelLlama3_70B},
		{"Premium - Creative", "premium", llm.LabelCreative, ModelLlama3_8B},

		// Edge Cases
		{"Empty Tier (Default Free)", "", llm.LabelCoding, ModelDeepSeekCoder},
		{"Unknown Intent (Default)", "free", "unknown_label", ModelLlama3_8B},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveModel(tt.tier, tt.intent)
			if got != tt.expected {
				t.Errorf("ResolveModel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Command to run tests -> go test ./internal/resolver -v
