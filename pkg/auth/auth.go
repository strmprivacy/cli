package auth

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

var TokenFile string

const (
	ApiAuthUrlFlag = "api-auth-url"
)

func login(cmd *cobra.Command) {
	Auth.login(cmd)
}

func revoke() {
	Auth.revoke()
}

func printAccessToken() {
	Auth.printAccessToken()
}

func RequireAuthenticationPreRun(cmd *cobra.Command, args []string) {
	_ = cmd.Root().PersistentPreRunE(cmd, args)
	accessToken := Auth.GetToken()

	if accessToken == nil {
		common.UnauthenticatedErrorWithExit()
	}
}
