package cmd

import (
	"context"

	"github.com/hiroyukiosaki/opa-llm-planner/internal/planner"
	"github.com/spf13/cobra"
)

var considerCmd = &cobra.Command{
	Use:   "consider",
	Short: "Generate new Rego rules for missing actions using LLM",
	RunE: func(cmd *cobra.Command, args []string) error {
		goalFile, _ := cmd.Flags().GetString("goal")
		currentFile, _ := cmd.Flags().GetString("current")
		policyDir, _ := cmd.Flags().GetString("policy")
		outFile, _ := cmd.Flags().GetString("out")
		appendMode, _ := cmd.Flags().GetBool("append")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		llmProvider, _ := cmd.Flags().GetString("llm-provider")

		return planner.RunConsider(context.Background(), planner.ConsiderOptions{
			GoalFile:    goalFile,
			CurrentFile: currentFile,
			PolicyDir:   policyDir,
			OutFile:     outFile,
			Append:      appendMode,
			DryRun:      dryRun,
			LLMProvider: llmProvider,
		})
	},
}

func init() {
	considerCmd.Flags().String("goal", "examples/goal.json", "Path to goal JSON file")
	considerCmd.Flags().String("current", "examples/current.json", "Path to current state JSON file")
	considerCmd.Flags().String("policy", "policies", "Path to directory containing Rego policy files")
	considerCmd.Flags().String("out", "", "Output file for generated Rego rules (default: stdout)")
	considerCmd.Flags().Bool("append", false, "Append generated rules to the output file instead of overwriting")
	considerCmd.Flags().Bool("dry-run", false, "Print generated rules without writing to file")
	considerCmd.Flags().String("llm-provider", "", "LLM provider: anthropic or openai (overrides LLM_PROVIDER env var)")
}
