package kafkaconsumer

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
)

var Cmd = &cobra.Command{
	Use:               "kafka (kafka-exporter-name)",
	Short:             "Read events via the kafka consumer (not for production purposes)",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, &args[0])
	},
	Args:              cobra.ExactArgs(1), // the kafka-exporter name
	ValidArgsFunction: kafka_exporter.NamesCompletion,
}

func init() {
	flags := Cmd.Flags()
	flags.String(common.ClientIdFlag, "", "client id to be used for receiving data")
	flags.String(common.ClientSecretFlag, "", "client secret to be used for receiving data")
	flags.String(GroupIdFlag, "", "kafka consumer group id. Uses a random value when not set")
	flags.String(KafkaBootstrapHostFlag, "export-bootstrap.kafka.strmprivacy.io:9092", "Kafka bootstrap brokers, separated by comma")
}
