package sim

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strmprivacy/strm/pkg/common"
)

const (
	IntervalFlag      = "interval"
	EventsApiUrlFlag  = "events-api-url"
	SessionRangeFlag  = "session-range"
	SessionPrefixFlag = "session-prefix"
	ConsentLevelsFlag = "consent-levels"
	QuietFlag         = "quiet"
	SchemaFlag        = "schema"
)

type StrmPrivacyEvent interface {
	Serialize(w io.Writer) error
}

type Sender interface {
	Send(event StrmPrivacyEvent, token string)
}

type ModernSender struct {
	Gateway, Schema string
	Client          http.Client
}

func (s ModernSender) Send(event StrmPrivacyEvent, token string) {
	b := &bytes.Buffer{}
	err := event.Serialize(b)
	common.CliExit(err)

	req, err := http.NewRequest("POST", s.Gateway, b)
	common.CliExit(err)

	req.Header.Set("Strm-Schema-Ref", s.Schema)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)

	if err != nil {
		common.CliExit(errors.New(fmt.Sprintf("An error occurred while simulating events.\nError: %v", err)))
	} else if resp.StatusCode != 204 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		common.CliExit(errors.New(fmt.Sprintf("The simulated event sent resulted in an error (expected status is not 204).\nError: %v, Response Body: %s", err, string(body))))
	}
}
