package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "opa-llm-planner",
	Short: "Generate execution plans from OPA policies and LLM",
	Long: `opa-llm-planner combines OPA (Open Policy Agent) Rego rules with LLM
to generate action plans from goal and current state JSON files.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(considerCmd)
}
