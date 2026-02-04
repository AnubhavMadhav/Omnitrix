package resolver

import (
	"strings"

	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/llm"
)

// Model Identifiers
const (
	// Ollama Models
	// For Setup/trying Omnitrix in your system -> Make sure you have pulled these in local (if running on local)
	ModelPhi3Mini       = "phi3:mini"
	ModelDeepSeekCoder  = "deepseek-coder"
	ModelLlama3_8BLocal = "llama3"
	ModelGemma2B        = "gemma:2b"

	// Groq Models
	ModelLlama3_70B = "llama-3.3-70b-versatile"
	ModelLlama3_8B  = "llama-3.1-8b-instant"
)

// User Tiers
const (
	TierFree    = "free"
	TierPremium = "premium"
)

// ResolveModel determines which LLM to use based on Intent and Tier.
func ResolveModel(tier string, intent llm.IntentLabel) string {
	// Normalize tier (default to free if empty or invalid)
	tier = strings.ToLower(tier)
	if tier != TierPremium {
		tier = TierFree
	}

	switch intent {
	case llm.LabelCoding:
		if tier == TierPremium {
			return ModelLlama3_70B
		}
		return ModelDeepSeekCoder

	case llm.LabelCreative:
		if tier == TierPremium {
			return ModelLlama3_8B
		}
		return ModelGemma2B

	case llm.LabelMath, llm.LabelSQL:
		if tier == TierPremium {
			return ModelLlama3_70B
		}
		return ModelPhi3Mini

	case llm.LabelGeneral:
		if tier == TierPremium {
			return ModelLlama3_70B
		}
		return ModelLlama3_8BLocal

	default:
		if tier == TierPremium {
			return ModelLlama3_70B // Safe default for Groq (Premium)
		}
		return ModelLlama3_8BLocal // Safe default for Local (Last Option)
	}
}
