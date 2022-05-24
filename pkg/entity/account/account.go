package account

import (
	"context"
	"github.com/strmprivacy/api-definitions-go/v2/api/account/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
)


var client account.AccountServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = account.NewAccountServiceClient(clientConnection)
}

func GetHandle() *account.GetAccountDetailsResponse {
	req := &account.GetAccountDetailsRequest{
		BillingId: auth.Auth.BillingId(),
	}
	details, err := client.GetAccountDetails(apiContext, req)
	common.CliExit(err)
	return details
}