package datasubject

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_subjects/v1"
	"google.golang.org/grpc"
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