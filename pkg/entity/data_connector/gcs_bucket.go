package data_connector

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
)

const (
	longDocGcs = `Creates a data connector for a Google Cloud Storage bucket. Service account credentials are required.

### Usage`
)

func createGcsBucketCmd() *cobra.Command {
	gcsBucket := &cobra.Command{
		Use:               "gcs [data-connector-name] [bucket-name]",
		Short:             "Create a Data Connector for a Google Cloud Storage Bucket",
		Long:              longDocGcs,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinterForType(cmd, cmd.Parent().Parent().Name())
		},
		Run: func(cmd *cobra.Command, args []string) {
			credentials := readCredentialsFile(cmd.Flags())
			dataConnectorName := &args[0]
			bucketName := &args[1]

			dataConnector := &entities.DataConnector{
				Ref: ref(dataConnectorName),
				Location: &entities.DataConnector_GoogleCloudStorageBucket{
					GoogleCloudStorageBucket: &entities.GoogleCloudStorageBucketLocation{
						BucketName:  *bucketName,
						Credentials: credentials,
					},
				},
			}
			create(dataConnector)
		},
		Args: cobra.ExactArgs(2),
	}

	flags := gcsBucket.Flags()
	flags.String(credentialsFileFlag, "", "file with service account JSON key credentials")
	_ = gcsBucket.MarkFlagRequired(credentialsFileFlag)
	_ = gcsBucket.MarkFlagFilename(credentialsFileFlag, "json")
	return gcsBucket
}
