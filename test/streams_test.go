package test

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestStreams(t *testing.T) {
	_ = newConfigDir()

	out := ExecuteCliAndGetOutput(t, "", "list", "streams")
	assert.Equal(t, out, `{}
`)
	out = ExecuteCliAndGetOutput(t, "", "create", "stream", "clitest")
	assert.Equal(t, out, `{"stream":{"ref":{"billingId":"testBillingId","name":"clitest"},"enabled":true,"limits":{"eventRate":"10000","eventCount":"10000000"},"credentials":[{"clientId":"clientId","clientSecret":"clientSecret"}],"maskedFields":{}}}
`)
	out = ExecuteCliAndGetOutput(t, "", "create", "stream", "clitest-with-tags", "--tags=foo,bar,baz")
	assert.Equal(t, out, `{"stream":{"ref":{"billingId":"testBillingId","name":"clitest-with-tags"},"enabled":true,"limits":{"eventRate":"10000","eventCount":"10000000"},"tags":["foo","bar","baz"],"credentials":[{"clientId":"clientId","clientSecret":"clientSecret"}],"maskedFields":{}}}
`)
	out = ExecuteCliAndGetOutput(t, "", "create", "stream", "--derived-from=clitest-with-tags", "--levels=2")
	assert.Equal(t, out, `{"stream":{"ref":{"billingId":"testBillingId","name":"clitest-with-tags-2"},"consentLevels":[2],"consentLevelType":"CUMULATIVE","enabled":true,"linkedStream":"clitest-with-tags","credentials":[{"clientId":"clientId","clientSecret":"clientSecret"}],"maskedFields":{}}}
`)
	out = ExecuteCliAndGetOutput(t, "", "get", "stream", "clitest-with-tags")
	assert.Equal(t, out, `{"streamTree":{"stream":{"ref":{"billingId":"testBillingId","name":"clitest-with-tags"},"enabled":true,"limits":{"eventRate":"10000","eventCount":"10000000"},"tags":["foo","bar","baz"],"credentials":[{"clientId":"clientId"}],"maskedFields":{}}}}
`)
	out = ExecuteCliAndGetOutput(t, "", "delete", "stream", "clitest-with-tags", "--recursive")
	assert.Equal(t, out, `{"streamTree":{"stream":{"ref":{"billingId":"testBillingId","name":"clitest-with-tags"},"enabled":true,"limits":{"eventRate":"10000","eventCount":"10000000"},"tags":["foo","bar","baz"],"credentials":[{"clientId":"clientId"}],"maskedFields":{}},"keyStream":{"ref":{"billingId":"testBillingId","name":"clitest-with-tags"}},"derived":[{"ref":{"billingId":"testBillingId","name":"clitest-with-tags-2"},"consentLevels":[2],"consentLevelType":"CUMULATIVE","enabled":true,"limits":{},"linkedStream":"clitest-with-tags","credentials":[{"clientId":"clientId"}],"maskedFields":{}}]}}
`)
}
