package auth

import (
	"fmt"
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

const (
	nonInteractiveTargetHostFlag      = "non-interactive"
	nonInteractiveTargetHostShortFlag = "n"
	nonInteractiveRemoteHostFlag      = "remote"
	nonInteractiveRemoteHostShortFlag = "r"
)

var longDocPrintToken = util.LongDocsUsage(`
Print the current (JWT) access token to the terminal that can be used in a http header. Note that the token is printed
on °stdout°, and the Expiry on °stderr° so it’s easy to capture the token for scripting use with

°°°bash
export token=$(strm auth print-access-token)
°°°

Note that this token might be expired, so a refresh may be required. Use token as follows:
'Authorization: Bearer &lt;token&gt;'

`)

func LoginCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login",
		Long: `Log a user in using its Console credentials and save the login token to disk,
to allow the CLI access to the STRM Privacy APIs.`,
		Run: func(cmd *cobra.Command, args []string) {
			login(cmd)
		},
		DisableAutoGenTag: true,
		Args:              cobra.ExactArgs(0),
	}

	flags := loginCmd.Flags()
	flags.BoolP(nonInteractiveTargetHostFlag, nonInteractiveTargetHostShortFlag, false, fmt.Sprintf("is the current host a headless system, without access to a browser? If true, the login process will wait for an authorization code flow result, that can be retrieved by using `%s auth login --%s`", common.RootCommandName, nonInteractiveRemoteHostFlag))
	flags.BoolP(nonInteractiveRemoteHostFlag, nonInteractiveRemoteHostShortFlag, false, "should the current host act as a remote login for a headless system? If true, an authorization code flow result will be printed, that can be used for the non-interactive target host.")

	loginCmd.MarkFlagsMutuallyExclusive(nonInteractiveTargetHostFlag, nonInteractiveRemoteHostFlag)

	return loginCmd
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
			if Auth.Email != "" {
				fmt.Println(fmt.Sprintf("Currently logged in as [%v]", Auth.Email))
			} else {
				common.UnauthenticatedErrorWithExit()
			}
		},
		DisableAutoGenTag: true,
	}
	return cmd
}
