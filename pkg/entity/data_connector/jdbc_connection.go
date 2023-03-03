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
	longDocJdbc = `Creates a data connector for a JDBC-compatible database.

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
	var databaseType entities.DatabaseType
	switch {
	case strings.Contains(jdbcUrl, "postgres"):
		databaseType = entities.DatabaseType_POSTGRES
	case strings.Contains(jdbcUrl, "bigquery"):
		databaseType = entities.DatabaseType_BIGQUERY
	case strings.Contains(jdbcUrl, "mysql"):
		databaseType = entities.DatabaseType_MYSQL
	case strings.Contains(jdbcUrl, "snowflake"):
		databaseType = entities.DatabaseType_SNOWFLAKE
	default:
		common.CliExit(errors.New(fmt.Sprintf("Unknown jdbc url (supported types: %s, %s, %s)",
			entities.DatabaseType_MYSQL.String(),
			entities.DatabaseType_POSTGRES.String(),
			entities.DatabaseType_BIGQUERY.String(),
		)))
	}

	return databaseType
}

func loadPrivateKey(jdbcUrl string, databaseType entities.DatabaseType) string {
	var privateKey string
	if databaseType == entities.DatabaseType_BIGQUERY {
		var filePath string
		for _, s := range strings.Split(jdbcUrl, ";") {
			if strings.Contains(s, "OAuthPvtKeyPath") {
				_, filePath, _ = strings.Cut(s, "OAuthPvtKeyPath=")
				break
			} else if strings.Contains(s, "ServiceAccountPrivateKey=") {
				_, filePath, _ = strings.Cut(s, "ServiceAccountPrivateKey=")
				break
			}
		}
		content, err := os.ReadFile(filePath)
		common.CliExit(err)
		privateKey = string(content)
	} else {
		privateKey = ""
	}
	return privateKey
}
