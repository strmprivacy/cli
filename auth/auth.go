// Interface with the Streammachine Security Token provider.
//
// Usage:
// authClient := &eventauth.Auth{Uri: sts}
// authClient.Authenticate(billingId, clientId, clientSecret)
// authClient.getToken() --> provides working token that needs to be put in an http
// Authorization: Bearer <token>
// header

package auth

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
)

var RootCommandName string
var TokenFile string

const (
	EventAuthHostFlag = "event-auth-host"
	ApiAuthUrlFlag    = "api-auth-url"
	PasswordFlag      = "password"
)

func login(apiHost string, s *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	password, _ := flags.GetString(PasswordFlag)
	if password == "" {
		password = askPassword()
		fmt.Println()
	}

	authClient := &Auth{Uri: apiHost}
	authClient.AuthenticateLogin(s, &password)
	_, billingId := authClient.GetToken(false)
	fmt.Println("Billing id:", billingId)
	filename := authClient.StoreLogin()
	fmt.Println("Saved login to:", filename)
}

func askPassword() string {
	fmt.Print("Enter password: ")
	fd := int(os.Stdin.Fd())
	state, _ := term.MakeRaw(fd)
	defer term.Restore(fd, state)
	pwBytes, _ := term.ReadPassword(fd)
	return string(pwBytes)
}

func printAccessToken(apiHost string) {
	authClient := &Auth{Uri: apiHost}
	authClient.LoadLogin()
	authClient.printToken()
}

func DoRefresh(apiHost string) {
	authClient := &Auth{Uri: apiHost}
	authClient.LoadLogin()
	authClient.refresh()
	authClient.StoreLogin()
}
