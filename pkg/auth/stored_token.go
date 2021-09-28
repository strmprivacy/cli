package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"time"
)

const StrmCredsFilePrefix = "strm-creds"

// The auth information that is persisted on disk
type storedToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	BillingId    string `json:"billingId"`
	Email        string `json:"email"`
}

type EmptyTokenError struct {
	tokenFilePath string
}

func (e *EmptyTokenError) Error() string {
	return fmt.Sprintf("Token from file %v is empty", e.tokenFilePath)
}

func (authorizer *Authenticator) LoadLogin() error {
	filename := authorizer.getSaveFilename()
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	} else if len(b) == 0 {
		return &EmptyTokenError{}
	} else {
		storedToken := unmarshalStoredToken(err, b)
		authorizer.populateValues(storedToken)

		return nil
	}
}

func (authorizer *Authenticator) populateValues(storedToken storedToken) {
	authorizer.storedToken = storedToken
	authorizer.tokenSource = createTokenSource(authorizer.storedToken)
	authorizer.billingId = &authorizer.storedToken.BillingId
	authorizer.Email = authorizer.storedToken.Email
}

func (authorizer *Authenticator) StoreLogin() string {
	filename := authorizer.getSaveFilename()
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	common.CliExit(err)
	b, err := json.Marshal(authorizer.storedToken)
	common.CliExit(err)
	err = ioutil.WriteFile(filename, b, 0644)
	common.CliExit(err)
	return filename
}

func (authorizer *Authenticator) getSaveFilename() string {
	if TokenFile == "" {
		u, err := url.Parse(authorizer.Uri)
		common.CliExit(err)
		filename := fmt.Sprintf("%v-%s.json", StrmCredsFilePrefix, u.Hostname())
		return path.Join(constants.ConfigPath, filename)
	} else {
		return TokenFile
	}
}

func unmarshalStoredToken(err error, b []byte) storedToken {
	storedToken := &storedToken{}
	err = json.Unmarshal(b, &storedToken)
	common.CliExit(err)
	return *storedToken
}

func createTokenSource(storedToken storedToken) oauth2.TokenSource {
	oauth2Token := &oauth2.Token{
		AccessToken:  storedToken.AccessToken,
		TokenType:    "bearer",
		RefreshToken: storedToken.RefreshToken,
		Expiry:       time.Unix(storedToken.ExpiresAt, 0),
	}

	ctx := context.Background()
	tokenSource := oAuth2Config.TokenSource(ctx, oauth2Token)
	return tokenSource
}
