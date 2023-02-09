package data_connector

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"os"
	"strings"
	"strmprivacy/strm/pkg/common"
)

const (
	longDocJdbc = `Creates a data connector for a jdbc database.

### Usage`
)

func createJdbcConnectionCmd() *cobra.Command {
	jdbcConnector := &cobra.Command{
		Use:               "jdbc (data-connector-name) \"<JDBC URL>\"",
		Short:             "Create a Data Connector for jdbc database",
		Long:              longDocJdbc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinterForType(cmd, cmd.Parent().Parent().Name())
		},
		Run: func(cmd *cobra.Command, args []string) {
			dataConnectorName := &args[0]
			jdbcUrl := &args[1]
			databaseType := determineDatabaseType(*jdbcUrl)
			privateKey := loadPrivateKey(*jdbcUrl, databaseType)

			dataConnector := &entities.DataConnector{
				Ref: ref(dataConnectorName),
				Location: &entities.DataConnector_JdbcLocation{
					JdbcLocation: &entities.JdbcLocation{
						JdbcUri:      *jdbcUrl,
						DatabaseType: databaseType,
						PrivateKey:   privateKey,
					},
				},
			}
			create(dataConnector, cmd)
		},
		Args: cobra.ExactArgs(2),
	}

	return jdbcConnector
}

func determineDatabaseType(jdbcUrl string) entities.DatabaseType {
	splitUrl := strings.Split(jdbcUrl, ":")

	var databaseType entities.DatabaseType

	switch db := splitUrl[1]; db {
	case "postgres":
		databaseType = entities.DatabaseType_POSTGRES
	case "datadirect":
		databaseType = entities.DatabaseType_BIGQUERY
	case "mysql":
		databaseType = entities.DatabaseType_MYSQL
	case "mongodb":
		databaseType = entities.DatabaseType_MONGODB
	default:
		common.CliExit(errors.New(fmt.Sprintf("Unknown jdbc url (supported types: %s, %s, %s, %s)",
			entities.DatabaseType_MYSQL.String(),
			entities.DatabaseType_POSTGRES.String(),
			entities.DatabaseType_BIGQUERY.String(),
			entities.DatabaseType_MONGODB.String(),
		)))
	}

	return databaseType
}

func loadPrivateKey(jdbcUrl string, databaseType entities.DatabaseType) string {
	var privateKey string
	if databaseType == entities.DatabaseType_BIGQUERY {
		filePath := strings.Split(jdbcUrl, "ServiceAccountPrivateKey=")[1]
		content, err := os.ReadFile(filePath)
		common.CliExit(err)
		privateKey = string(content)
	} else {
		privateKey = ""
	}
	return privateKey
}
