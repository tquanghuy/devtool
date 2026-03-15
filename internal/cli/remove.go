package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"devtool/internal/manager"
	"github.com/spf13/cobra"
)

var removeForce bool
var removeNonInteractive bool

var removeCmd = &cobra.Command{
	Use:   "remove <tool-name>|<instance-id>",
	Short: "Remove a tool from the managed list",
	Long: `Remove a tool or specific tool instance from devtool's managed list.

If multiple instances of a port-bound tool exist, you will be prompted to select one.
Use --force to terminate a running tool before removal without prompting.

Examples:
  devtool remove docker
  devtool remove postgres
  devtool remove postgres-5433
  devtool remove postgres --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName := args[0]

		selectFn := func(identifiers []string) (int, error) {
			fmt.Println("Multiple instances found:")
			for i, id := range identifiers {
				fmt.Printf("  %d) %s\n", i+1, id)
			}
			fmt.Printf("Select instance to remove [1-%d]: ", len(identifiers))

			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil {
				return -1, err
			}
			line = strings.TrimSpace(line)
			n, err := strconv.Atoi(line)
			if err != nil || n < 1 || n > len(identifiers) {
				return -1, fmt.Errorf("invalid selection: %q", line)
			}
			return n - 1, nil
		}

		confirmFn := func(name string) (bool, error) {
			fmt.Printf("Tool %s is running. Force stop it? [y/N]: ", name)
			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil {
				return false, err
			}
			line = strings.TrimSpace(strings.ToLower(line))
			if line == "y" || line == "yes" {
				return true, nil
			}
			return false, nil
		}

		return manager.RemoveTool(toolName, removeForce, removeNonInteractive, selectFn, confirmFn)
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "Bypasses the prompt to gracefully stop the tool and forcefully terminates it before removal.")
	removeCmd.Flags().BoolVar(&removeNonInteractive, "non-interactive", false, "Fails if the tool is currently running (requires manual stop first) or if multiple instances exist without a specific <instance-id> provided.")
	rootCmd.AddCommand(removeCmd)
}
