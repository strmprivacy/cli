package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/auth"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication command",
}

func init() {
	AuthCmd.AddCommand(auth.LoginCmd())
	AuthCmd.AddCommand(auth.PrintTokenCmd())
	AuthCmd.AddCommand(auth.RefreshCmd())
}
