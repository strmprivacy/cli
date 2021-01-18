package io.streammachine.api.cli.commands

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.required
import com.github.ajalt.clikt.parameters.types.choice
import com.github.ajalt.clikt.parameters.types.enum
import com.github.ajalt.clikt.parameters.types.int
import com.github.ajalt.clikt.parameters.types.restrictTo
import com.github.kittinunf.fuel.gson.jsonBody
import com.github.kittinunf.fuel.httpDelete
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.fuel.httpPut
import io.streammachine.api.cli.common.ExporterCreateRequest
import io.streammachine.api.cli.common.ExporterType
import io.streammachine.api.cli.common.SinkType
import io.streammachine.api.cli.common.printResponse

class Exporters : CliktCommand(
    name = "exporters",
    help = "View exporters and their details."
) {
    init {
        subcommands(
            ExportersList(),
            ExportersGet(),
            ExportersCreate(),
            ExportersDelete()
        )
    }

    override fun run() = Unit
}

class ExportersList : CliktCommand(
    name = "list",
    help = "List your exporters"
) {
    override fun run() {
        "/v1/exporters"
            .httpGet()
            .printResponse()
    }
}

open class ExportersGet : CliktCommand(
    name = "get",
    help = "Get a specific exporter by name for a stream"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exporter")
    internal val exporterName by option(
        "-en",
        "--exporter-name",
        help = "Name of the exporter that is linked to this stream."
    ).required()

    override fun run() {
        "/v1/exporters/$streamName/$exporterName"
            .httpGet()
            .printResponse()
    }
}

open class ExportersDelete : CliktCommand(
    name = "delete",
    help = "Delete a specific exporter by name for a stream"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    internal val exporterName by option(
        "-en",
        "--exporter-name",
        help = "Name of the exporter that is linked to this stream."
    ).required()

    override fun run() {
        "/v1/exporters/$streamName/$exporterName"
            .httpDelete()
            .printResponse()
    }
}

open class ExportersCreate : CliktCommand(
    name = "create",
    help = "Create a new exporter"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")

    internal val exporterName by option(
        "-en",
        "--exporter-name",
        help = "Name of the exporter."
    ).required()

    internal val sinkName by option(
        "-sn",
        "--sink-name",
        help = "The name of the sink that should be used for this exporter."
    ).required()

    internal val sinkType by option(
        "-st",
        "--sink-type",
        help = "The type of sink that should be used for this exporter."
    ).choice(
        *SinkType.getTypes(),
        ignoreCase = true
    ).enum<SinkType>(ignoreCase = true).required()

    internal val intervalSecs by option(
        "-i",
        "--interval",
        help = "The interval in seconds between each batch that is exported to the configured sink."
    ).int()
        .restrictTo(30..3600)
        .required()

    internal val pathPrefix by option(
        "-pp",
        "--path-prefix",
        help = "Optional path prefix. Every object that is exported to the configured sink will have this path prepended to the resource destination."
    )

    override fun run() {
        "/v1/exporters/$streamName"
            .httpPut()
            .jsonBody(
                ExporterCreateRequest(
                    exporterName,
                    sinkName,
                    sinkType,
                    intervalSecs,
                    ExporterType.BATCH,
                    pathPrefix,
                    null
                )
            )
            .printResponse()
    }
}
