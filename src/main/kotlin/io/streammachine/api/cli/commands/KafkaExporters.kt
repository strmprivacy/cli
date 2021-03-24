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

class KafkaExporters : CliktCommand(
    name = "kafka-exporters",
    help = "Kafka exporters and details."
) {
    init {
        subcommands(
            KafkaExportersGet(),
            KafkaExportersCreate(),
            KafkaExportersDelete()
        )
    }

    override fun run() = Unit
}

open class KafkaExportersGet : CliktCommand(
    name = "get",
    help = "Get a specific exporter by name for a stream"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exporter")

    override fun run() {
        "/v1/kafka-exporters/$streamName".httpGet().printResponse()
    }
}

open class KafkaExportersDelete : CliktCommand(
    name = "delete",
    help = "Delete a Kafka exporter by stream name"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    override fun run() {
        "/v1/kafka-exporters/$streamName".httpDelete().printResponse()
    }
}

open class KafkaExportersCreate : CliktCommand(
    name = "create",
    help = "Create a new Kafka exporter"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    override fun run() {
        "/v1/kafka-exporters/$streamName".httpPut().printResponse()
    }
}