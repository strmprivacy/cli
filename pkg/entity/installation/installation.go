package installation

import (
	"context"
	"fmt"
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
	req := &installations.GetInstallationRequest{Id: *id}
	response, err := client.GetInstallation(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	req := &installations.ListInstallationsRequest{}
	response, err := client.ListInstallations(apiContext, req)
	fmt.Println(response)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	installationIds := make([]string, 0, len(response.Installations))
	for _, s := range response.Installations {
		installationIds = append(installationIds, s.Id)
	}
	return installationIds, cobra.ShellCompDirectiveNoFileComp
}
