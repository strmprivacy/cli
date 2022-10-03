package data_connector

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"strmprivacy/strm/pkg/util"
)

const (
	storageAccountUriFlag = "storage-account-uri"
	tenantIdFlag          = "tenant-id"
	clientIdFlag          = "client-id"
	clientSecretFlag      = "client-secret"
	longDocBlobStorage    = `Creates a data connector for an Azure Blob Storage container. Authentication is based on
Client Secret Credentials, i.e. of an Application Principal.

### Usage`
	projectName = "project"
)

func createAzureBlobStorageCmd() *cobra.Command {
	azureBlobStorage := &cobra.Command{
		Use:               "azure-blob-storage [data-connector-name] [container-name]",
		Short:             "Create a Data Connector for an Azure Blob Storage Container",
		Long:              longDocBlobStorage,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinterForType(cmd, cmd.Parent().Parent().Name())
		},
		Run: func(cmd *cobra.Command, args []string) {
			dataConnectorName := &args[0]
			containerName := &args[1]

			dataConnector := &entities.DataConnector{
				Ref: ref(dataConnectorName),
				Location: &entities.DataConnector_AzureBlobStorageContainer{
					AzureBlobStorageContainer: &entities.AzureBlobStorageContainerLocation{
						StorageAccountUri: util.GetStringAndErr(cmd.Flags(), storageAccountUriFlag),
						ContainerName:     *containerName,
						ClientSecretCredential: &entities.AzureClientSecretCredential{
							TenantId:     util.GetStringAndErr(cmd.Flags(), tenantIdFlag),
							ClientId:     util.GetStringAndErr(cmd.Flags(), clientIdFlag),
							ClientSecret: util.GetStringAndErr(cmd.Flags(), clientSecretFlag),
						},
					},
				},
			}
			create(dataConnector, nil)
		},
		Args: cobra.ExactArgs(2),
	}

	flags := azureBlobStorage.Flags()
	flags.String(storageAccountUriFlag, "", "full URI of the Azure Storage account, e.g. \"https://myaccount.blob.core.windows.net\"")
	flags.String(tenantIdFlag, "", "tenant ID of the application principal's client secret credentials")
	flags.String(clientIdFlag, "", "client ID of the application principal's client secret credentials")
	flags.String(clientSecretFlag, "", "client secret of the application principal's client secret credentials")
	flags.String(projectName, "", `Project name to create resource in`)
	_ = azureBlobStorage.MarkFlagRequired(storageAccountUriFlag)
	_ = azureBlobStorage.MarkFlagRequired(tenantIdFlag)
	_ = azureBlobStorage.MarkFlagRequired(clientIdFlag)
	_ = azureBlobStorage.MarkFlagRequired(clientSecretFlag)
	return azureBlobStorage
}
