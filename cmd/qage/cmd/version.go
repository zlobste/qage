package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zlobste/qage/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show the version, commit, and build date of qage.`,
	RunE:  runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Println(version.String())
	return nil
}
