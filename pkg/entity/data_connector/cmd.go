package data_connector

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_connectors/v1"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

const (
	credentialsFileFlag = "credentials-file"
)

var longDoc = util.LongDocsUsage(`
A Data Connector represents a location from which data can be read, or to which data can be written.  For
example, an AWS S3 bucket, a Google Cloud Storage bucket or a JDBC database connection. By itself, a Data Connector does nothing.  A Data Connector
with valid credentials is required when creating a Batch Exporter or Batch Job.
`)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "data-connector (name)",
		Short:             "Get Data Connector by name",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			get(&args[0], recursive)
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: NamesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "data-connectors",
		Short:             "List Data Connectors",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			list(recursive)
		},
	}
}

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "data-connector (name ...)",
		Short:             "Delete one or more Data Connectors by name",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			for i := range args {
				del(&args[i], recursive)
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: NamesCompletion,
	}
}

func CreateCmd() *cobra.Command {
	dataConnector := &cobra.Command{
		Use:               "data-connector",
		Short:             "Create a Data Connector",
		Long:              longDoc,
		DisableAutoGenTag: true,
	}
	dataConnector.AddCommand(createS3BucketCmd())
	dataConnector.AddCommand(createGcsBucketCmd())
	dataConnector.AddCommand(createAzureBlobStorageCmd())
	dataConnector.AddCommand(createJdbcConnectionCmd())
	return dataConnector
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	req := &data_connectors.ListDataConnectorsRequest{
		ProjectId: common.ProjectId,
	}
	response, err := Client.ListDataConnectors(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.DataConnectors))
	for _, dataConnector := range response.DataConnectors {
		names = append(names, dataConnector.Ref.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
