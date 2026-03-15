package cli

import (
	"fmt"

	"devtool/internal/manager"
	"github.com/spf13/cobra"
)

var addPort int
var addNonInteractive bool

var addCmd = &cobra.Command{
	Use:   "add <tool-name>",
	Short: "Add a supported tool to the managed list",
	Long: `Add a tool to devtool's managed list.

Singleton tools (e.g. docker, telepresence) may only be added once.
Port-bound tools (e.g. postgres, mysql) can be added multiple times on different ports.

Examples:
  devtool add docker
  devtool add postgres
  devtool add postgres --port 5433`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName := args[0]
		if err := manager.AddTool(toolName, addPort, addNonInteractive); err != nil {
			return err
		}
		fmt.Printf("Successfully added %s to managed list.\n", toolName)
		return nil
	},
}

func init() {
	addCmd.Flags().IntVar(&addPort, "port", 0, "Specifically request a port for port-bound tools. Ignored or causes error for singletons.")
	addCmd.Flags().BoolVar(&addNonInteractive, "non-interactive", false, "Fails instead of prompting if port mapping resolution is required for port-bound tools.")
	rootCmd.AddCommand(addCmd)
}
