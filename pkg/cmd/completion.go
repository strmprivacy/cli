package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(strm completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ strm completion bash > /etc/bash_completion.d/strm
  # macOS:
  $ strm completion bash > /usr/local/etc/bash_completion.d/strm

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ strm completion zsh > "${fpath[1]}/_strm"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ strm completion fish > ~/.config/fish/completions/strm.fish

PowerShell:

  PS> strm completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> strm completion powershell > strm.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}
