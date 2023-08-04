package data_subject

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/data_subjects/v1"
	"google.golang.org/grpc"
	"os"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
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

func list(cmd *cobra.Command) {
	var err error
	flags := cmd.Flags()
	req := &data_subjects.ListDataSubjectsRequest{}
	req.PageSize = util.GetInt32AndErr(flags, pageSizeFlag)
	req.PageToken = util.GetStringAndErr(flags, pageTokenFlag)
	response, err := client.ListDataSubjects(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(args []string, cmd *cobra.Command) {
	if len(args) == 0 {

		fmt.Fprintln(os.Stderr, `Call with at least one data-subject argument

You can retrieve valid arguments via 'list data-subjects' `)
		os.Exit(0)
	}
	req := &data_subjects.DeleteDataSubjectsRequest{DataSubjects: args}
	response, err := client.DeleteDataSubjects(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
