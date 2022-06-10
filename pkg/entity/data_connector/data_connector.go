package data_connector

import (
	"context"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_connectors/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var Client data_connectors.DataConnectorsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	Client = data_connectors.NewDataConnectorsServiceClient(clientConnection)
}

func list(recursive bool) {
	req := &data_connectors.ListDataConnectorsRequest{
		BillingId:          auth.Auth.BillingId(),
		ProjectId:          common.ProjectId,
		Recursive:          recursive,
		IncludeCredentials: false,
	}
	response, err := Client.ListDataConnectors(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(name *string, recursive bool) {
	req := &data_connectors.GetDataConnectorRequest{
		Ref:                ref(name),
		Recursive:          recursive,
		IncludeCredentials: false,
	}
	response, err := Client.GetDataConnector(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(name *string, recursive bool) {
	req := &data_connectors.DeleteDataConnectorRequest{
		Ref:       ref(name),
		Recursive: recursive,
	}
	response, err := Client.DeleteDataConnector(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func create(dataConnector *entities.DataConnector) {
	req := &data_connectors.CreateDataConnectorRequest{
		DataConnector: dataConnector,
	}
	response, err := Client.CreateDataConnector(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func ref(name *string) *entities.DataConnectorRef {
	return &entities.DataConnectorRef{
		BillingId: auth.Auth.BillingId(),
		ProjectId: common.ProjectId,
		Name: *name,
	}
}

func readCredentialsFile(flags *pflag.FlagSet) string {
	fn := util.GetStringAndErr(flags, credentialsFileFlag)
	buf, err := ioutil.ReadFile(fn)
	common.CliExit(err)
	return string(buf)
}
