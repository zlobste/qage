package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zlobste/qage/pkg/qage"
)

var selftestCmd = &cobra.Command{
	Use:   "selftest",
	Short: "Run internal validation tests",
	Long: `Run internal validation tests to verify that qage is working correctly.

This tests key generation, encoding/decoding, and age integration.
No output means all tests passed.`,
	Example: `  qage selftest`,
	RunE:    runSelftest,
}

func runSelftest(cmd *cobra.Command, args []string) error {
	if err := qage.Selftest(); err != nil {
		return fmt.Errorf("selftest failed: %w", err)
	}
	return nil
}
