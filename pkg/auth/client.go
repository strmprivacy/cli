package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"streammachine.io/strm/pkg/common"
	"time"
)

var Client Auth
var ConfigPath string

func SetupClient(apiAuthUrl string) {
	Client = Auth{Uri: apiAuthUrl}
}

// Auth is the entity that interacts with authorization endpoints.
// Used both for user logins and events
type Auth struct {
	Uri       string // don't add the path, only http[s]://<some.host>
	token     *token // filled in by AuthenticateEvent and refresh
	TokenPath string
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// token defines sts json response. Dont' change.
// For event tokens billingId and Email are not filled in.
type token struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int    `json:"expiresAt"`
	BillingId    string `json:"billingId"`
	Email        string `json:"email"`
}

type EmptyTokenError struct {
	tokenFilePath string
}

func (e *EmptyTokenError) Error() string {
	return fmt.Sprintf("Token from file %v is empty", e.tokenFilePath)
}

func (authorizer *Auth) GetToken(quiet bool) (string, string) {
	if int64(authorizer.token.ExpiresAt)-30 < time.Now().Unix() {
		if !quiet {
			println("Refreshing STS token")
		}
		authorizer.refresh()
	}
	return authorizer.token.IdToken, authorizer.token.BillingId
}

func (authorizer *Auth) refresh() {
	b, err := json.Marshal(authorizer.token)
	common.CliExit(err)
	resp, err := http.Post(authorizer.Uri+"/refresh", "application/json; charset=UTF-8", bytes.NewReader(b))
	common.CliExit(err)
	defer resp.Body.Close()
	authorizer.handleAuthResponse(resp)
}

func (authorizer *Auth) AuthenticateLogin(email, password *string) {
	type authLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	postBody, err := json.Marshal(authLogin{
		Email:    *email,
		Password: *password,
	})
	common.CliExit(err)

	resp, err := http.Post(authorizer.Uri+"/auth", "application/json; charset=UTF-8", bytes.NewReader(postBody))
	common.CliExit(err)
	if resp.StatusCode != 200 {
		errorBody, _ := io.ReadAll(resp.Body)
		var errorJson = ErrorResponse{Code: 0, Message: ""}
		err2 := json.Unmarshal(errorBody, &errorJson)
		common.CliExit(err2)
		common.CliExit("Error authenticating: " + errorJson.Message)
	}
	defer resp.Body.Close()
	authorizer.handleAuthResponse(resp)
}

func (authorizer *Auth) AuthenticateEvent(billingId, clientId, secret string) {
	// authJson json format for posting to sts.  Don't change
	type authJson struct {
		BillingId    string `json:"billingId"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}

	postBody, err := json.Marshal(authJson{
		BillingId: billingId, ClientId: clientId, ClientSecret: secret,
	})
	common.CliExit(err)

	resp, err := http.Post(authorizer.Uri+"/auth", "application/json; charset=UTF-8", bytes.NewReader(postBody))
	common.CliExit(err)
	defer resp.Body.Close()
	authorizer.handleAuthResponse(resp)
}

func (authorizer *Auth) handleAuthResponse(resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	common.CliExit(err)
	err = json.Unmarshal(body, &authorizer.token)
	if &authorizer.token.IdToken == nil {
		common.CliExit("Cannot get ID token from auth response")
	}
	common.CliExit(err)
}

func (authorizer *Auth) StoreLogin() string {
	filename := authorizer.getSaveFilename()
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	common.CliExit(err)
	b, err := json.Marshal(authorizer.token)
	common.CliExit(err)
	err = ioutil.WriteFile(filename, b, 0644)
	common.CliExit(err)
	return filename

}

func (authorizer *Auth) getSaveFilename() string {
	if TokenFile == "" {
		u, err := url.Parse(authorizer.Uri)
		common.CliExit(err)
		filename := fmt.Sprintf("strm-creds-%s.json", u.Hostname())
		return path.Join(ConfigPath, filename)
	} else {
		return TokenFile
	}
}

func (authorizer *Auth) LoadLogin() error {
	filename := authorizer.getSaveFilename()
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	} else if len(b) == 0 {
		return &EmptyTokenError{}
	} else {
		err = json.Unmarshal(b, &authorizer.token)
		common.CliExit(err)

		return nil
	}
}

func (authorizer *Auth) printToken() {
	if authorizer.token != nil {
		fmt.Println(authorizer.token.IdToken)
		// these go to stderr, so the token is easy to capture in a script
		println("Expires at:", time.Unix(int64(authorizer.token.ExpiresAt), 0).String())
		println("Billing id:", authorizer.token.BillingId)
	} else {
		common.MissingIdTokenError()
	}

}
