package context

import (
	"fmt"
	"io/ioutil"
	"path"
	"streammachine.io/strm/pkg/auth"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
	"strings"
)

type configuration struct {
	ConfigPath     string
	ConfigFilepath string
	Contents       []byte
	SavedEntities  []string
}

func showConfiguration() {
	configFilepath := findConfigFile()
	contents, err := ioutil.ReadFile(configFilepath)
	common.CliExit(err)

	configuration := configuration{
		ConfigPath:     util.ConfigPath,
		ConfigFilepath: configFilepath,
		Contents:       contents,
		SavedEntities:  listSavedEntities(util.ConfigPath),
	}

	printer.Print(configuration)
}

type savedEntity struct {
	Path     string
	Contents []byte
}

func entityInfo(args []string) {
	filepath := path.Join(util.ConfigPath, args[0]+".json")
	contents, err := ioutil.ReadFile(filepath)
	common.CliExit(err)

	entity := savedEntity{Path: filepath, Contents: contents}
	printer.Print(entity)
}

func listSavedEntities(p string) []string {
	files, err := ioutil.ReadDir(p)
	common.CliExit(err)

	entityTypeTemplate := path.Base(p) + "/%v"
	var entities = make([]string, 0)

	for _, f := range files {
		if f.IsDir() {
			entities = append(entities, listSavedEntities(path.Join(p, f.Name()))...)
		} else if strings.HasSuffix(f.Name(), "json") && !strings.HasPrefix(f.Name(), auth.StrmCredsFilePrefix) {
			entities = append(entities, fmt.Sprintf(entityTypeTemplate, strings.Replace(f.Name(), ".json", "", -1)))
		}
	}

	return entities
}

func findConfigFile() string {
	files, err := ioutil.ReadDir(util.ConfigPath)
	common.CliExit(err)

	var configFilepath string

	for _, f := range files {
		if (strings.HasSuffix(f.Name(), "json") || strings.HasSuffix(f.Name(), "yaml")) && strings.HasPrefix(f.Name(), constants.DefaultConfigFilename) && len(f.Name()) == constants.DefaultConfigFilenameLength {
			configFilepath = path.Join(util.ConfigPath, f.Name())
		}
	}

	return configFilepath
}
