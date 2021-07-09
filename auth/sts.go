package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

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

func (authorizer *Auth) GetToken(quiet bool) (string, string) {
	if int64(authorizer.token.ExpiresAt)-30 < time.Now().Unix() {
		if !quiet {
			println("Refreshing sts token")
		}
		authorizer.refresh()
	}
	return authorizer.token.IdToken, authorizer.token.BillingId
}

func (authorizer *Auth) refresh() {
	b, err := json.Marshal(authorizer.token)
	cobra.CheckErr(err)
	resp, err := http.Post(authorizer.Uri+"/refresh", "application/json; charset=UTF-8", bytes.NewReader(b))
	cobra.CheckErr(err)
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
	cobra.CheckErr(err)

	resp, err := http.Post(authorizer.Uri+"/auth", "application/json; charset=UTF-8", bytes.NewReader(postBody))
	cobra.CheckErr(err)
	if resp.StatusCode != 200 {
		errorBody, _ := io.ReadAll(resp.Body)
		var errorJson = ErrorResponse{Code: 0, Message: ""}
		err2 := json.Unmarshal(errorBody, &errorJson)
		cobra.CheckErr(err2)
		cobra.CheckErr("Error authenticating: " + errorJson.Message)
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
	cobra.CheckErr(err)

	resp, err := http.Post(authorizer.Uri+"/auth", "application/json; charset=UTF-8", bytes.NewReader(postBody))
	cobra.CheckErr(err)
	defer resp.Body.Close()
	authorizer.handleAuthResponse(resp)
}

func (authorizer *Auth) handleAuthResponse(resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)
	err = json.Unmarshal(body, &authorizer.token)
	if &authorizer.token.IdToken == nil {
		cobra.CheckErr("Cannot get ID token from auth response")
	}
	cobra.CheckErr(err)
}

func (authorizer *Auth) StoreLogin() string {
	filename := authorizer.getSaveFilename()
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	cobra.CheckErr(err)
	b, err := json.Marshal(authorizer.token)
	cobra.CheckErr(err)
	err = ioutil.WriteFile(filename, b, 0644)
	cobra.CheckErr(err)
	return filename

}

func (authorizer *Auth) getSaveFilename() string {
	if TokenFile == "" {
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		u, err := url.Parse(authorizer.Uri)
		cobra.CheckErr(err)
		filename := fmt.Sprintf("strm-creds-%s.json", u.Hostname())
		return path.Join(home, ".config", "stream-machine", filename)
	} else {
		return TokenFile
	}
}

func (authorizer *Auth) LoadLogin() {
	filename := authorizer.getSaveFilename()
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		cobra.CheckErr(fmt.Sprintf("No login information found. Use: `%v auth login` first.", CommandName))
	}
	err = json.Unmarshal(b, &authorizer.token)
	cobra.CheckErr(err)
}

func (authorizer *Auth) printToken() {
	fmt.Println(authorizer.token.IdToken)
	// these go to stderr, so the token is easy to capture in a script
	println("Expires at", time.Unix(int64(authorizer.token.ExpiresAt), 0).String())
	println("Billing-id", authorizer.token.BillingId)
}
