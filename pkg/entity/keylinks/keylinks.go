package keylinks

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/data_subjects/v1"
	"google.golang.org/grpc"
	"os"
	"strmprivacy/strm/pkg/common"
)

const (
	pageSizeFlag  = "page-size"
	pageTokenFlag = "page-token"
)

var client data_subjects.DataSubjectsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_subjects.NewDataSubjectsServiceClient(clientConnection)
}

func list(args []string, cmd *cobra.Command) {
	if len(args) == 0 {

		fmt.Fprintln(os.Stderr, `Call with at least one data-subject argument

You can retrieve valid arguments via 'list data-subjects' `)
		os.Exit(0)
	}
	var err error
	req := &data_subjects.ListDataSubjectKeylinksRequest{}
	req.DataSubjectId = args
	response, err := client.ListDataSubjectKeylinks(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}
