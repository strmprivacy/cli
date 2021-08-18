package test

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

type TestConfig struct {
	billingId         string
	email             string
	password          string
	s3UserName        string
	s3AccessKeyId     string
	s3SecretAccessKey string
}

var _testConfig TestConfig

func testConfig() *TestConfig {
	if (TestConfig{}) == _testConfig {
		_ = godotenv.Load()

		_testConfig = TestConfig{
			billingId:         os.Getenv("STRM_TEST_USER_BILLING_ID"),
			email:             os.Getenv("STRM_TEST_USER_EMAIL"),
			password:          os.Getenv("STRM_TEST_USER_PASSWORD"),
			s3UserName:        os.Getenv("STRM_TEST_S3_USER_NAME"),
			s3AccessKeyId:     os.Getenv("STRM_TEST_S3_ACCESS_KEY_ID"),
			s3SecretAccessKey: os.Getenv("STRM_TEST_S3_SECRET_ACCESS_KEY"),
		}

		if _testConfig == (TestConfig{}) {
			println("Error: Configuration parameters not present!")
		}
	}

	return &_testConfig
}

func newConfigDir() string {
	var err error
	configDir, err = ioutil.TempDir("", "test")
	if err != nil {
		println(fmt.Sprintf("error: %v", err))
	}
	defaultTokenFileName = ""

	_ = os.Setenv("TZ", "UTC")
	_ = os.Setenv("STRM_CONFIG_PATH", configDir)

	_ = os.Setenv("STRM_EVENT_AUTH_HOST", "https://auth.dev.strm.services")
	_ = os.Setenv("STRM_EVENTS_GATEWAY", "https://in.dev.strm.services/event")
	_ = os.Setenv("STRM_API_AUTH_URL", "https://api.dev.streammachine.io/v1")
	_ = os.Setenv("STRM_API_HOST", "apis.dev.streammachine.io:443")

	return configDir
}

type TokenFile struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int    `json:"expiresAt"`
	BillingId    string `json:"billingId"`
	Email        string `json:"email"`
}

var configDir string
var defaultTokenFileName string

func ExecuteCliAndGetOutput(t *testing.T, tokenFile string, cmd ...string) string {
	command := executeCli(t, tokenFile, cmd...)
	out, _ := command.CombinedOutput()
	s := string(out)

	s = replaceSecretsWithPropertyNames(s)
	return s
}

func executeCli(t *testing.T, tokenFile string, cmd ...string) *exec.Cmd {
	if len(tokenFile) == 0 {
		if len(defaultTokenFileName) == 0 {
			defaultTokenFile, _ := ioutil.TempFile(configDir, "strm_*.json")
			defaultTokenFileName = defaultTokenFile.Name()

			initializeStrmEntities(t, defaultTokenFileName)
		}
		tokenFile = defaultTokenFileName
	}

	path, _ := os.Getwd()
	rootPath := path + "/../"

	cmd = append(cmd, "--token-file="+tokenFile)
	cmd = append(cmd, "--output=json-raw")

	command := exec.Command(rootPath+"dist/dstrm", cmd...)
	command.Dir = rootPath
	return command
}

func initializeStrmEntities(t *testing.T, tokenFileName string) {
	// Create a default login for this test run
	loginOut := ExecuteCliAndGetOutput(t, tokenFileName, "auth", "login", testConfig().email, "--password="+testConfig().password)
	assert.Equal(t, loginOut, "Billing id: testBillingId\nSaved login to: "+defaultTokenFileName+"\n")

	// Remove all resources (ugly implementation until we have plain text output in the CLI)
	nameMatcher := regexp.MustCompile(`"name":"([^"]+)"`)

	sinksOut := ExecuteCliAndGetOutput(t, tokenFileName, "list", "sinks")
	allSinkNames := nameMatcher.FindAllStringSubmatch(sinksOut, -1)
	for i := 0; i < len(allSinkNames); i++ {
		out := ExecuteCliAndGetOutput(t, tokenFileName, "delete", "sink", allSinkNames[i][1], "--recursive")
		if !strings.HasPrefix(out, "{") {
			t.Error("delete sink " + allSinkNames[i][1] + " failed with error: " + out)
		}
	}

	streamsOut := ExecuteCliAndGetOutput(t, tokenFileName, "list", "streams")
	allStreamNames := nameMatcher.FindAllStringSubmatch(streamsOut, -1)
	for i := 0; i < len(allStreamNames); i++ {
		out := ExecuteCliAndGetOutput(t, tokenFileName, "delete", "stream", allStreamNames[i][1], "--recursive")
		if !strings.HasPrefix(out, "{") {
			t.Error("delete stream " + allSinkNames[i][1] + " failed with error: " + out)
		}
	}
}

func replaceSecretsWithPropertyNames(out string) string {
	clientIdReplacer := regexp.MustCompile(`clientId":"([^"]+)"`)
	out = clientIdReplacer.ReplaceAllString(out, `clientId":"clientId"`)

	clientSecretReplacer := regexp.MustCompile(`clientSecret":"([^"]+)"`)
	out = clientSecretReplacer.ReplaceAllString(out, `clientSecret":"clientSecret"`)

	s3UserNameReplacer := regexp.MustCompile(`UserName\\":\\"([^"]+)\\"`)
	out = s3UserNameReplacer.ReplaceAllString(out, `UserName\":\"UserName\"`)

	s3AccessKeyIdReplacer := regexp.MustCompile(`AccessKeyId\\":\\"([^"]+)\\"`)
	out = s3AccessKeyIdReplacer.ReplaceAllString(out, `AccessKeyId\":\"AccessKeyId\"`)

	s3SecretAccessKeyReplacer := regexp.MustCompile(`SecretAccessKey\\":\\"([^"]+)\\"`)
	out = s3SecretAccessKeyReplacer.ReplaceAllString(out, `SecretAccessKey\":\"SecretAccessKey\"`)

	testBillingIdReplacer := regexp.MustCompile("(" + testConfig().billingId + ")")
	out = testBillingIdReplacer.ReplaceAllString(out, "testBillingId")

	return out
}

func CreateNonExistingTokenFileName() string {
	tokenDir, _ := ioutil.TempDir("", "strm_test")
	tokenFileName := tokenDir + "/nonexisting.json"
	return tokenFileName
}
