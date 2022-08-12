package member

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/organizations/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/projects/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
)

const (
	projectFlag      = "project"
	organizationFlag = "organization"
)

var projectsClient projects.ProjectsServiceClient
var organizationsClient organizations.OrganizationsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	projectsClient = projects.NewProjectsServiceClient(clientConnection)
	organizationsClient = organizations.NewOrganizationsServiceClient(clientConnection)
}

func list(cmd *cobra.Command) []*entities.User {
	flags := cmd.Flags()
	organizationBool, err := flags.GetBool(organizationFlag)
	common.CliExit(err)
	projectBool, err := flags.GetBool(projectFlag)
	common.CliExit(err)
	if organizationBool == projectBool {
		common.CliExit(errors.New("strm list members requires organizations or projects flag, not neither nor both"))
	}
	if projectBool {
		return listProjectMembers()
	} else {
		return listOrganizationMembers()
	}
}

func listProjectMembers() []*entities.User {
	projectId := project.GetProjectId(common.GetActiveProject())
	req := &projects.ListProjectMembersRequest{
		ProjectId: projectId,
	}
	response, err := projectsClient.ListProjectMembers(apiContext, req)
	common.CliExit(err)
	return response.ProjectMembers
}

func listOrganizationMembers() []*entities.User {
	req := &organizations.ListOrganizationMembersRequest{}
	response, err := organizationsClient.ListOrganizationMembers(apiContext, req)
	common.CliExit(err)
	return response.OrganizationMembers
}
