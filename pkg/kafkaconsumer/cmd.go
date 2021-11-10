package kafkaconsumer

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/kafka_exporter"
)

var Cmd = &cobra.Command{
	Use:   "kafka (kafka-exporter-name)",
	Short: "Read events via the kafka consumer (not for production purposes)",
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, &args[0])
	},
	Args: cobra.ExactArgs(1), // the kafka-exporter name
	ValidArgsFunction: kafka_exporter.NamesCompletion,
}

func init() {
	flags := Cmd.Flags()
	flags.String(common.ClientIdFlag, "", "Client id to be used for receiving data")
	flags.String(common.ClientSecretFlag, "", "Client secret to be used for receiving data")
	flags.String(GroupIdFlag, "", "Kafka consumer group id. Uses a random value when not set")
}