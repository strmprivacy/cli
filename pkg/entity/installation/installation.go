package installation

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/installations/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
)

var client installations.InstallationsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = installations.NewInstallationsServiceClient(clientConnection)
}

func get(id *string, _ *cobra.Command) {
	req := &installations.GetInstallationRequest{InstallationId: *id}
	response, err := client.GetInstallation(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func list() {
	req := &installations.ListInstallationsRequest{}
	response, err := client.ListInstallations(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	req := &installations.ListInstallationsRequest{}
	response, err := client.ListInstallations(apiContext, req)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	installationIds := make([]string, 0, len(response.Installations))
	for _, s := range response.Installations {
		installationIds = append(installationIds, s.Id)
	}
	return installationIds, cobra.ShellCompDirectiveNoFileComp
}
