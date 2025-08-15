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

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Show identity metadata",
	Long: `Show metadata about a qage identity including the cryptographic suite,
key lengths, and public recipient.`,
	Example: `  # Inspect from file
  qage inspect -i ~/.qage/key

  # Inspect from stdin
  cat ~/.qage/key | qage inspect`,
	RunE: runInspect,
}

var inspectIdentity string

func init() {
	inspectCmd.Flags().StringVarP(&inspectIdentity, "identity", "i", "-", "identity file ('-' for stdin)")
}

func runInspect(cmd *cobra.Command, args []string) error {
	identity, comment, err := readIdentity()
	if err != nil {
		return err
	}

	return displayIdentityInfo(identity, comment)
}

func readIdentity() (*qage.Identity, string, error) {
	// Read identity
	var r io.Reader = os.Stdin
	if inspectIdentity != "-" {
		f, err := os.Open(inspectIdentity)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open identity file: %w", err)
		}
		defer func() {
			if closeErr := f.Close(); closeErr != nil {
				// Log error but don't override main error
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
		return nil, "", fmt.Errorf("failed to read identity: %w", err)
	}
	if identityLine == "" {
		return nil, "", fmt.Errorf("no identity found in input")
	}

	// Parse identity
	var identity *qage.Identity
	var comment string
	var err error

	if strings.HasPrefix(identityLine, "QAGE-SECRET-KEY-1 ") {
		// File format
		identity, comment, err = qage.ParseIdentityFile(identityLine)
	} else {
		// Direct bech32
		identity, err = qage.ParseIdentity(identityLine)
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to parse identity: %w", err)
	}

	return identity, comment, nil
}

func displayIdentityInfo(identity *qage.Identity, comment string) error {
	// Show metadata
	fmt.Printf("Type: qage identity\n")
	fmt.Printf("Suite: %s\n", identity.Suite())
	if comment != "" {
		fmt.Printf("Comment: %s\n", comment)
	}

	// Get recipient
	recipient := identity.Recipient()
	recipientStr, err := recipient.String()
	if err != nil {
		return fmt.Errorf("failed to encode recipient: %w", err)
	}

	fmt.Printf("Public recipient: %s\n", recipientStr)
	fmt.Printf("Recipient length: %d characters\n", len(recipientStr))

	// Get identity string
	identityStr, err := identity.String()
	if err != nil {
		return fmt.Errorf("failed to encode identity: %w", err)
	}
	fmt.Printf("Identity length: %d characters\n", len(identityStr))

	return nil
}
