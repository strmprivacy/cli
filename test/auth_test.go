package test

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestAuthAccessTokenOutputsAnErrorWhenNotLoggedIn(t *testing.T) {
	_ = newConfigDir()
	tokenFileName := CreateNonExistingTokenFileName()

	out := ExecuteCliAndGetOutput(t, tokenFileName, "auth", "print-access-token")

	assert.Equal(t, out, "Error: No login information found. Use: `dstrm auth login` first.\n")
}
