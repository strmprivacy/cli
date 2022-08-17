package auth

import (
	"fmt"
	"github.com/spf13/cobra"
)

var longDocPrintToken = `
Print the current (JWT) access token to the terminal that can be used in a http header. Note that the token is printed
on ` + "`stdout`" + `, and the Expiry on ` + "`stderr`" + ` so it’s easy to capture the token for scripting use with

` + "```" + `bash
export token=$(strm auth access-token)
` + "```" + `

Note that this token might be expired, so a refresh may be required. Use token as follows:
'Authorization: Bearer &lt;token&gt;'

### Usage
`

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

func RevokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke",
		Long:  `Revoke your current login session and stored credentials.`,
		Run: func(cmd *cobra.Command, args []string) {
			revoke()
		},
		DisableAutoGenTag: true,
		Args:              cobra.ExactArgs(0),
	}
	return cmd
}

func PrintTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print-access-token",
		Short: "Print your current access-token to stdout",
		Long:  longDocPrintToken,
		Run: func(cmd *cobra.Command, args []string) {
			printAccessToken()
		},
		DisableAutoGenTag: true,
	}
	return cmd
}

func ShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show your current login credentials",
		Long:  `Show the email address of your current login credentials`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Currently logged in as [%v]", Auth.Email))
		},
		DisableAutoGenTag: true,
	}
	return cmd
}
