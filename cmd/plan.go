package cmd

import (
	"context"

	"github.com/hiroyukiosaki/opa-llm-planner/internal/planner"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate an action plan from goal/current state using OPA policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		goalFile, _ := cmd.Flags().GetString("goal")
		currentFile, _ := cmd.Flags().GetString("current")
		policyDir, _ := cmd.Flags().GetString("policy")
		outFile, _ := cmd.Flags().GetString("out")
		useLLM, _ := cmd.Flags().GetBool("llm")
		llmProvider, _ := cmd.Flags().GetString("llm-provider")

		return planner.Run(context.Background(), planner.PlanOptions{
			GoalFile:    goalFile,
			CurrentFile: currentFile,
			PolicyDir:   policyDir,
			OutFile:     outFile,
			UseLLM:      useLLM,
			LLMProvider: llmProvider,
		})
	},
}

func init() {
	planCmd.Flags().String("goal", "examples/goal.json", "Path to goal JSON file")
	planCmd.Flags().String("current", "examples/current.json", "Path to current state JSON file")
	planCmd.Flags().String("policy", "policies", "Path to directory containing Rego policy files")
	planCmd.Flags().String("out", "", "Output file for the plan JSON (default: stdout)")
	planCmd.Flags().Bool("llm", false, "Use LLM to enrich actions with descriptions and parameters")
	planCmd.Flags().String("llm-provider", "", "LLM provider: anthropic or openai (overrides LLM_PROVIDER env var)")
}
