package sim

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strmprivacy/strm/pkg/common"
)

type LegacySender struct {
	Gateway, Schema string
	Client          http.Client
}

func (s LegacySender) Send(event StrmPrivacyEvent, token string) {
	b := &bytes.Buffer{}
	err := event.Serialize(b)
	common.CliExit(err)
	req, err := http.NewRequest("POST", s.Gateway, b)
	common.CliExit(err)
	req.Header.Set("Strm-Schema-Id", s.Schema)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.Client.Do(req)
	if err != nil || resp.StatusCode != 204 {
		if resp != nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			fmt.Printf("%v %s\n", err, string(body))
		}
	}
}
