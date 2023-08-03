package stream

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/streams/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.ListCommandName:   listTablePrinter{},
			common.OutputFormatTable + common.GetCommandName:    getTablePrinter{},
			common.OutputFormatTable + common.DeleteCommandName: deletePrinter{},
			common.OutputFormatTable + common.CreateCommandName: createTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatPlain + common.GetCommandName:    getPlainPrinter{},
			common.OutputFormatPlain + common.DeleteCommandName: deletePrinter{},
			common.OutputFormatPlain + common.CreateCommandName: createPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}
type createTablePrinter struct{}

type deletePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*streams.ListStreamsResponse)
	printTable(listResponse.Streams)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*streams.GetStreamResponse)
	printTable([]*v1.StreamTree{getResponse.StreamTree})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*streams.CreateStreamResponse)
	printTable([]*v1.StreamTree{{Stream: createResponse.Stream}})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*streams.ListStreamsResponse)
	printPlain(listResponse.Streams)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*streams.GetStreamResponse)
	printPlain([]*v1.StreamTree{getResponse.StreamTree})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*streams.CreateStreamResponse)
	printPlain([]*v1.StreamTree{{Stream: createResponse.Stream}})
}

func (p deletePrinter) Print(data interface{}) {
	fmt.Println("Stream has been deleted")
}

func printTable(streamTreeArray []*v1.StreamTree) {
	rows := make([]table.Row, 0, len(streamTreeArray))
	// topics are only set for CCD installations. We don't yet have the installation
	// id inside the StreamTree, so we have this workaround to only show a topic
	// column if any of the shown streams has a topic
	hasTopics := false

	var policies = map[string]string{}

	for _, stream := range streamTreeArray {
		if stream.Stream.PolicyId != "" && len(policies) == 0 {
			_, p := policy.PoliciesNameIdMap()
			policies = p
		}

		row := table.Row{
			stream.Stream.Ref.Name,
			len(stream.Stream.LinkedStream) != 0,
			stream.Stream.ConsentLevels,
			stream.Stream.Enabled,
			policies[stream.Stream.PolicyId],
		}
		if len(stream.Stream.KafkaTopic) != 0 {
			row = append(row, stream.Stream.KafkaTopic)
			hasTopics = true
		}
		rows = append(rows, row)
	}

	headers := table.Row{
		"Stream",
		"Derived",
		"Purposes",
		"Enabled",
		"Policy Name",
	}
	if hasTopics {
		headers = append(headers, "Kafka Topic")
	}
	util.RenderTable(headers, rows)
}

func printPlain(streamTreeArray []*v1.StreamTree) {
	var names string
	lastIndex := len(streamTreeArray) - 1

	for index, stream := range streamTreeArray {
		names = names + stream.Stream.Ref.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
