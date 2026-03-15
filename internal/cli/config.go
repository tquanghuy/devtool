package cli

import (
	"fmt"
	"syscall"

	"devtool/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage devtool configuration",
}

var setCredentialCmd = &cobra.Command{
	Use:   "set-credential [user]",
	Short: "Safely store a database password in the native keychain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user := args[0]
		fmt.Printf("Enter password for %s: ", user)
		bytepw, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return err
		}
		if err := config.SetPassword(user, string(bytepw)); err != nil {
			return fmt.Errorf("failed to save password: %w", err)
		}
		fmt.Printf("Password for user '%s' saved successfully in keychain.\n", user)
		return nil
	},
}

func init() {
	configCmd.AddCommand(setCredentialCmd)
	rootCmd.AddCommand(configCmd)
}
