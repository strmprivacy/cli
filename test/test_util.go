package test

import (
	"context"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"streammachine.io/strm/pkg/util"
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

	_ = os.Setenv("STRM_EVENTS_AUTH_URL", "https://auth.dev.strm.services")
	_ = os.Setenv("STRM_EVENTS_API_URL", "https://in.dev.strm.services/event")
	_ = os.Setenv("STRM_API_AUTH_URL", "https://accounts.dev.streammachine.io")
	_ = os.Setenv("STRM_API_HOST", "apis.dev.streammachine.io:443")
	_ = os.Setenv("STRM_HEADLESS", "true")

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
	performCliLogin(t, tokenFileName)

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

func performCliLogin(t *testing.T, tokenFileName string) {
	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		out := ExecuteCliAndGetOutput(t, tokenFileName, "auth", "login")
		assert.Matches(t, out, ".*https://accounts\\.dev\\.streammachine\\.io/auth/realms/users/protocol/openid-connect/auth.*")
		assert.Matches(t, out, ".*You are now logged in as \\[clitest-dev@streammachine\\.io\\]\\.")

		return nil
	})
	eg.Go(loginInBrowser)

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func loginInBrowser() error {
	chrome := launcher.New().Headless(true).MustLaunch()
	connection := rod.New().ControlURL(chrome).MustConnect()
	page := connection.MustPage("http://localhost:10000")

	page.MustElement("#username").MustInput(testConfig().email)
	page.MustElement("#password").MustInput(testConfig().password)
	page.MustElement("button[name=login]").MustClick()
	page.MustWaitLoad()

	return nil
}
func TryLoad(m proto.Message, s string) (proto.Message, error) {
	err := protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal([]byte(s), m)
	return m, err
}

func assertProtoEquals(t *testing.T, actual proto.Message, expected proto.Message) {
	if !proto.Equal(actual, expected) {
		printer := util.ProtoMessageJsonPrettyPrinter{}
		fmt.Println("Assertion failure: different proto messages")
		fmt.Println("expected:")
		printer.Print(expected)
		fmt.Println("actual:")
		printer.Print(actual)
		t.Fail()
	}
}

func ExecuteAndVerify(t *testing.T, expected proto.Message, args ...string) {
	/*
		we need a proto message of the same type as the expected, and since
		golang does not (yet) have generics, we can do it this way.
		the contents of the cloned message are NOT used, only its address for the TryLoad
		call
	*/
	outputMessage := proto.Clone(expected)
	output := ExecuteCliAndGetOutput(t, "", args...)
	if strings.HasPrefix(output, "Error") {
		t.Error(output)
	}
	out, err := TryLoad(outputMessage, output)
	if err != nil {
		fmt.Println("Can't execute", args)
		t.Fail()
	}
	assertProtoEquals(t, out, expected)
}
