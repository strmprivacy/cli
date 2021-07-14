package auth

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/util"
)

func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [email]",
		Short: "Login",
		Run: func(cmd *cobra.Command, args []string) {
			login(&args[0], cmd)
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

func RefreshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "RefreshCmd an existing access-token",
		Long: `Not really necessary, the CLI will auto-refresh.
`,
		Run: func(cmd *cobra.Command, args []string) {
			Refresh()
		},
	}
	return cmd
}

func apiHost(cmd *cobra.Command) string {
	return util.GetStringAndErr(cmd.Flags(), ApiAuthUrlFlag)
}
