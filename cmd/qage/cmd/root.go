// Package cmd provides the CLI commands for qage.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qage",
	Short: "Post-quantum hybrid age recipients",
	Long: `qage provides post-quantum secure recipients for age encryption.

It uses a hybrid X25519 + ML-KEM-768 key encapsulation mechanism to provide
security against both classical and quantum computers.

Examples:
  # Generate a new keypair
  qage keygen -o key.txt --comment "laptop"

  # Get the public recipient
  qage pub -i key.txt

  # Use with age
  age -R $(qage pub -i key.txt) -o secret.age secret.txt
  age -d -i key.txt secret.age`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(keygenCmd)
	rootCmd.AddCommand(pubCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(selftestCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(docsCmd)
}

// GetRootCmd returns the root command for testing.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// NewRootCmd creates a fresh root command for testing.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qage",
		Short: "Post-quantum hybrid age recipients",
		Long: `qage provides post-quantum secure recipients for age encryption.

It uses a hybrid X25519 + ML-KEM-768 key encapsulation mechanism to provide
security against both classical and quantum computers.

Examples:
  # Generate a new keypair
  qage keygen -o key.txt --comment "laptop"

  # Get the public recipient
  qage pub -i key.txt

  # Use with age
  age -R $(qage pub -i key.txt) -o secret.age secret.txt
  age -d -i key.txt secret.age`,
	}

	cmd.AddCommand(keygenCmd)
	cmd.AddCommand(pubCmd)
	cmd.AddCommand(inspectCmd)
	cmd.AddCommand(selftestCmd)
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(completionCmd)
	cmd.AddCommand(docsCmd)

	return cmd
}
