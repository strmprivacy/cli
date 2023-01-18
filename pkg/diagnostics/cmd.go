package diagnostics

import "github.com/spf13/cobra"

const (
	quasiIdentifierFlagName    = "qi"
	sensitiveAttributeFlagName = "sa"
	dataFileFlagName           = "file"
	metricsFlagName            = "metrics"
)

func EvaluateCmd() *cobra.Command {
	diagnostics := &cobra.Command{
		Use:               "diagnostics",
		Short:             "Evaluate privacy diagnostics for your dataset",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			evaluate(cmd)
		},
		Args: cobra.ExactArgs(0),
	}
	flags := diagnostics.Flags()
	flags.StringP(dataFileFlagName, "F", "", `the path to the csv file to be evaluated`)
	flags.String(quasiIdentifierFlagName, "",
		`list of comma separated quasi identifier columns: [qi1,qi2,..]`)
	flags.String(sensitiveAttributeFlagName, "",
		`list of comma separated sensitive attributes columns: [sa1,sa2,..]`)
	flags.String(metricsFlagName, "k-anonymity, l-diversity, t-closeness",
		`list of metrics to calculate [k-anonymity, l-diversity, t-closeness]. defaults to all`)

	return diagnostics
}
