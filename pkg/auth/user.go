package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trietsch/oauth2cli"
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
	return authenticator.AccessToken()
}

func (authenticator *Authenticator) printAccessToken() {
	accessToken := authenticator.AccessToken()

	if accessToken != nil {
		fmt.Println(*accessToken)
	} else {
		common.UnauthenticatedErrorWithExit()
	}
}

func (authenticator *Authenticator) AccessToken() *string {
	if authenticator.tokenSource == nil {
		return nil
	} else {
		tokens, err := authenticator.tokenSource.Token()
		if err != nil {
			authenticator.revoke()
			common.CliExit(errors.New(fmt.Sprintf("Your session has expired. Please re-login using: %s auth login", common.RootCommandName)))
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

func (authenticator *Authenticator) login(cmd *cobra.Command) {
	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	baseOauthCliConfig := oauth2cli.Config{OAuth2Config: oAuth2Config()}

	nonInteractiveTarget, _ := cmd.Flags().GetBool(nonInteractiveTargetHostFlag)
	nonInteractiveRemote, _ := cmd.Flags().GetBool(nonInteractiveRemoteHostFlag)

	if nonInteractiveTarget {
		baseOauthCliConfig.NonInteractive = true
		baseOauthCliConfig.NonInteractivePromptText = fmt.Sprintf("On a machine with access to a browser, use %s auth login --%s to retrieve a valid code:\n", common.RootCommandName, nonInteractiveRemoteHostFlag)
		baseOauthCliConfig.OAuth2Config.RedirectURL = "http://localhost:10000"
		eg.Go(authenticator.handleLogin(ctx, baseOauthCliConfig))
	} else {
		ready := make(chan string, 1)
		defer close(ready)
		port := findFreePort()

		baseOauthCliConfig.LocalServerReadyChan = ready
		baseOauthCliConfig.LocalServerBindAddress = strings.Split(fmt.Sprintf("127.0.0.1:%d", port), ",")
		baseOauthCliConfig.SuccessRedirectURL = "https://strmprivacy.io/auth-success"
		baseOauthCliConfig.FailureRedirectURL = "https://strmprivacy.io/auth-failure"
		baseOauthCliConfig.AuthCodeOptions = []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("prompt", "login")}

		eg.Go(startBrowserLoginFlow(ready, ctx))

		if nonInteractiveRemote {
			eg.Go(authenticator.handleCode(ctx, baseOauthCliConfig))
		} else {
			eg.Go(authenticator.handleLogin(ctx, baseOauthCliConfig))
		}
	}

	if err := eg.Wait(); err != nil {
		common.CliExit(errors.New(fmt.Sprintf("Login failed, please check the logs for details at %v", common.LogFileName())))
	}
}

func (authenticator *Authenticator) handleCode(ctx context.Context, cfg oauth2cli.Config) func() error {
	return func() error {
		codeAndConfig, err := oauth2cli.GetCodeAndConfig(ctx, cfg)
		common.CliExit(err)

		fmt.Println(fmt.Sprintf("\nUse the following text on the headless host to login:\n%v", *codeAndConfig))

		return nil
	}
}

func (authenticator *Authenticator) handleLogin(ctx context.Context, cfg oauth2cli.Config) func() error {
	return func() error {
		oAuthToken, err := oauth2cli.GetToken(ctx, cfg)
		if err != nil {
			rootCauseErr := getRootCause(err)

			switch rootCauseErr.(type) {
			case base64.CorruptInputError:
				common.CliExit(errors.New(fmt.Sprintf("\nInvalid base64 encoded input. Make sure that the input you provide is retrieved using %s auth login --%s", common.RootCommandName, nonInteractiveRemoteHostFlag)))
			case *json.SyntaxError:
				common.CliExit(errors.New(fmt.Sprintf("\nMalformed JSON input. Make sure that the input you provide is retrieved using %s auth login --%s", common.RootCommandName, nonInteractiveRemoteHostFlag)))
			case *oauth2.RetrieveError:
				retrieveErr := (rootCauseErr).(*oauth2.RetrieveError)
				common.CliExit(errors.New(fmt.Sprintf("\nUnable to exchange authorization code with token (HTTP Code: %v, Body: %v)", retrieveErr.Response.Status, string(retrieveErr.Body))))
			default:
				common.CliExit(fmt.Errorf("\n%w", rootCauseErr))
			}
		}

		authenticator.populateValues(oauthTokenToStoredToken(*oAuthToken))
		authenticator.storeLogin()

		fmt.Println(fmt.Sprintf("\nYou are now logged in as [%v].", authenticator.Email))

		return nil
	}
}

func getRootCause(err error) error {
	nextErr := errors.Unwrap(err)

	if nextErr != nil {
		return getRootCause(nextErr)
	} else {
		return err
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
