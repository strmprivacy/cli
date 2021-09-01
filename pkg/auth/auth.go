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
	"context"
	"fmt"
	"github.com/int128/oauth2cli"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
)

var TokenFile string

const (
	EventAuthHostFlag = "event-auth-host"
	ApiAuthUrlFlag    = "api-auth-url"
	PasswordFlag      = "password"
)

func login() {
	oautLogin()
	/*
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
	*/
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

func oautLogin() {

	ready := make(chan string, 1)
	defer close(ready)
	cfg := oauth2cli.Config{
		OAuth2Config: oauth2.Config{
			ClientID:     "k04TCdBzSedw3JXDpCg0LejEU7kf0nch",
			ClientSecret: "8Hw5SbyGF-jrhyJSHOfGsU5KWGBpXhuQ-YWczNI2bBy6nICbkaM3qxIyflEGHIXh",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://dev-p6wb4vof.eu.auth0.com/authorize",
				TokenURL: "https://dev-p6wb4vof.eu.auth0.com/oauth/token",
			},
			Scopes: []string{"offline_access", "profile", "email", "openid"},
		},
		AuthCodeOptions: []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("audience", "https://apis.streammachine.io"),
		},
		//TokenRequestOptions:    pkce.TokenRequestOptions(),
		LocalServerReadyChan:   ready,
		LocalServerCertFile:    "",
		LocalServerKeyFile:     "",
		LocalServerBindAddress: strings.Split("127.0.0.1:10000", ","),
		Logf:                   log.Printf,
	}

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			log.Printf("Open %s", url)
			if err := browser.OpenURL(url); err != nil {
				log.Printf("could not open the browser: %s", err)
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		token, err := oauth2cli.GetToken(ctx, cfg)
		if err != nil {
			return fmt.Errorf("could not get a token: %w", err)
		}
		log.Printf("You got a valid token until %s", token.Expiry)
		log.Printf("Access Token:\n\n%v\n", token.AccessToken)
		log.Printf("Refresh Token:\n\n%v\n", token.RefreshToken)
		log.Printf("Raw:\n\n%v\n", token)
		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Fatalf("authorization error: %s", err)
	}
}
