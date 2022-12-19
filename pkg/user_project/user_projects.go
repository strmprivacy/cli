package user_project

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
)

const ActiveProjectFilename = "active_projects.json"

var Projects UsersProjects

// UsersProjects is the printed json format of the different active projects
// per past or currently logged-in user
type UsersProjects struct {
	Users []UserProject `json:"users"`
}

type UserProject struct {
	Email         string `json:"email"`
	ActiveProject string `json:"active_project"`
}

func (projects *UsersProjects) Init(projectId string) {
	projects.Users = []UserProject{{Email: GetUserEmail(), ActiveProject: projectId}}
}

func (projects *UsersProjects) GetCurrentProjectByEmail() string {
	activeProject := ""
	email := GetUserEmail()
	for _, user := range projects.Users {
		if user.Email == email {
			activeProject = user.ActiveProject
		}
	}
	return activeProject
}

func (projects *UsersProjects) SetActiveProject(project string) {
	email := GetUserEmail()
	// The `Project` var is already initialized when setActiveProject
	// is called
	added := false
	for _, user := range projects.Users {
		if user.Email == email {
			user.ActiveProject = project
			added = true
		}
	}

	if !added {
		projects.Users = append(projects.Users, UserProject{
			Email:         email,
			ActiveProject: project,
		})
	}
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

func GetActiveProject() string {
	activeProjectFilePath := path.Join(common.ConfigPath, ActiveProjectFilename)

	bytes, err := os.ReadFile(activeProjectFilePath)
	common.CliExit(err)
	activeProjects := UsersProjects{}
	_ = json.Unmarshal(bytes, &activeProjects)
	activeProject := activeProjects.GetCurrentProjectByEmail()
	Projects = activeProjects
	log.Infoln("Current active user_project is: " + activeProject)
	return activeProject
}
