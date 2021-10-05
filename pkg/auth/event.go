package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"streammachine.io/strm/pkg/common"
	"time"
)

var token *eventToken

type eventToken struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int    `json:"expiresAt"`
}

type authJson struct {
	BillingId    string `json:"billingId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func GetEventToken(billingId string, clientId string, secret string) string {
	if token == nil {
		authenticate(billingId, clientId, secret)
	} else if int64((*token).ExpiresAt)-30 < time.Now().Unix() {
		refresh()
	}

	return (*token).IdToken
}

func authenticate(billingId string, clientId string, secret string) {
	postBody, err := json.Marshal(authJson{
		BillingId: billingId, ClientId: clientId, ClientSecret: secret,
	})
	common.CliExit(err)

	resp, err := http.Post(common.EventAuthHost+"/auth", "application/json; charset=UTF-8", bytes.NewReader(postBody))
	common.CliExit(err)
	defer resp.Body.Close()

	handleAuthResponse(resp)
}

func refresh() {
	b, err := json.Marshal(token)
	common.CliExit(err)
	resp, err := http.Post(common.EventAuthHost+"/refresh", "application/json; charset=UTF-8", bytes.NewReader(b))
	common.CliExit(err)
	defer resp.Body.Close()
	handleAuthResponse(resp)
}

func handleAuthResponse(resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	common.CliExit(err)

	eventToken := eventToken{}

	err = json.Unmarshal(body, &eventToken)
	if &eventToken.IdToken == nil {
		common.CliExit("Cannot get ID token from auth response")
	}
	common.CliExit(err)

	token = &eventToken
}
