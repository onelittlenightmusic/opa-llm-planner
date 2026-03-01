package planner

import (
	"context"
	"fmt"
	"os"

	"github.com/hiroyukiosaki/opa-llm-planner/internal/llm"
	opaeval "github.com/hiroyukiosaki/opa-llm-planner/internal/opa"
)

// ConsiderOptions holds options for the consider command.
type ConsiderOptions struct {
	GoalFile    string
	CurrentFile string
	PolicyDir   string
	OutFile     string
	Append      bool
	DryRun      bool
	LLMProvider string
}

// RunConsider evaluates missing actions, generates Rego rules via LLM, validates them,
// and writes (or prints) the result.
func RunConsider(ctx context.Context, opts ConsiderOptions) error {
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

	if len(missing) == 0 {
		fmt.Println("No missing actions found; no new rules needed.")
		return nil
	}

	fmt.Printf("Missing actions: %v\n", missing)

	client, err := llm.NewClient(opts.LLMProvider)
	if err != nil {
		return fmt.Errorf("creating LLM client: %w", err)
	}

	regoSrc, err := client.GenerateRegoRules(ctx, missing, goal, current)
	if err != nil {
		return fmt.Errorf("generating Rego rules: %w", err)
	}

	// Validate the generated Rego
	if err := opaeval.ValidateRego(ctx, "generated.rego", regoSrc); err != nil {
		return fmt.Errorf("generated Rego is invalid: %w\n\nGenerated source:\n%s", err, regoSrc)
	}

	if opts.DryRun {
		fmt.Println("--- Generated Rego (dry-run) ---")
		fmt.Println(regoSrc)
		return nil
	}

	if opts.OutFile == "" || opts.OutFile == "-" {
		fmt.Println(regoSrc)
		return nil
	}

	if opts.Append {
		f, err := os.OpenFile(opts.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("opening output file: %w", err)
		}
		defer f.Close()
		if _, err := fmt.Fprintln(f, regoSrc); err != nil {
			return fmt.Errorf("writing rules: %w", err)
		}
	} else {
		if err := os.WriteFile(opts.OutFile, []byte(regoSrc), 0644); err != nil {
			return fmt.Errorf("writing rules: %w", err)
		}
	}

	fmt.Printf("Rego rules written to %s\n", opts.OutFile)
	return nil
}
