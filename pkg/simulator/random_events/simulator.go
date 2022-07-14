package random_events

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strmprivacy/strm/pkg/common"
)

type Simulator struct {
	Gateway string
	Schema  string
	Client  http.Client
}

func (s Simulator) Send(event StrmPrivacyEvent) {
	b := &bytes.Buffer{}
	err := event.Serialize(b)
	common.CliExit(err)

	req, err := http.NewRequest("POST", s.Gateway, b)
	common.CliExit(err)

	req.Header.Set("Strm-Schema-Ref", s.Schema)

	resp, err := s.Client.Do(req)

	if err != nil {
		common.CliExit(errors.New(fmt.Sprintf("An error occurred while simulating events.\nError: %v", err)))
	} else if resp.StatusCode != 204 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		common.CliExit(errors.New(fmt.Sprintf("The simulated event sent resulted in an error (expected status is not 204).\nError: %v, Response Body: %s", err, string(body))))
	}
}
