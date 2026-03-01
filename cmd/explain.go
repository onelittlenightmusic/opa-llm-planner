package cmd

import (
	"context"

	"github.com/hiroyukiosaki/opa-llm-planner/internal/planner"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Show why actions are missing using OPA trace",
	Long: `Evaluate the OPA policy with tracing enabled and print the full trace,
showing exactly which rules fired (or did not fire) and why each action is missing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		goalFile, _ := cmd.Flags().GetString("goal")
		currentFile, _ := cmd.Flags().GetString("current")
		policyDir, _ := cmd.Flags().GetString("policy")
		withLocation, _ := cmd.Flags().GetBool("location")

		return planner.RunExplain(context.Background(), planner.ExplainOptions{
			GoalFile:     goalFile,
			CurrentFile:  currentFile,
			PolicyDir:    policyDir,
			WithLocation: withLocation,
		})
	},
}

func init() {
	explainCmd.Flags().String("goal", "examples/goal.json", "Path to goal JSON file")
	explainCmd.Flags().String("current", "examples/current.json", "Path to current state JSON file")
	explainCmd.Flags().String("policy", "policies", "Path to directory containing Rego policy files")
	explainCmd.Flags().Bool("location", false, "Include source file location in trace output")
}
