package member

import (
	"context"
	"github.com/strmprivacy/api-definitions-go/v2/api/projects/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
)

var client projects.ProjectsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = projects.NewProjectsServiceClient(clientConnection)
}

func list(args []string) *projects.ListProjectMembersResponse {
	projectName := ""
	if len(args) > 0 {
		projectName = args[0]
	}
	projectId := project.GetProjectId(projectName)
	req := &projects.ListProjectMembersRequest{
		ProjectId: projectId,
	}
	response, err := client.ListProjectMembers(apiContext, req)
	common.CliExit(err)
	return response
}
