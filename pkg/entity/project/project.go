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
	"strmprivacy/strm/pkg/user_project"
	"strmprivacy/strm/pkg/util"
)

var client projects.ProjectsServiceClient
var apiContext context.Context

const (
	descriptionFlag   = "description"
	addMembersFlag    = "add-member"
	removeMembersFlag = "remove-member"
)

var projectsAndActiveProject *ProjectsAndActiveProject

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = projects.NewProjectsServiceClient(clientConnection)
}

func ListProjects() *projects.ListProjectsResponse {
	return &projects.ListProjectsResponse{
		Projects: ListProjectsWithActive().Projects,
	}
}

func ListProjectsWithActive() ProjectsAndActiveProject {
	if projectsAndActiveProject == nil {
		if apiContext == nil {
			common.CliExit(errors.New(fmt.Sprintf("No login information found. Use: %s auth login first.", common.RootCommandName)))
		}

		req := &projects.ListProjectsRequest{}
		response, err := client.ListProjects(apiContext, req)
		common.CliExit(err)

		projectsAndActiveProject = &ProjectsAndActiveProject{
			Projects:      response.Projects,
			activeProject: user_project.GetActiveProject(),
		}
	}

	return *projectsAndActiveProject
}

func GetProject(projectName string) ProjectsAndActiveProject {
	for _, project := range ListProjectsWithActive().Projects {
		if project.Name == projectName {
			return ProjectsAndActiveProject{
				Projects:      []*entities.Project{project},
				activeProject: user_project.GetActiveProject(),
			}
		}
	}
	return ProjectsAndActiveProject{}
}

func create(projectName *string, cmd *cobra.Command) ProjectsAndActiveProject {
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
	return ProjectsAndActiveProject{
		Projects:      []*entities.Project{response.Project},
		activeProject: user_project.GetActiveProject(),
	}
}

func GetProjectIdFromName(projectName string) string {
	activeProject := ""
	if projectName == "" {
		activeProject = user_project.GetActiveProject()
	} else {
		activeProject = projectName
	}
	resolvedProject := GetProject(activeProject)
	if len(resolvedProject.Projects) == 0 {
		common.CliExit(errors.New(fmt.Sprintf("Project '%v' does not exist, or you do not have access "+
			"to it.", activeProject)))
	}
	return resolvedProject.Projects[0].Id
}

func GetProjectId(cmd *cobra.Command) string {
	flags := cmd.Flags()
	projectName, _ := flags.GetString(common.ProjectNameFlag)

	return GetProjectIdFromName(projectName)
}

func manage(args []string, cmd *cobra.Command) {
	flags := cmd.Flags()
	membersToAdd, err := flags.GetStringArray(addMembersFlag)
	membersToRemove, err := flags.GetStringArray(removeMembersFlag)

	projectName := ""
	if len(args) > 0 {
		projectName = args[0]
	}
	projectId := GetProjectIdFromName(projectName)

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

func get(projectName string) ProjectsAndActiveProject {
	projectId := GetProjectIdFromName(projectName)
	req := &projects.GetProjectRequest{ProjectId: projectId}

	response, err := client.GetProject(apiContext, req)
	common.CliExit(err)

	return ProjectsAndActiveProject{
		Projects:      []*entities.Project{response.Project},
		activeProject: user_project.GetActiveProject(),
	}
}

func del(projectName string) {
	projectId := GetProjectIdFromName(projectName)
	req := &projects.DeleteProjectRequest{
		ProjectId: projectId,
	}
	_, err := client.DeleteProject(apiContext, req)
	common.CliExit(err)
	return
}

type ProjectsAndActiveProject struct {
	Projects      []*entities.Project
	activeProject string
}
