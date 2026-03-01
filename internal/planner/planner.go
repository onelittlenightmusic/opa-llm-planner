package planner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hiroyukiosaki/opa-llm-planner/internal/llm"
	opaeval "github.com/hiroyukiosaki/opa-llm-planner/internal/opa"
	"github.com/hiroyukiosaki/opa-llm-planner/internal/types"
)

// PlanOptions holds options for the plan command.
type PlanOptions struct {
	GoalFile    string
	CurrentFile string
	PolicyDir   string
	OutFile     string
	UseLLM      bool
	LLMProvider string
}

// Run executes the plan logic and writes the result to OutFile.
func Run(ctx context.Context, opts PlanOptions) error {
	goal, err := loadJSON(opts.GoalFile)
	if err != nil {
		return fmt.Errorf("reading goal: %w", err)
	}

	current, err := loadJSON(opts.CurrentFile)
	if err != nil {
		return fmt.Errorf("reading current: %w", err)
	}

	evaluator := opaeval.NewEvaluator(opts.PolicyDir)
	missing, err := evaluator.EvaluateMissing(ctx, goal, current)
	if err != nil {
		return fmt.Errorf("OPA evaluation: %w", err)
	}

	actions := make([]types.Action, 0, len(missing))
	for _, name := range missing {
		actions = append(actions, types.Action{
			Type:   name,
			Status: "pending",
		})
	}

	if opts.UseLLM && len(actions) > 0 {
		client, err := llm.NewClient(opts.LLMProvider)
		if err != nil {
			return fmt.Errorf("creating LLM client: %w", err)
		}
		for i, action := range actions {
			enriched, err := client.EnrichAction(ctx, action, goal, current)
			if err != nil {
				return fmt.Errorf("enriching action %q: %w", action.Type, err)
			}
			actions[i] = enriched
		}
	}

	goalID, _ := goal["id"].(string)
	plan := types.Plan{
		PlanID:  uuid.New().String(),
		GoalID:  goalID,
		Actions: actions,
	}

	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling plan: %w", err)
	}

	if opts.OutFile == "" || opts.OutFile == "-" {
		fmt.Println(string(data))
		return nil
	}

	if err := os.WriteFile(opts.OutFile, data, 0644); err != nil {
		return fmt.Errorf("writing plan: %w", err)
	}
	fmt.Printf("Plan written to %s\n", opts.OutFile)
	return nil
}

func loadJSON(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
