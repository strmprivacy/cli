package kafka_user

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
)

const (
	nameFlag = "name"
	saveFlag = "save"
)

var content, _ = ioutil.ReadFile("pkg/entity/kafka_user/docstring.md")

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-user [name ...]",
		Short:             "Delete one or more Kafka users",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for i, _ := range args {
				del(&args[i])
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: namesCompletion,
	}
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-user [name]",
		Short:             "Get Kafka user",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: namesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-users [kafka-exporter-name]",
		Short:             "List Kafka users",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the kafka exporter name
		ValidArgsFunction: kafka_exporter.NamesCompletion,
	}
}

func CreateCmd() *cobra.Command {
	kafkaUser := &cobra.Command{
		Use:               "kafka-user [exporter-name]",
		Short:             "Create a Kafka user on a Kafka exporter",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			streamName := &args[0]
			create(streamName, cmd)

		},
		Args:              cobra.ExactArgs(1), // the kafka-exporter name
		ValidArgsFunction: kafka_exporter.NamesCompletion,
	}
	flags := kafkaUser.Flags()
	flags.Bool(saveFlag, false, "save the result in the config directory")

	return kafkaUser
}
