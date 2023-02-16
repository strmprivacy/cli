package data_connector

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_connectors/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	return configurePrinterForType(command, command.Parent().Name())
}

func configurePrinterForType(command *cobra.Command, commandType string) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)
	recursive, err := command.Flags().GetBool(common.RecursiveFlagName)

	if err != nil {
		recursive = false
	}
	p := availablePrinters(recursive)[outputFormat+commandType]

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters(recursive bool) map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.ListCommandName:   listTablePrinter{recursive: recursive},
			common.OutputFormatTable + common.GetCommandName:    getTablePrinter{recursive: recursive},
			common.OutputFormatTable + common.DeleteCommandName: deletePrinter{recursive: recursive},
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

type listTablePrinter struct {
	recursive bool
}
type getTablePrinter struct {
	recursive bool
}
type createTablePrinter struct{}

type deletePrinter struct {
	recursive bool
}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_connectors.ListDataConnectorsResponse)
	printTable(listResponse.DataConnectors, p.recursive)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*data_connectors.GetDataConnectorResponse)
	printTable([]*entities.DataConnector{getResponse.DataConnector}, p.recursive)
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*data_connectors.CreateDataConnectorResponse)
	printTable([]*entities.DataConnector{createResponse.DataConnector}, false)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_connectors.ListDataConnectorsResponse)
	printPlain(listResponse.DataConnectors)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*data_connectors.GetDataConnectorResponse)
	printPlain([]*entities.DataConnector{getResponse.DataConnector})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*data_connectors.CreateDataConnectorResponse)
	printPlain([]*entities.DataConnector{createResponse.DataConnector})
}

func (p deletePrinter) Print(data interface{}) {
	if p.recursive {
		fmt.Println("Data Connector and dependent entities have been deleted")
	} else {
		fmt.Println("Data Connector has been deleted")
	}
}

func printTable(dataConnectors []*entities.DataConnector, recursive bool) {
	rows := make([]table.Row, 0, len(dataConnectors))

	for _, dataConnector := range dataConnectors {
		var locationName string
		var locationType string

		switch location := dataConnector.Location.(type) {
		case *entities.DataConnector_S3Bucket:
			locationType = "Amazon S3 Bucket"
			locationName = location.S3Bucket.BucketName
		case *entities.DataConnector_GoogleCloudStorageBucket:
			locationType = "Google Cloud Storage Bucket"
			locationName = location.GoogleCloudStorageBucket.BucketName
		case *entities.DataConnector_AzureBlobStorageContainer:
			locationType = "Azure Blob Storage Container"
			locationName = location.AzureBlobStorageContainer.ContainerName
		case *entities.DataConnector_JdbcLocation:
			locationType = "JDBC Database"
			locationName = location.JdbcLocation.DatabaseType.String()
		}

		var row table.Row

		if recursive {
			row = table.Row{
				dataConnector.Ref.Name,
				locationType,
				locationName,
				len(dataConnector.DependentEntities.BatchExporters),
				len(dataConnector.DependentEntities.BatchJobs),
			}
		} else {
			row = table.Row{
				dataConnector.Ref.Name,
				locationType,
				locationName,
			}
		}

		rows = append(rows, row)
	}

	var header table.Row

	if recursive {
		header = table.Row{
			"Name",
			"Type",
			"Target Name",
			"# Batch Exporters",
			"# Batch Jobs",
		}
	} else {
		header = table.Row{
			"Name",
			"Type",
			"Target Name",
		}
	}

	util.RenderTable(header, rows)
}

func printPlain(dataConnectors []*entities.DataConnector) {
	var names []string

	for _, dataConnector := range dataConnectors {
		names = append(names, dataConnector.Ref.Name)
	}

	util.RenderPlain(strings.Join(names, "\n"))
}
