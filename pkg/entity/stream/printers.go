package stream

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	v1 "github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/streams/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case "json":
		return util.GenericPrettyJsonPrinter{}
	case "json-raw":
		return util.GenericRawJsonPrinter{}
	case "table":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return listTablePrinter{}
		case constants.GetCommandName:
			return getTablePrinter{}
		case constants.DeleteCommandName:
			return deletePrinter{}
		case constants.CreateCommandName:
			return createTablePrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	case "plain":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return listPlainPrinter{}
		case constants.GetCommandName:
			return getPlainPrinter{}
		case constants.DeleteCommandName:
			return deletePrinter{}
		case constants.CreateCommandName:
			return createPlainPrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	default:
		return util.GenericPrettyJsonPrinter{}
	}
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}
type createTablePrinter struct{}

type deletePrinter struct{}

func (p listTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*streams.ListStreamsResponse)
	printTable(listResponse.Streams)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*streams.GetStreamResponse)
	printTable([]*v1.StreamTree{getResponse.StreamTree})
}

func (p createTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*streams.CreateStreamResponse)
	printTable([]*v1.StreamTree{{Stream: createResponse.Stream}})
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*streams.ListStreamsResponse)
	printPlain(listResponse.Streams)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*streams.GetStreamResponse)
	printPlain([]*v1.StreamTree{getResponse.StreamTree})
}

func (p createPlainPrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*streams.CreateStreamResponse)
	printPlain([]*v1.StreamTree{{Stream: createResponse.Stream}})
}

func (p deletePrinter) Print(_ proto.Message) {
	fmt.Println("Stream has been deleted")
}

func printTable(streamTreeArray []*v1.StreamTree) {
	rows := make([]table.Row, 0, len(streamTreeArray))

	for _, stream := range streamTreeArray {
		var consentLevelType string

		if stream.Stream.ConsentLevelType != entities.ConsentLevelType_CONSENT_LEVEL_TYPE_UNSPECIFIED {
			consentLevelType = stream.Stream.ConsentLevelType.String()
		} else {
			consentLevelType = ""
		}

		rows = append(rows, table.Row{
			stream.Stream.Ref.Name,
			len(stream.Stream.LinkedStream) != 0,
			consentLevelType,
			stream.Stream.ConsentLevels,
			stream.Stream.Enabled,
			stream.KeyStream != nil,
		})
	}

	util.RenderTable(
		table.Row{
			"Stream",
			"Derived",
			"Consent Level Type",
			"Consent Levels",
			"Enabled",
			"Has Key Stream",
		},
		rows,
	)
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
