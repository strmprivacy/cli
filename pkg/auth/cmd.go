package auth

import (
	"github.com/spf13/cobra"
)

func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login",
		Run: func(cmd *cobra.Command, args []string) {
			login()
		},
		Args: cobra.ExactArgs(0),
	}
	return cmd
}

func PrintTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print-access-token",
		Short: "Print your current access-token to stdout",
		Long: `Prints an access token that can be used in an http header.
Note that this token might be expired, so a refresh may be required.
Use token as follows:
'Authorization: Bearer <token>'
`,
		Run: func(cmd *cobra.Command, args []string) {
			printAccessToken()
		},
	}
	return cmd
}
