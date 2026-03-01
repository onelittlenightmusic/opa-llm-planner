package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/hiroyukiosaki/opa-llm-planner/internal/types"
)

// LLMClient defines the interface for LLM interactions.
type LLMClient interface {
	// EnrichAction fills in description and parameters for an action using LLM.
	EnrichAction(ctx context.Context, action types.Action, goal, current map[string]interface{}) (types.Action, error)
	// GenerateRegoRules generates Rego rules for missing actions.
	GenerateRegoRules(ctx context.Context, missing []string, goal, current map[string]interface{}) (string, error)
}

// NewClient creates an LLMClient for the given provider.
// provider is "anthropic" or "openai". Falls back to LLM_PROVIDER env var.
func NewClient(provider string) (LLMClient, error) {
	if provider == "" {
		provider = os.Getenv("LLM_PROVIDER")
	}
	switch provider {
	case "anthropic":
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
		}
		return NewAnthropicClient(apiKey), nil
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is not set")
		}
		return NewOpenAIClient(apiKey), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider %q; use 'anthropic' or 'openai'", provider)
	}
}
