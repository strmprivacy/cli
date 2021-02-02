package io.streammachine.api.cli.commands

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.ajalt.clikt.parameters.groups.OptionGroup
import com.github.ajalt.clikt.parameters.groups.cooccurring
import com.github.ajalt.clikt.parameters.groups.mutuallyExclusiveOptions
import com.github.ajalt.clikt.parameters.options.convert
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.required
import com.github.ajalt.clikt.parameters.types.choice
import com.github.ajalt.clikt.parameters.types.enum
import com.github.ajalt.clikt.parameters.types.file
import com.github.kittinunf.fuel.gson.jsonBody
import com.github.kittinunf.fuel.httpDelete
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.fuel.httpPut
import io.streammachine.api.cli.common.Common.readJsonAsString
import io.streammachine.api.cli.common.SinkCreateRequest
import io.streammachine.api.cli.common.SinkType
import io.streammachine.api.cli.common.printResponse

class Sinks : CliktCommand(
    name = "sinks",
    help = "View sinks and their details."
) {
    init {
        subcommands(
            SinksList(),
            SinksCreate(),
            SinksDelete()
        )
    }

    override fun run() = Unit
}

class SinksList : CliktCommand(
    name = "list",
    help = "List your sinks (optionally filter on a sink by name)"
) {
    class NameAndType : OptionGroup() {
        val name by option("-n", "--name", help = "Name of the sink to list.").required()
        val sinkType by option("-st", "--sink-type")
            // TODO verify sink type uppercase
            .choice(*SinkType.getTypes())
            .enum<SinkType>(ignoreCase = true)
            .required()
    }

    private val nameAndType by NameAndType().cooccurring()

    override fun run() {
        nameAndType?.let { getSink(it.name, it.sinkType) } ?: getAll()
    }

    private fun getSink(sink: String, sinkType: SinkType) = "/v1/sinks/$sink"
        .httpGet(listOf("sinkType" to sinkType))
        .printResponse()

    private fun getAll() = "/v1/sinks"
        .httpGet()
        .printResponse()
}

class SinksCreate : CliktCommand(
    name = "create",
    help = "Create a new sink"
) {
    private val sinkType by argument(
        "sink-type",
        help = "The type of the sink that is created."
    ).choice(*SinkType.getTypes()).enum<SinkType>(ignoreCase = true)
    private val sinkName by argument("sink-name", help = "The name of the sink that is created.")
    private val bucketName by option(
        "-bn",
        "--bucket-name",
        help = "The name of the bucket (without scheme)"
    ).required()

    private val credentials by mutuallyExclusiveOptions(
        option(
            "-c",
            "--credentials",
            help = "A valid JSON credentials body, when using ${SinkType.S3.name}. Leave empty when using ${SinkType.GCLOUD.name}."
        ),
        option(
            "-cf",
            "--credentials-file",
            help = "A file that contains a valid JSON credentials body, when using ${SinkType.S3.name}. Leave empty when using ${SinkType.GCLOUD.name}."
        ).file().convert { it.readJsonAsString() },
        help = "Either provide credentials directly, or through a file. If both are provided, file is leading."
    )

    override fun run() {
        "/v1/sinks"
            .httpPut()
            .jsonBody(
                SinkCreateRequest(
                    sinkType,
                    sinkName,
                    bucketName,
                    credentials
                )
            )
            .printResponse()
    }
}

class SinksDelete : CliktCommand(
    name = "delete",
    help = "Delete an existing sink"
) {
    private val sinkName by argument("sink-name", help = "The name of the sink that is created.")
    private val sinkType by argument(
        "sink-type",
        help = "The type of the sink that is created."
    ).choice(*SinkType.getTypes()).enum<SinkType>(ignoreCase = true)

    override fun run() {
        "/v1/sinks/$sinkName"
            .httpDelete(listOf("sinkType" to sinkType))
            .printResponse()
    }
}
