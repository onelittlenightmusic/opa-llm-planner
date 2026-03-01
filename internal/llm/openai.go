package llm

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/hiroyukiosaki/opa-llm-planner/internal/types"
)

// OpenAIClient implements LLMClient using OpenAI.
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates an OpenAIClient with the given API key.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{client: openai.NewClient(apiKey)}
}

// EnrichAction asks GPT-4o to fill in description and parameters for the action.
func (c *OpenAIClient) EnrichAction(ctx context.Context, action types.Action, goal, current map[string]interface{}) (types.Action, error) {
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

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens: 1024,
	})
	if err != nil {
		return action, fmt.Errorf("openai API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return action, fmt.Errorf("empty response from OpenAI")
	}

	responseText := resp.Choices[0].Message.Content

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

// GenerateRegoRules asks GPT-4o to generate Rego rules for the missing actions.
func (c *OpenAIClient) GenerateRegoRules(ctx context.Context, missing []string, goal, current map[string]interface{}) (string, error) {
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

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens: 2048,
	})
	if err != nil {
		return "", fmt.Errorf("openai API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
