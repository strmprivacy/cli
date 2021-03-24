package io.streammachine.api.cli.commands

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.kittinunf.fuel.httpDelete
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.fuel.httpPost
import io.streammachine.api.cli.common.printResponse

class KafkaExportersUsers : CliktCommand(
    name = "kafka-exporters-users",
    help = "View users on kafka-exporters and their details."
) {
    init {
        subcommands(
            kafkaExportersUsersList(),
            KafkaExportersUsersGetOne(),
            KafkaExportersUsersCreate(),
            KafkaExportersUsersDelete()
        )
    }

    override fun run() = Unit
}

open class kafkaExportersUsersList : CliktCommand(
    name = "list",
    help = "get all users on a certain stream"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exporter")

    override fun run() {
        "/v1/kafka-exporters-users/$streamName"
            .httpGet()
            .printResponse()
    }
}
open class KafkaExportersUsersGetOne : CliktCommand(
    name = "get1",
    help = "Get a specific exporter by name for a stream"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exporter")
    internal val clientId by argument("client-id", help = "client Id")

    override fun run() {
        "/v1/kafka-exporters-users/$streamName/$clientId"
            .httpGet()
            .printResponse()
    }
}

open class KafkaExportersUsersDelete : CliktCommand(
    name = "delete",
    help = "Delete a specific exporter by name for a stream"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    internal val clientId by argument("client-id", help = "client Id")
    override fun run() {
        "/v1/kafka-exporters-users/$streamName/$clientId"
            .httpDelete()
            .printResponse()
    }
}

open class KafkaExportersUsersCreate : CliktCommand(
    name = "create",
    help = "Create a new user on a kafka exporter"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    override fun run() {
        "/v1/kafka-exporters-users/$streamName"
            .httpPost()
            .printResponse()
    }
}