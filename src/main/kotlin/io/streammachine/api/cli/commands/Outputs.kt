package io.streammachine.api.cli.commands

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.ajalt.clikt.parameters.options.multiple
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.required
import com.github.ajalt.clikt.parameters.options.unique
import com.github.ajalt.clikt.parameters.types.choice
import com.github.ajalt.clikt.parameters.types.enum
import com.github.ajalt.clikt.parameters.types.int
import com.github.kittinunf.fuel.gson.jsonBody
import com.github.kittinunf.fuel.httpDelete
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.fuel.httpPut
import io.streammachine.api.cli.common.*

class Outputs : CliktCommand(
    name = "outputs",
    help = "View outputs and their details."
) {
    init {
        subcommands(
            OutputsList(),
            OutputsCreate(),
            OutputsDelete(),
            OutputStreamExporters()
        )
    }

    override fun run() = Unit
}

class OutputsList : CliktCommand(
    name = "list",
    help = "List your outputs (optionally filter on an output by name)"
) {
    private val name by argument("name", help = "The name of the stream for which outputs are listed.")
    private val outputName by option("-on", "--output-name", help = "Optional filter on an output name.")

    override fun run() {
        outputName?.let { getOutput(it) } ?: getAll()
    }

    private fun getOutput(output: String) = "/v1/streams/$name/outputs/$output"
        .httpGet()
        .printResponse()

    private fun getAll() = "/v1/streams/$name/outputs"
        .httpGet()
        .printResponse()
}

class OutputsCreate : CliktCommand(
    name = "create",
    help = "Create a new output"
) {
    private val name by argument("name", help = "The name of the stream for which an output is created.")
    private val outputName by argument("output-name", help = "The name of the output that is created.")

    private val consentLevels by option("-cl", "--consent-level", help = "The consent levels for this output.")
        .int().multiple().unique()

    private val consentLevelType by option("-clt", "--consent-level-type").choice(
        *ConsentLevelType.getTypes(),
        ignoreCase = true
    ).enum<ConsentLevelType>(ignoreCase = true).required()

    private val description by option(
        "-d",
        "--description",
        help = "Optional description to describe the purpose of this output."
    )

    private val tags by option("-t", "--tag", help = "Optional tag for this output.").multiple().unique()

    private val decrypterName by option(
        "-dn",
        "--decrypter-name",
        help = "Name for the decrypter that will produce this output."
    )

    override fun run() {
        "/v1/streams/$name/outputs"
            .httpPut()
            .jsonBody(
                OutputStreamCreateRequest(
                    outputName,
                    description,
                    tags.toList(),
                    consentLevelType,
                    consentLevels.toList(),
                    decrypterName
                )
            )
            .printResponse()
    }
}

class OutputsDelete : CliktCommand(
    name = "delete",
    help = "Delete an existing output"
) {
    private val name by argument("name", help = "The name of the stream for which an output is deleted.")
    private val outputName by argument("output-name", help = "The name of the output that is deleted.")

    override fun run() {
        "/v1/streams/$name/outputs/$outputName"
            .httpDelete()
            .printResponse()
    }
}

class OutputStreamExporters : CliktCommand(
    name = "exporters",
    help = "View exporters and their details."
) {
    init {
        subcommands(
            OutputStreamExportersGet(),
            OutputStreamExportersDelete()
        )
    }

    override fun run() = Unit
}

class OutputStreamExportersGet : StreamExportersGet() {
    private val outputStreamName by option(
        "-osn",
        "--output-stream-name",
        help = "Name of the output stream that is linked to this exporter."
    ).required()

    override fun run() {
        "/v1/streams/$streamName/outputs/$outputStreamName/exporters/$name"
            .httpGet()
            .printResponse()
    }
}

class OutputStreamExportersCreate : StreamExportersCreate() {
    private val outputStreamName by option(
        "-osn",
        "--output-stream-name",
        help = "Name of the output stream that is linked to this exporter."
    ).required()

    override fun run() {
        "/v1/streams/$streamName/outputs/$outputStreamName/exporters"
            .httpPut()
            .jsonBody(
                ExporterCreateRequest(
                    name,
                    sinkName,
                    sinkType,
                    intervalSecs,
                    ExporterType.BATCH,
                    pathPrefix,
                    null,
                    false
                )
            )
            .printResponse()
    }
}

class OutputStreamExportersDelete : StreamExportersDelete() {
    private val outputStreamName by option(
        "-osn",
        "--output-stream-name",
        help = "Name of the output stream that is linked to this exporter."
    ).required()

    override fun run() {
        "/v1/streams/$streamName/outputs/$outputStreamName/exporters/$name"
            .httpDelete()
            .printResponse()
    }
}