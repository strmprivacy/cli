package auth

import (
	"fmt"
	"golang.org/x/oauth2"
	"streammachine.io/strm/pkg/common"
	"strings"
	"time"
)

var Auth Authenticator

func SetupAuth(authUrl string) {
	Auth = Authenticator{Uri: authUrl}
}

// Authenticator is the entity that interacts with authorization endpoints.
// Used both for user logins and events
type Authenticator struct {
	Uri         string
	storedToken storedToken
	tokenSource oauth2.TokenSource
	billingId *string
	Email     string
}

func (authorizer *Authenticator) BillingId() string {
	if authorizer.billingId == nil {
		common.MissingIdTokenError()
	}

	return *authorizer.billingId
}

func (authorizer *Authenticator) GetToken() *string {
	accessToken, _ := authorizer.accessTokenAndExpiry()

	return accessToken
}

func (authorizer *Authenticator) BillingIdAbsent() bool {
	return len(strings.TrimSpace(authorizer.BillingId())) == 0
}

func (authorizer *Authenticator) printAccessToken() {
	accessToken, _ := authorizer.accessTokenAndExpiry()
	if accessToken != nil {
		fmt.Println(*accessToken)
	} else {
		common.MissingIdTokenError()
	}
}

func (authorizer *Authenticator) accessTokenAndExpiry() (*string, *time.Time) {
	if authorizer.tokenSource == nil {
		return nil, nil
	} else {
		tokens, err := authorizer.tokenSource.Token()
		common.CliExit(err)
		if authorizer.storedToken.AccessToken != tokens.AccessToken {
			authorizer.populateValues(oauthTokenToStoredToken(*tokens))
			authorizer.StoreLogin()
		}
		return &tokens.AccessToken, &tokens.Expiry
	}
}


