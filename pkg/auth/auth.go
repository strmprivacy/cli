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

var TokenFile string

const (
	EventAuthHostFlag = "event-auth-host"
	ApiAuthUrlFlag    = "api-auth-url"
	PasswordFlag      = "password"
)

func login(s *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	password, _ := flags.GetString(PasswordFlag)
	if password == "" {
		password = askPassword()
		fmt.Println()
	}

	Client.AuthenticateLogin(s, &password)
	_, billingId := Client.GetToken(false)
	fmt.Println("Billing id:", billingId)
	filename := Client.StoreLogin()
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

func printAccessToken() {
	Client.printToken()
}

func Refresh() {
	Client.LoadLogin()
	Client.refresh()
	Client.StoreLogin()
}
