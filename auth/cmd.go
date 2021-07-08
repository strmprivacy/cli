package auth

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/utils"
)

func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [email]",
		Short: "Login",
		Run: func(cmd *cobra.Command, args []string) {
			login(apiHost(cmd), &args[0], cmd)
		},
		Args: cobra.ExactArgs(1), // the stream name
	}
	flags := cmd.Flags()
	flags.String(PasswordFlag, "", "password")
	return cmd
}

func PrintTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access-token",
		Short: "print an access-token to stdout",
		Long: `Prints an access token that can be used in an http header.
Use token as follows:
'Authorization: Bearer <token>'
`,
		Run: func(cmd *cobra.Command, args []string) {
			printAccessToken(apiHost(cmd))
		},
	}
	return cmd
}

func Refresh() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh an existing access-token",
		Long: `Not really necessary, the cli will auto-refresh.
`,
		Run: func(cmd *cobra.Command, args []string) {
			DoRefresh(apiHost(cmd))
		},
	}
	return cmd
}

func apiHost(cmd *cobra.Command) string {
	return utils.GetStringAndErr(cmd.Flags(), ApiAuthUrlFlag)
}
