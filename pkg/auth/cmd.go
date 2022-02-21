package auth

import (
	"github.com/spf13/cobra"
	"io/ioutil"
)

func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login",
		Long: `Log a user in using its Console credentials and save the login token to disk, 
to allow the CLI access to the STRM Privacy APIs.`,
		Run: func(cmd *cobra.Command, args []string) {
			login()
		},
		DisableAutoGenTag: true,
		Args:              cobra.ExactArgs(0),
	}
	return cmd
}
func PrintTokenCmd() *cobra.Command {
	var content, _ = ioutil.ReadFile("pkg/auth/docstring_print_token.md")
	cmd := &cobra.Command{
		Use:   "print-access-token",
		Short: "Print your current access-token to stdout",
		Long:  string(content),
		Run: func(cmd *cobra.Command, args []string) {
			printAccessToken()
		},
		DisableAutoGenTag: true,
	}
	return cmd
}
