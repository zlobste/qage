package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for qage.

To load completions:

Bash:
  $ source <(qage completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ qage completion bash > /etc/bash_completion.d/qage
  # macOS:
  $ qage completion bash > $(brew --prefix)/etc/bash_completion.d/qage

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ qage completion zsh > "${fpath[1]}/_qage"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ qage completion fish | source
  # To load completions for each session, execute once:
  $ qage completion fish > ~/.config/fish/completions/qage.fish

PowerShell:
  PS> qage completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> qage completion powershell > qage.ps1
  # and source this file from your PowerShell profile.`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:                  runCompletion,
}

func runCompletion(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
	case "zsh":
		return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
	case "fish":
		return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
	}
	return nil
}
