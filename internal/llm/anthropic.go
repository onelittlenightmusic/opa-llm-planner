package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/hiroyukiosaki/opa-llm-planner/internal/types"
)

const claudeModel = anthropic.Model("claude-sonnet-4-6")

// AnthropicClient implements LLMClient using Anthropic Claude.
type AnthropicClient struct {
	client *anthropic.Client
}

// NewAnthropicClient creates an AnthropicClient with the given API key.
func NewAnthropicClient(apiKey string) *AnthropicClient {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicClient{client: client}
}

// EnrichAction asks Claude to fill in description and parameters for the action.
func (c *AnthropicClient) EnrichAction(ctx context.Context, action types.Action, goal, current map[string]interface{}) (types.Action, error) {
	goalJSON, _ := json.MarshalIndent(goal, "", "  ")
	currentJSON, _ := json.MarshalIndent(current, "", "  ")

	prompt := fmt.Sprintf(`You are a planning assistant. Given an action type and context, provide a description and parameters for the action.

Action type: %s

Goal:
%s

Current state:
%s

Respond with a JSON object containing "description" (string) and "parameters" (object) fields only. Do not include any other text.`,
		action.Type, string(goalJSON), string(currentJSON))

	msg, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(claudeModel),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		return action, fmt.Errorf("anthropic API error: %w", err)
	}

	if len(msg.Content) == 0 {
		return action, fmt.Errorf("empty response from Claude")
	}

	responseText := msg.Content[0].Text

	var enriched struct {
		Description string                 `json:"description"`
		Parameters  map[string]interface{} `json:"parameters"`
	}
	if err := json.Unmarshal([]byte(responseText), &enriched); err != nil {
		action.Description = responseText
		return action, nil
	}

	action.Description = enriched.Description
	action.Parameters = enriched.Parameters
	return action, nil
}

// GenerateRegoRules asks Claude to generate Rego rules for the missing actions.
func (c *AnthropicClient) GenerateRegoRules(ctx context.Context, missing []string, goal, current map[string]interface{}) (string, error) {
	goalJSON, _ := json.MarshalIndent(goal, "", "  ")
	currentJSON, _ := json.MarshalIndent(current, "", "  ")
	missingJSON, _ := json.MarshalIndent(missing, "", "  ")

	prompt := fmt.Sprintf(`You are an OPA (Open Policy Agent) Rego expert. Generate Rego rules for missing actions.

Missing actions:
%s

Goal:
%s

Current state:
%s

Generate Rego rules in the "planner" package that define when each missing action should be taken.
Each rule should follow this pattern:

missing[action] {
  <conditions based on goal and current state>
  action := "<action_name>"
}

Return only valid Rego code, starting with "package planner". Do not include any explanation or markdown.`,
		string(missingJSON), string(goalJSON), string(currentJSON))

	msg, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(claudeModel),
		MaxTokens: anthropic.F(int64(2048)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		return "", fmt.Errorf("anthropic API error: %w", err)
	}

	if len(msg.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return msg.Content[0].Text, nil
}
