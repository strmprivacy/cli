package data_connector

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"strmprivacy/strm/pkg/util"
)

const (
	assumeRoleArnFlag = "assume-role-arn"
	longDocS3         = `Creates a data connector for an AWS S3 bucket. An ARN can be specified in case a role should be assumed.

### Usage`
)

func createS3BucketCmd() *cobra.Command {
	s3Bucket := &cobra.Command{
		Use:               "s3 (data-connector-name) (bucket-name)",
		Short:             "Create a Data Connector for an AWS S3 Bucket",
		Long:              longDocS3,
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
				Location: &entities.DataConnector_S3Bucket{
					S3Bucket: &entities.AwsS3BucketLocation{
						BucketName:    *bucketName,
						Credentials:   credentials,
						AssumeRoleArn: util.GetStringAndErr(cmd.Flags(), assumeRoleArnFlag),
					},
				},
			}
			create(dataConnector, cmd)
		},
		Args: cobra.ExactArgs(2),
	}

	flags := s3Bucket.Flags()
	flags.String(credentialsFileFlag, "", "file with JSON AWS Access Key Credentials")
	flags.String(assumeRoleArnFlag, "", "ARN of the role to assume")
	_ = s3Bucket.MarkFlagRequired(credentialsFileFlag)
	_ = s3Bucket.MarkFlagFilename(credentialsFileFlag, "json")
	return s3Bucket
}
