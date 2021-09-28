package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"streammachine.io/strm/pkg/common"
)

func (authorizer *Authenticator) AuthenticateEvent(billingId, clientId, secret string) {
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
	//authorizer.handleAuthResponse(resp)
}
