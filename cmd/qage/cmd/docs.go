package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate CLI documentation",
	Long:   `Generate CLI documentation in markdown format.`,
	Hidden: true, // Hide from normal help, mainly for maintainers
	RunE:   runDocs,
}

var docsOutput string

func init() {
	docsCmd.Flags().StringVarP(&docsOutput, "output", "o", "./docs", "output directory for docs")
}

func runDocs(cmd *cobra.Command, args []string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(docsOutput, 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	// Generate markdown docs
	err := doc.GenMarkdownTree(cmd.Root(), docsOutput)
	if err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}

	fmt.Printf("Documentation generated in %s\n", docsOutput)

	// List generated files
	entries, err := os.ReadDir(docsOutput)
	if err != nil {
		return err
	}

	fmt.Println("Generated files:")
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".md" {
			fmt.Printf("  %s\n", entry.Name())
		}
	}

	return nil
}
