package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/int128/oauth2cli"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/streammachineio/api-definitions-go/api/account/v1"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"net/http"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/entity"
	"strings"
)

var oAuth2Config = oauth2.Config{
	ClientID: "cli",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.dev.streammachine.io/auth/auth",
		TokenURL: "https://accounts.dev.streammachine.io/auth/token",
	},
	Scopes: []string{"offline_access", "email"},
}

func (authorizer *Authenticator) Login() {
	ready := make(chan string, 1)
	defer close(ready)

	cfg := oauth2cli.Config{
		OAuth2Config:           oAuth2Config,
		LocalServerReadyChan:   ready,
		LocalServerBindAddress: strings.Split("127.0.0.1:10000", ","),
		LocalServerSuccessHTML: constants.AuthSuccessHTML,
	}

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(startBrowserLoginFlow(ready, ctx))
	eg.Go(authorizer.handleLogin(ctx, cfg))

	if err := eg.Wait(); err != nil {
		common.CliExit(fmt.Sprintf("Login failed, please check the logs for details at %v", common.LogFileName()))
	}
}

func (authorizer *Authenticator) handleLogin(ctx context.Context, cfg oauth2cli.Config) func() error {
	return func() error {
		oAuthToken, err := oauth2cli.GetToken(ctx, cfg)
		common.CliExit(err)

		authorizer.populateValues(oauthTokenToStoredToken(*oAuthToken))
		authorizer.StoreLogin()

		fmt.Println(fmt.Sprintf("\nYou are now logged in as [%v].", authorizer.Email))

		return nil
	}
}

func oauthTokenToStoredToken(t oauth2.Token) storedToken {
	return storedToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    t.Expiry.Unix(),
		BillingId:    getLegacyBillingId(t.AccessToken),
		Email:        getEmailFromClaims(t),
	}
}

func startBrowserLoginFlow(ready chan string, ctx context.Context) func() error {
	return func() error {
		select {
		case localhostCallbackUrl := <-ready:
			err := browser.OpenURL(localhostCallbackUrl)
			if err != nil {
				browserOpenError := fmt.Sprintf("Unable to open browser for authentication: %s", err)
				log.Error(browserOpenError)
				common.CliExit(browserOpenError)
			} else {
				fmt.Println("Follow the login flow in your browser, which is opened automatically. If not, open the following URL to complete the login:")
				fmt.Println(fmt.Sprintf("\n    %v", authorizationCodeFlowUrl(localhostCallbackUrl)))
			}
			return nil
		case <-ctx.Done():
			contextError := fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
			log.Error(contextError)
			common.CliExit(contextError)
			return nil
		}
	}
}

func getLegacyBillingId(accessToken string) string {
	clientConnection, ctx := entity.SetupGrpc(common.ApiHost, &accessToken)
	accountClient := account.NewAccountServiceClient(clientConnection)

	response, err := accountClient.GetLegacyBillingId(ctx, &account.GetLegacyBillingIdRequest{})
	common.CliExit(err)
	return response.BillingId
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
		common.CliExit("Could not retrieve authorization code flow URL. Please retry or contact Stream Machine support if the problem persists.")
	}

	return locationResponse.Request.URL.String()
}
