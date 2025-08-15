package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/zlobste/qage/pkg/qage"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new qage identity",
	Long: `Generate a new qage identity with X25519 + ML-KEM-768 hybrid keys.

The identity will be printed to stdout unless -o is specified.`,
	Example: `  # Generate a key to stdout
  qage keygen --comment "laptop"

  # Generate a key to file
  qage keygen -o ~/.qage/key --comment "laptop"`,
	RunE: runKeygen,
}

var (
	keygenOutput  string
	keygenComment string
)

func init() {
	keygenCmd.Flags().StringVarP(&keygenOutput, "output", "o", "", "output file (default: stdout)")
	keygenCmd.Flags().StringVarP(&keygenComment, "comment", "c", "", "comment for the key")
}

func runKeygen(cmd *cobra.Command, args []string) error {
	// Generate new identity
	identity, err := qage.NewIdentity()
	if err != nil {
		return fmt.Errorf("failed to generate identity: %w", err)
	}

	// Format for file
	formatted, err := identity.FormatFile(keygenComment)
	if err != nil {
		return fmt.Errorf("failed to format identity: %w", err)
	}

	// Output
	w := cmd.OutOrStdout()
	if keygenOutput != "" {
		// Ensure directory exists
		if mkdirErr := os.MkdirAll(filepath.Dir(keygenOutput), 0700); mkdirErr != nil {
			return fmt.Errorf("failed to create directory: %w", mkdirErr)
		}

		f, err := os.OpenFile(keygenOutput, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() {
			if closeErr := f.Close(); closeErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
			}
		}()
		w = f
	}

	_, err = fmt.Fprintln(w, formatted)
	return err
}
