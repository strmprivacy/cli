package context

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
	"os"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
	"strings"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), constants.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Name()]

	if p == nil {
		var allowedValues string

		switch command.Name() {
		case entityInfoCommandName:
			allowedValues = constants.ContextOutputFormatFlagAllowedValuesText
		case configCommandName:
			allowedValues = constants.ConfigOutputFormatFlagAllowedValuesText
		}

		common.CliExit(fmt.Sprintf("Output format '%v' is not supported for '%v'. Allowed values: %v", command.CommandPath(), outputFormat, allowedValues))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return map[string]util.Printer{
		constants.OutputFormatJsonRaw + entityInfoCommandName:  jsonRawPrinter{},
		constants.OutputFormatJson + entityInfoCommandName:     jsonPrettyPrinter{},
		constants.OutputFormatFilepath + entityInfoCommandName: filepathPrinter{},
		constants.OutputFormatPlain + configCommandName:        plainPrinter{},
	}
}

type filepathPrinter struct{}
type jsonRawPrinter struct{}
type jsonPrettyPrinter struct{}
type plainPrinter struct{}

func (p filepathPrinter) Print(data interface{}) {
	entity, _ := (data).(savedEntity)

	fmt.Println(entity.Path)
}

func (p jsonRawPrinter) Print(data interface{}) {
	entity, _ := (data).(savedEntity)
	compactJson := util.CompactJson(entity.Contents)
	fmt.Println(string(compactJson.Bytes()))
}

func (p jsonPrettyPrinter) Print(data interface{}) {
	entity, _ := (data).(savedEntity)
	prettyJson := util.PrettifyJson(util.CompactJson(entity.Contents))
	fmt.Println(string(prettyJson.Bytes()))
}

func (p plainPrinter) Print(data interface{}) {
	config, _ := (data).(configuration)

	l := list.NewWriter()
	l.SetOutputMirror(os.Stdout)

	fmt.Println(fmt.Sprintf("Configuration directory: %v", config.ConfigPath))
	fmt.Println(fmt.Sprintf("Configuration file: %v", config.ConfigFilepath))
	fmt.Println(fmt.Sprintf("Configuration file contents: \n\n    %v", strings.ReplaceAll(string(config.Contents), "\n", "\n    ")))

	for _, entity := range config.SavedEntities {
		l.AppendItem(entity)
	}

	fmt.Println("Saved entities:")
	l.Render()
}
