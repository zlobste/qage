package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zlobste/qage/pkg/qage"
)

var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Extract public recipient from identity",
	Long: `Extract the public recipient string from a qage identity.

The recipient can be used with age -R or age --recipient.`,
	Example: `  # From file
  qage pub -i ~/.qage/key

  # From stdin
  cat ~/.qage/key | qage pub

  # Use with age
  age -R $(qage pub -i ~/.qage/key) -o secret.age secret.txt`,
	RunE: runPub,
}

var pubIdentity string

func init() {
	pubCmd.Flags().StringVarP(&pubIdentity, "identity", "i", "-", "identity file ('-' for stdin)")
}

func runPub(cmd *cobra.Command, args []string) error {
	// Read identity
	var r io.Reader = os.Stdin
	if pubIdentity != "-" {
		f, err := os.Open(pubIdentity)
		if err != nil {
			return fmt.Errorf("failed to open identity file: %w", err)
		}
		defer func() {
			if closeErr := f.Close(); closeErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
			}
		}()
		r = f
	}

	// Read first non-empty, non-comment line
	scanner := bufio.NewScanner(r)
	var identityLine string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			identityLine = line
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read identity: %w", err)
	}
	if identityLine == "" {
		return fmt.Errorf("no identity found in input")
	}

	// Parse identity
	var identity *qage.Identity
	var err error

	if strings.HasPrefix(identityLine, "QAGE-SECRET-KEY-1 ") {
		// File format
		identity, _, err = qage.ParseIdentityFile(identityLine)
	} else {
		// Direct bech32
		identity, err = qage.ParseIdentity(identityLine)
	}

	if err != nil {
		return fmt.Errorf("failed to parse identity: %w", err)
	}

	// Get recipient
	recipient := identity.Recipient()
	recipientStr, err := recipient.String()
	if err != nil {
		return fmt.Errorf("failed to encode recipient: %w", err)
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), recipientStr)
	return err
}
