package io.streammachine.api.cli

import com.github.ajalt.clikt.completion.completionOption
import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.context
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.options.flag
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.versionOption
import io.streammachine.api.cli.Strm.Companion.COMMAND
import io.streammachine.api.cli.commands.*
import io.streammachine.api.cli.common.ColorHelpFormatter
import io.streammachine.api.cli.common.Common
import io.streammachine.api.cli.common.UpdateCheck.printUpdateMessageIfAvailable
import io.streammachine.api.cli.common.initializeFuel

fun main(args: Array<String>) = Strm()
    .context { helpFormatter = ColorHelpFormatter }
    .completionOption(help = "Generate the completion script for the Stream Machine CLI. Usage = $COMMAND --generate-completion [bash zsh fish] > /completion/script/location/strm-completions.sh")
    .versionOption(Common.VERSION, names = setOf("-v", "--version"), message = { "Stream Machine CLI version: $it" })
    .main(args)

//fun main(args: Array<String>) = printUpdateMessageIfAvailable()

class Strm : CliktCommand(
    name = COMMAND,
    help = "Command Line Interface for https://streammachine.io",
    epilog = "Docs: https://docs.streammachine.io - Gitter: https://gitter.im/stream-machine"
) {
    companion object {
        internal const val COMMAND = "strm"
    }

    private val verbose by option("--verbose", hidden = true).flag(default = false)

    init {
        subcommands(
            Authentication(),
            Streams(),
            Outputs(),
            Exporters(),
            ConsentLevels(),
            Sinks()
        )
    }

    override fun run() {
        Common.VERBOSE_LOGGING = verbose
        printUpdateMessageIfAvailable()
        initializeFuel()
    }
}
