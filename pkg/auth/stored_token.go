package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strmprivacy/strm/pkg/common"
	"time"
)

const StrmCredsFilePrefix = "strm-creds"

// The auth information that is persisted on disk
type storedToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	Email        string `json:"email"`
}

type EmptyTokenError struct {
	tokenFilePath string
}

func (e *EmptyTokenError) Error() string {
	return fmt.Sprintf("Token from file %v is empty", e.tokenFilePath)
}

func (authenticator *Authenticator) LoadLogin() error {
	filename := authenticator.getSaveFilename()
	b, err := os.ReadFile(filename)

	if errors.Is(err, os.ErrNotExist) {
		return common.UnauthenticatedError()
	} else if len(b) == 0 {
		return &EmptyTokenError{}
	} else if err != nil {
		log.Errorln(fmt.Sprintf("Unexpected exception occurred while trying to read file %v", filename), err)
		return err
	} else {
		storedToken := unmarshalStoredToken(err, b)
		authenticator.populateValues(storedToken)

		return nil
	}
}

func (authenticator *Authenticator) populateValues(storedToken storedToken) {
	authenticator.storedToken = storedToken
	authenticator.tokenSource = createTokenSource(authenticator.storedToken)
	authenticator.Email = authenticator.storedToken.Email
}

func (authenticator *Authenticator) storeLogin() string {
	filename := authenticator.getSaveFilename()
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	common.CliExit(err)
	b, err := json.Marshal(authenticator.storedToken)
	common.CliExit(err)
	err = os.WriteFile(filename, b, 0644)
	common.CliExit(err)
	return filename
}

func (authenticator *Authenticator) getSaveFilename() string {
	if TokenFile == "" {
		u, err := url.Parse(common.ApiAuthHost)
		common.CliExit(err)
		filename := fmt.Sprintf("%v-%s.json", StrmCredsFilePrefix, u.Hostname())
		return path.Join(common.ConfigPath(), filename)
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
	oAuth2Config := oAuth2Config()
	tokenSource := oAuth2Config.TokenSource(ctx, oauth2Token)
	return tokenSource
}
