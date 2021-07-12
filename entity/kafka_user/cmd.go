package kafka_user

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/entity/kafka_exporter"
)

const (
	nameFlag = "name"
	saveFlag = "save"
)

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-user [name]",
		Short: "Delete a Kafka user",
		Run: func(cmd *cobra.Command, args []string) {
			del(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: existingNames,
	}
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-user [name]",
		Short: "Get Kafka user",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: existingNames,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-users [kafka-exporter-name]",
		Short: "List Kafka users",
		Run: func(cmd *cobra.Command, args []string) {
			list(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the kafka exporter name
		ValidArgsFunction: kafka_exporter.ExistingNamesCompletion,
	}
}

func CreateCmd() *cobra.Command {
	kafkaUser := &cobra.Command{
		Use:   "kafka-user [exporter-name]",
		Short: "Create a Kafka user on a Kafka exporter",
		Run: func(cmd *cobra.Command, args []string) {
			streamName := &args[0]
			create(streamName, cmd)

		},
		Args:              cobra.ExactArgs(1), // the kafka-exporter name
		ValidArgsFunction: kafka_exporter.ExistingNamesCompletion,
	}
	flags := kafkaUser.Flags()
	flags.Bool(saveFlag, false, "save the result in the config directory")

	return kafkaUser
}
