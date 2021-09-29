package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"streammachine.io/strm/pkg/common"
)

type eventToken struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int    `json:"expiresAt"`
}

func GetEventToken(billingId string, clientId string, secret string) string {
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

	resp, err := http.Post(common.EventAuthHost+"/auth", "application/json; charset=UTF-8", bytes.NewReader(postBody))
	common.CliExit(err)
	defer resp.Body.Close()

	return handleAuthResponse(resp)
}

func handleAuthResponse(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	common.CliExit(err)
	eventToken := eventToken{}

	err = json.Unmarshal(body, &eventToken)
	if &eventToken.IdToken == nil {
		common.CliExit("Cannot get ID token from auth response")
	}
	common.CliExit(err)

	return eventToken.IdToken
}
