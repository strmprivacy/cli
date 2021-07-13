package test

import (
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestAuthAccessTokenOutputsTheAccessTokenWhenLoggedIn(t *testing.T) {
	_ = newConfigDir()
	out := ExecuteCliAndGetOutput(t, "test/simple-token.json", "auth", "access-token")

	assert.Equal(t, out, "id.token.test\nExpires at: 2021-07-02 15:42:54 +0000 UTC\nBilling id: my.billing.id\n")
}

func TestAuthAccessTokenOutputsAnErrorWhenNotLoggedIn(t *testing.T) {
	_ = newConfigDir()
	tokenFileName := CreateNonExistingTokenFileName()

	out := ExecuteCliAndGetOutput(t, tokenFileName, "auth", "access-token")

	assert.Equal(t, out, "Error: No login information found. Use: `dstrm auth login` first.\n")
}

func TestAuthLoginWithAPasswordSpecifiedLogsTheUserIn(t *testing.T) {
	_ = newConfigDir()
	tokenFileName := CreateNonExistingTokenFileName()

	out := ExecuteCliAndGetOutput(t, tokenFileName, "auth", "login", testConfig().email, "--password="+testConfig().password)

	assert.Matches(t, out, "Billing id: testBillingId\nSaved login to: "+tokenFileName+"\n")
}

func TestAuthRefreshRefreshesTheToken(t *testing.T) {
	_ = newConfigDir()
	tokenFileName := CreateNonExistingTokenFileName()

	out := ExecuteCliAndGetOutput(t, tokenFileName, "auth", "login", testConfig().email, "--password="+testConfig().password)

	assert.Matches(t, out, "Billing id: testBillingId\nSaved login to: "+tokenFileName+"\n")
	var oldTokenBytes, oldTokenFileError = ioutil.ReadFile(tokenFileName)
	if oldTokenFileError != nil {
		t.Error(oldTokenFileError)
		return
	}
	var oldTokenFile TokenFile
	var oldTokenJsonError = json.Unmarshal(oldTokenBytes, &oldTokenFile)
	if oldTokenJsonError != nil {
		t.Error(oldTokenJsonError)
		return
	}

	if oldTokenFile.IdToken == "" {
		t.Error("idToken is empty after login!")
		return
	}

	time.Sleep(1 * time.Second)

	out = ExecuteCliAndGetOutput(t, tokenFileName, "auth", "refresh")
	assert.Equal(t, out, "")

	var newTokenBytes, newTokenFileError = ioutil.ReadFile(tokenFileName)
	if newTokenFileError != nil {
		t.Error(newTokenFileError)
		return
	}

	var newTokenFile TokenFile
	var newTokenJsonError = json.Unmarshal(newTokenBytes, &newTokenFile)
	if newTokenJsonError != nil {
		t.Error(newTokenJsonError)
		return
	}

	if newTokenFile.IdToken == "" {
		t.Error("idToken is empty after refresh!")
		return
	}

	if newTokenFile.IdToken == oldTokenFile.IdToken {
		t.Error("refreshed idToken is equal to old idToken!")
		return
	}
}
