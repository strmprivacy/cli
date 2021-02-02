package io.streammachine.api.cli.commands

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.ajalt.clikt.parameters.options.multiple
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
import io.streammachine.api.cli.common.*

class Streams : CliktCommand(
    name = "streams",
    help = "View streams and their details."
) {
    init {
        subcommands(
            StreamsList(),
            StreamsCreate(),
            StreamsDelete(),
            StreamExporters()
        )
    }

    override fun run() = Unit
}

class StreamsList : CliktCommand(
    name = "list",
    help = "List your streams (optionally filter on a stream by name)"
) {
    private val name by option("-n", "--name", help = "Stream name; provide a name to filter on")

    override fun run() {
        name?.let { getStream(it) } ?: getAll()
    }

    private fun getStream(name: String) = "/v1/streams/$name"
        .httpGet()
        .printResponse()

    private fun getAll() = "/v1/streams"
        .httpGet()
        .printResponse()
}

class StreamsCreate : CliktCommand(
    name = "create",
    help = "Create a new stream"
) {
    private val name by argument("name", help = "The name of the stream that should be created.")
    private val description by option(
        "-d",
        "--description",
        help = "Optional description to describe the purpose of this stream"
    )
    private val tags by option("-t", "--tag", help = "Optional tag for this stream.").multiple()

    override fun run() {
        "/v1/streams"
            .httpPut()
            .jsonBody(StreamCreateRequest(name, description, tags))
            .printResponse()
    }
}

class StreamsDelete : CliktCommand(
    name = "delete",
    help = "Delete an existing stream"
) {
    private val name by argument("name", help = "The name of the stream that should be deleted.")

    override fun run() {
        "/v1/streams/$name"
            .httpDelete()
            .printResponse("Stream '$name' has been deleted.")
    }
}

class StreamExporters : CliktCommand(
    name = "exporters",
    help = "View exporters and their details."
) {
    init {
        subcommands(
            StreamExportersGet(),
            StreamExportersDelete()
        )
    }

    override fun run() = Unit
}

open class StreamExportersGet : CliktCommand(
    name = "get",
    help = "Get a specific exporter by name for a stream"
) {
    internal val name by argument("name", help = "Name of the exporter that should be listed")
    internal val streamName by option(
        "-sn",
        "--stream-name",
        help = "Name of the stream that is linked to this exporter."
    ).required()

    override fun run() {
        "/v1/streams/$streamName/exporters/$name"
            .httpGet()
            .printResponse()
    }
}

open class StreamExportersCreate : CliktCommand(
    name = "create",
    help = "Create a new exporter"
) {
    internal val name by argument("name", help = "The name of the exporter that should be created.")

    internal val streamName by option(
        "-n",
        "--stream-name",
        help = "Name of the stream that is linked to this exporter."
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
        "/v1/streams/$streamName/exporters"
            .httpPut()
            .jsonBody(
                ExporterCreateRequest(
                    name,
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

open class StreamExportersDelete : CliktCommand(
    name = "delete",
    help = "Delete an existing exporter"
) {
    internal val name by argument("name", help = "The name of the exporter that should be deleted.")
    internal val streamName by option(
        "-sn",
        "--stream-name",
        help = "Name of the stream that is linked to this exporter."
    ).required()

    override fun run() {
        "/v1/streams/$streamName/exporters/$name"
            .httpDelete()
            .printResponse()
    }
}
