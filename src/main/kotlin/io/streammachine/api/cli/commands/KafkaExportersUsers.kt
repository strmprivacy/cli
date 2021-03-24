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
    help = "Get all users on a certain Kafka exporter"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exporter")

    override fun run() {
        "/v1/kafka-exporters-users/$streamName".httpGet().printResponse()
    }
}
open class KafkaExportersUsersGetOne : CliktCommand(
    name = "get",
    help = "Get a specific user on a Kafka exporter by stream name"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    internal val clientId by argument("client-id", help = "The client id that identifies the Kafka Consumer")

    override fun run() {
        "/v1/kafka-exporters-users/$streamName/$clientId".httpGet().printResponse()
    }
}

open class KafkaExportersUsersDelete : CliktCommand(
    name = "delete",
    help = "Delete a specific exporter by name for a stream"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    internal val clientId by argument("client-id", help = "The client id that identifies the Kafka Consumer")
    override fun run() {
        "/v1/kafka-exporters-users/$streamName/$clientId".httpDelete().printResponse()
    }
}

open class KafkaExportersUsersCreate : CliktCommand(
    name = "create",
    help = "Create a new user on an existing Kafka exporter"
) {
    internal val streamName by argument("stream-name", help = "Name of the stream that is being exported")
    override fun run() {
        "/v1/kafka-exporters-users/$streamName".httpPost().printResponse()
    }
}