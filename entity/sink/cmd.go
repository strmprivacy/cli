package sink

import "github.com/spf13/cobra"

const (
	sinkTypeFlag        = "sink-type"
	credentialsFileFlag = "credentials-file"
)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sink [name]",
		Short: "get sink by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: ExistingNamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sinks",
		Short: "List sinks",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sink [name]",
		Short: "Delete sinks",
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool("recursive")
			del(&args[0], recursive)
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ExistingNamesCompletion,
	}
}
func CreateCmd() *cobra.Command {
	sink := &cobra.Command{
		Use:   "sink [sink-name] [bucket-name]",
		Short: "Create sink",
		Run: func(cmd *cobra.Command, args []string) {
			create(&args[0], &args[1], cmd)
		},
		Args: cobra.ExactArgs(2),
	}

	flags := sink.Flags()
	flags.String(sinkTypeFlag, "", "S3 or GCLOUD")
	flags.String(credentialsFileFlag, "", "file with credentials")
	_ = sink.MarkFlagRequired(credentialsFileFlag)
	_ = sink.MarkFlagFilename(credentialsFileFlag, "json")
	return sink

}