package user_project

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
)

const activeProjectFilename = "active_projects.json"

var ActiveProjectFilepath = path.Join(common.ConfigPath(), activeProjectFilename)

var Projects *UsersProjectsContext

// UsersProjectsContext is the printed json format of the different active projects
// per past or currently logged-in user
type UsersProjectsContext struct {
	Users []UserProjectContext `json:"users"`
}

type UserProjectContext struct {
	Email         string `json:"email"`
	ActiveProject string `json:"active_project"`
	ZedToken      string `json:"zed_token"`
}

func (projects *UsersProjectsContext) GetCurrentProjectByEmail() string {
	activeProject := ""
	email := GetUserEmail()
	for _, user := range projects.Users {
		if user.Email == email {
			activeProject = user.ActiveProject
		}
	}
	return activeProject
}

func (projects *UsersProjectsContext) SetActiveProject(project string) {
	email := GetUserEmail()
	added := false
	for index, user := range projects.Users {
		if user.Email == email {
			(*projects).Users[index].ActiveProject = project
			added = true
		}
	}

	if !added {
		projects.Users = append(projects.Users, UserProjectContext{
			Email:         email,
			ActiveProject: project,
		})
	}

	storeUserProjectContext()
}

func GetUserEmail() string {
	if auth.Auth.Email != "" {
		return auth.Auth.Email
	}
	if err := auth.Auth.LoadLogin(); err != nil {
		common.CliExit(err)
	}
	return auth.Auth.Email
}

func initializeUsersProjectsContext() {
	if Projects == nil {
		activeProjectFilePath := path.Join(common.ConfigPath(), activeProjectFilename)

		bytes, err := os.ReadFile(activeProjectFilePath)
		common.CliExit(err)
		activeProjects := UsersProjectsContext{}
		_ = json.Unmarshal(bytes, &activeProjects)
		Projects = &activeProjects
	}
}

func GetActiveProject() string {
	initializeUsersProjectsContext()
	activeProject := Projects.GetCurrentProjectByEmail()
	log.Infoln("Current active project is: " + activeProject)
	return activeProject
}

func SetZedToken(zedToken string) {
	initializeUsersProjectsContext()
	email := GetUserEmail()
	for index, user := range Projects.Users {
		// If there is no entry for the user, a zed token will be added the next time, when it is present
		if user.Email == email {
			(*Projects).Users[index].ZedToken = zedToken
		}
	}

	storeUserProjectContext()
}

func GetZedToken() *string {
	initializeUsersProjectsContext()
	email := GetUserEmail()
	for _, user := range Projects.Users {
		if user.Email == email && user.ZedToken != "" {
			return &user.ZedToken
		}
	}
	return nil
}

func storeUserProjectContext() {
	projects, err := json.Marshal(Projects)
	common.CliExit(err)

	err = os.WriteFile(ActiveProjectFilepath, projects, 0644)
	common.CliExit(err)
}
