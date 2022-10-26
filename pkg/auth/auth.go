package auth

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

var TokenFile string

const (
	EventsAuthUrlFlag = "events-auth-url"
	ApiAuthUrlFlag    = "api-auth-url"
)

func login() {
	Auth.login()
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
