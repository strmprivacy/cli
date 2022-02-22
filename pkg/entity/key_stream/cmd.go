package key_stream

import (
	"github.com/spf13/cobra"
)

var longDoc = `Key Streams are a restricted feature. For now, enabling and disabling key streams can not be done through the Console or
CLI. Contact us for more information.

A Key Stream can be enabled on a stream and it contains all encryption keys that are used on this stream. Normally STRM
Privacy fully manages and stores the encryption keys that are used, but with a key stream, clients can get access to the
keys and decrypt events themselves.

Usage of key streams places a lot more responsibility in the hands of the client, so this feature requires careful
consideration before using.

With regard to the data flow, STRM Privacy generates a new encryption key whenever an event with a new "key link" (which
can be seen as a
"session" concept, in that it links separate events together) is received. This encryption key is stored internally,
rotated after a certain period and, if there is a key stream, put on the Key Stream approximately at the same time as
the event is put on the input stream.

### Usage`

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "key-stream [name]",
		Short:             "Get key stream by name",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "key-streams",
		Short:             "List key streams",
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
