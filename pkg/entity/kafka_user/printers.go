package kafka_user

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_users/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case constants.OutputFormatJson:
		return util.GenericPrettyJsonPrinter{}
	case constants.OutputFormatJsonRaw:
		return util.GenericRawJsonPrinter{}
	case constants.OutputFormatTable:
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
	case constants.OutputFormatPlain:
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
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, constants.OutputFormatFlagAllowedValuesText))
		return nil
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
	listResponse, _ := (data).(*kafka_users.ListKafkaUsersResponse)
	printTable(listResponse.KafkaUsers)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_users.GetKafkaUserResponse)
	printTable([]*entities.KafkaUser{getResponse.KafkaUser})
}

func (p createTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*kafka_users.CreateKafkaUserResponse)
	printTable([]*entities.KafkaUser{createResponse.KafkaUser})
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_users.ListKafkaUsersResponse)
	printPlain(listResponse.KafkaUsers)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_users.GetKafkaUserResponse)
	printPlain([]*entities.KafkaUser{getResponse.KafkaUser})
}

func (p createPlainPrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*kafka_users.CreateKafkaUserResponse)
	printPlain([]*entities.KafkaUser{createResponse.KafkaUser})
}

func (p deletePrinter) Print(_ proto.Message) {
	fmt.Println("Kafka Exporter has been deleted")
}

func printTable(kafkaUsers []*entities.KafkaUser) {
	rows := make([]table.Row, 0, len(kafkaUsers))

	for _, user := range kafkaUsers {
		rows = append(rows, table.Row{
			user.Ref.Name,
			user.KafkaExporterName,
			user.Topic,
			user.ClientId,
			user.ClientSecret,
		})
	}

	util.RenderTable(
		table.Row{
			"Kafka User",
			"Kafka Exporter",
			"Topic",
			"Client ID",
			"Client Secret",
		},
		rows,
	)
}

func printPlain(kafkaUsers []*entities.KafkaUser) {
	var names string
	lastIndex := len(kafkaUsers) - 1

	for index, user := range kafkaUsers {
		names = names + user.Ref.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
