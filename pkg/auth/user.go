package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/int128/oauth2cli"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"os"
	"strings"
	"strmprivacy/strm/pkg/common"
)

var Auth = Authenticator{}

type Authenticator struct {
	storedToken storedToken
	tokenSource oauth2.TokenSource
	Email       string
}

func oAuth2Config() oauth2.Config {
	return oauth2.Config{
		ClientID: "cli",
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%v/auth/auth", common.ApiAuthHost),
			TokenURL: fmt.Sprintf("%v/auth/token", common.ApiAuthHost),
		},
		Scopes: []string{"offline_access", "email"},
	}
}

func (authenticator *Authenticator) GetToken() *string {
	return authenticator.accessToken()
}

func (authenticator *Authenticator) printAccessToken() {
	accessToken := authenticator.accessToken()

	if accessToken != nil {
		fmt.Println(*accessToken)
	} else {
		common.MissingIdTokenError()
	}
}

func (authenticator *Authenticator) accessToken() *string {
	if authenticator.tokenSource == nil {
		return nil
	} else {
		tokens, err := authenticator.tokenSource.Token()
		if err != nil {
			authenticator.revoke()
			common.CliExit(errors.New(fmt.Sprintf("Your session has expired. Please re-login using: `%s auth login`", common.RootCommandName)))
		}
		if authenticator.storedToken.AccessToken != tokens.AccessToken {
			authenticator.populateValues(oauthTokenToStoredToken(*tokens))
			authenticator.storeLogin()
		}
		return &tokens.AccessToken
	}
}

func (authenticator *Authenticator) revoke() {
	filename := authenticator.getSaveFilename()
	err := os.Remove(filename)
	common.CliExit(err)
}

func (authenticator *Authenticator) login() {
	ready := make(chan string, 1)
	defer close(ready)

	port := findFreePort()

	cfg := oauth2cli.Config{
		OAuth2Config:           oAuth2Config(),
		LocalServerReadyChan:   ready,
		LocalServerBindAddress: strings.Split(fmt.Sprintf("127.0.0.1:%d", port), ","),
		SuccessRedirectURL:     "https://strmprivacy.io/auth-success",
		FailureRedirectURL:     "https://strmprivacy.io/auth-failure",
		AuthCodeOptions: []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("prompt", "login"),
		},
	}

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(startBrowserLoginFlow(ready, ctx))
	eg.Go(authenticator.handleLogin(ctx, cfg))

	if err := eg.Wait(); err != nil {
		common.CliExit(errors.New(fmt.Sprintf("Login failed, please check the logs for details at %v", common.LogFileName())))
	}
}

func (authenticator *Authenticator) handleLogin(ctx context.Context, cfg oauth2cli.Config) func() error {
	return func() error {
		oAuthToken, err := oauth2cli.GetToken(ctx, cfg)
		common.CliExit(err)

		authenticator.populateValues(oauthTokenToStoredToken(*oAuthToken))
		authenticator.storeLogin()

		fmt.Println(fmt.Sprintf("\nYou are now logged in as [%v].", authenticator.Email))

		return nil
	}
}

func oauthTokenToStoredToken(t oauth2.Token) storedToken {
	return storedToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    t.Expiry.Unix(),
		Email:        getEmailFromClaims(t),
	}
}

func startBrowserLoginFlow(ready chan string, ctx context.Context) func() error {
	return func() error {
		select {
		case localCallbackServerUrl := <-ready:
			headless := os.Getenv("STRM_HEADLESS")

			if headless != "true" {
				err := browser.OpenURL(localCallbackServerUrl)

				if err != nil {
					browserOpenError := errors.New(fmt.Sprintf("Unable to open browser for authentication: %s", err))
					log.Error(browserOpenError)
					common.CliExit(browserOpenError)
				}
			}

			fmt.Println("Follow the login flow in your browser, which is opened automatically. If not, open the following URL to complete the login:")
			fmt.Println(fmt.Sprintf("\n    %v", authorizationCodeFlowUrl(localCallbackServerUrl)))

			return nil
		case <-ctx.Done():
			contextError := fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
			log.Error(contextError)
			common.CliExit(contextError)
			return nil
		}
	}
}

func getEmailFromClaims(t oauth2.Token) string {
	parsedAccessToken, _ := jwt.Parse(t.AccessToken, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	claims := parsedAccessToken.Claims.(jwt.MapClaims)
	return (claims["email"]).(string)
}

func authorizationCodeFlowUrl(url string) string {
	locationResponse, err := http.Get(url)
	if err != nil {
		common.CliExit(errors.New("Could not retrieve authorization code flow URL. Please retry or contact STRM Privacy support if the problem persists."))
	}

	return locationResponse.Request.URL.String()
}

func findFreePort() int {
	headless := os.Getenv("STRM_HEADLESS")

	if headless != "true" {
		port := 10000
		foundOpenPort := false

		for port < 10010 {
			host := fmt.Sprintf("127.0.0.1:%d", port)

			ln, err := net.Listen("tcp", host)
			if err != nil {
				log.Infof(fmt.Sprintf("Can't listen on port %d: %s", port, err))
				port = port + 1
				continue
			}

			_ = ln.Close()
			foundOpenPort = true
			break
		}

		if !foundOpenPort {
			common.CliExit(errors.New("Unable to find free port in range 10000 <= port <= 10009. Please check your running applications and make sure that a port in this range is free."))
		}

		return port
	} else {
		return 10000
	}
}
