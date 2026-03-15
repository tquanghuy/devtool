package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"devtool/internal/checker"
	"devtool/internal/config"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of all configured developer tools.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		orchestrator := checker.NewOrchestrator(cfg)
		orchestrator.Register(&checker.PostgresChecker{})
		orchestrator.Register(&checker.MySQLChecker{})
		orchestrator.Register(&checker.DockerChecker{})
		orchestrator.Register(&checker.TelepresenceChecker{})

		// Context with 3-second timeout SC-003
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		fmt.Println("Checking services...")
		results := orchestrator.RunAll(ctx)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "TOOL\tSTATUS\tLATENCY\tMESSAGE")

		for _, res := range results {
			statusStr := "UP"
			errorStr := "-"
			if !res.IsUp {
				statusStr = "DOWN"
				if res.Error != nil {
					errorStr = res.Error.Error()
				} else {
					errorStr = "timeout or unknown error"
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", res.Name, statusStr, res.Duration.Round(time.Millisecond).String(), errorStr)
		}
		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
