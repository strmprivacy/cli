package diagnostics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

type Metrics struct {
	K int            `json:"k"`
	L map[string]int `json:"L"`
	T float64        `json:"T"`
}

func evaluate(cmd *cobra.Command) {
	flags := cmd.Flags()
	dataFile := util.GetStringAndErr(flags, dataFileFlagName)
	qi := getFlags(cmd, quasiIdentifierFlagName)
	sa := getFlags(cmd, sensitiveAttributeFlagName)
	metrics := getFlags(cmd, metricsFlagName)
	b, fw := buildForm(dataFile, qi, sa, metrics)

	request, err := createRequest(b, fw)

	// TODO change client
	client := &http.Client{}
	response, err := client.Do(request)
	common.CliExit(err)

	err = fw.Close()
	common.CliExit(err)

	body, err := io.ReadAll(response.Body)
	common.CliExit(err)

	output := Metrics{}
	err = json.Unmarshal(body, &output)
	common.CliExit(err)
	printer.Print(output)
}

func createRequest(b *bytes.Buffer, formWriter *multipart.Writer) (*http.Request, error) {
	// TODO change url
	request, err := http.NewRequest("POST", "http://localhost:8080/upload", b)
	common.CliExit(err)
	contentType := fmt.Sprintf("multipart/form-data;boundary=%v", formWriter.Boundary())
	request.Header.Add("Content-Type", contentType)
	return request, err
}

func buildForm(dataFile string, qi []string, sa []string, metrics []string) (*bytes.Buffer, *multipart.Writer) {
	b := &bytes.Buffer{}
	formWriter := multipart.NewWriter(b)

	part, err := formWriter.CreateFormFile("upload", filepath.Base(dataFile))
	common.CliExit(err)

	data, err := os.ReadFile(dataFile)
	common.CliExit(err)

	_, err = part.Write(data)
	common.CliExit(err)

	for ix, value := range qi {
		key := "qi[" + strconv.Itoa(ix) + "]"
		err = formWriter.WriteField(key, value)
		common.CliExit(err)
	}
	for ix, value := range sa {
		key := "sa[" + strconv.Itoa(ix) + "]"
		err = formWriter.WriteField(key, value)
		common.CliExit(err)
	}
	for ix, value := range metrics {
		key := "metrics[" + strconv.Itoa(ix) + "]"
		err = formWriter.WriteField(key, value)
		common.CliExit(err)
	}
	return b, formWriter
}

func getFlags(cmd *cobra.Command, flag string) []string {
	columns := util.GetStringAndErr(cmd.Flags(), flag)
	return strings.Split(columns, ",")
}
