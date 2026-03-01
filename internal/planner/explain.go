package planner

import (
	"context"
	"fmt"
	"os"
	"strings"

	opaeval "github.com/hiroyukiosaki/opa-llm-planner/internal/opa"
)

// ExplainOptions holds options for the explain command.
type ExplainOptions struct {
	GoalFile     string
	CurrentFile  string
	PolicyDir    string
	WithLocation bool
}

// RunExplain evaluates missing actions with OPA tracing and prints the trace.
func RunExplain(ctx context.Context, opts ExplainOptions) error {
	goal, err := loadJSON(opts.GoalFile)
	if err != nil {
		return fmt.Errorf("reading goal: %w", err)
	}

	current, err := loadJSON(opts.CurrentFile)
	if err != nil {
		return fmt.Errorf("reading current: %w", err)
	}

	evaluator := opaeval.NewEvaluator(opts.PolicyDir)

	fmt.Println("=== OPA Trace: why actions are missing ===")
	fmt.Println()

	actions, err := evaluator.ExplainMissing(ctx, goal, current, os.Stdout, opts.WithLocation)
	if err != nil {
		return fmt.Errorf("OPA evaluation: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Result ===")
	if len(actions) == 0 {
		fmt.Println("No missing actions.")
	} else {
		fmt.Printf("Missing actions: [%s]\n", strings.Join(actions, ", "))
	}
	return nil
}
