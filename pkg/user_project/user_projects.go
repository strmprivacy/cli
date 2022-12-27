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
	added := false
	for index, user := range projects.Users {
		if user.Email == email {
			(*projects).Users[index].ActiveProject = project
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

func LoadActiveProject() {
	activeProjectFilePath := path.Join(common.ConfigPath(), ActiveProjectFilename)

	bytes, err := os.ReadFile(activeProjectFilePath)
	common.CliExit(err)
	activeProjects := UsersProjects{}
	_ = json.Unmarshal(bytes, &activeProjects)
	Projects = activeProjects
}

func GetActiveProject() string {
	LoadActiveProject()
	activeProject := Projects.GetCurrentProjectByEmail()
	log.Infoln("Current active project is: " + activeProject)
	return activeProject
}
