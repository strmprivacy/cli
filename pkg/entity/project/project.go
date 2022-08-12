package project

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/projects/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var client projects.ProjectsServiceClient
var apiContext context.Context

const (
	descriptionFlag   = "description"
	addMembersFlag    = "add-member"
	removeMembersFlag = "remove-member"
)

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = projects.NewProjectsServiceClient(clientConnection)
}

func ListProjects() *projects.ListProjectsResponse {
	if apiContext == nil {
		common.CliExit(errors.New(fmt.Sprint("No login information found. Use: `dstrm auth login` first.")))
	}
	req := &projects.ListProjectsRequest{}
	response, err := client.ListProjects(apiContext, req)
	common.CliExit(err)
	return response
}

func GetProject(projectName string) *entities.Project {
	for _, project := range ListProjects().Projects {
		if project.Name == projectName {
			return project
		}
	}
	return nil
}

func create(projectName *string, cmd *cobra.Command) *projects.CreateProjectResponse {
	flags := cmd.Flags()
	description := util.GetStringAndErr(flags, descriptionFlag)
	req := &projects.CreateProjectRequest{
		Project: &entities.Project{
			Name:        *projectName,
			Description: description,
		},
	}
	response, err := client.CreateProject(apiContext, req)
	common.CliExit(err)
	return response
}

func GetProjectId(projectName string) string {
	activeProject := ""
	if projectName == "" {
		activeProject = common.GetActiveProject()
	} else {
		activeProject = projectName
	}
	resolvedProject := GetProject(activeProject)
	if resolvedProject == nil {
		common.CliExit(errors.New(fmt.Sprintf("Project '%v' does not exist, or you do not have access "+
			"to it.", activeProject)))
	}
	return resolvedProject.Id
}

func manage(args []string, cmd *cobra.Command) {
	flags := cmd.Flags()
	membersToAdd, err := flags.GetStringArray(addMembersFlag)
	membersToRemove, err := flags.GetStringArray(removeMembersFlag)

	projectName := ""
	if len(args) > 0 {
		projectName = args[0]
	}
	projectId := GetProjectId(projectName)

	if len(membersToAdd) > 0 {
		addReq := &projects.AddProjectMembersRequest{
			Emails:    membersToAdd,
			ProjectId: projectId,
		}
		_, err = client.AddProjectMembers(apiContext, addReq)
		common.CliExit(err)
	}

	if len(membersToRemove) > 0 {
		removeReq := &projects.RemoveProjectMembersRequest{
			Emails:    membersToRemove,
			ProjectId: projectId,
		}
		_, err = client.RemoveProjectMembers(apiContext, removeReq)
		common.CliExit(err)
	}
	return
}
