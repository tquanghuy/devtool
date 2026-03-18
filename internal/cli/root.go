package cli

import (
	"devtool/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devtool",
	Short: "A tool to manage developer services and configurations",
	Long: `devtool helps check the status of local databases (Postgres, MySQL)
and external services like Docker and Telepresence concurrently.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root command flags if needed
}
