package schema

import (
	"github.com/spf13/cobra"
)

var longDoc = `### Usage`

func DeleteCmd() *cobra.Command {
	schema := &cobra.Command{
		Use:   "schema (reference)",
		Short: "Delete Schema by reference",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the schema reference
		DisableAutoGenTag: true,
		ValidArgsFunction: RefsCompletion,
	}

	return schema
}

func ActivateCmd() *cobra.Command {
	schema := &cobra.Command{
		Use:   "schema (reference)",
		Short: "Set the state of a Schema to ACTIVATED",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			activate(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the schema reference
		DisableAutoGenTag: true,
		ValidArgsFunction: RefsCompletion,
	}

	return schema
}

func ArchiveCmd() *cobra.Command {
	schema := &cobra.Command{
		Use:   "schema (reference)",
		Short: "Set the state of a Schema to ARCHIVED",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			archive(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the schema reference
		DisableAutoGenTag: true,
		ValidArgsFunction: RefsCompletion,
	}

	return schema
}

func GetCmd() *cobra.Command {
	getSchema := &cobra.Command{
		Use:               "schema [name]",
		Short:             "Get schema by name",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: RefsCompletion,
	}
	flags := getSchema.Flags()
	flags.String(kafkaClusterFlag, "", "if present, find the corresponding Confluent ID for the given Kafka cluster")

	return getSchema
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "schemas",
		Short:             "List schemas",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}

func CreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:               "schema (handle/name/version)",
		Short:             "create a schema",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd, &args[0])
		},
		Args: cobra.ExactArgs(1),
	}
	flags := createCmd.Flags()
	flags.String(definitionFlag, "", "filename of the definition (either Simple Schema, Avro or Json)")
	_ = createCmd.MarkFlagRequired(definitionFlag)
	flags.String(schemaTypeFlag, "AVRO", "type of schema")
	flags.Bool(publicFlag, false, "should the schema become public")
	return createCmd
}
