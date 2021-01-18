package io.streammachine.api.cli.commands

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.required
import com.github.ajalt.clikt.parameters.types.int
import com.github.kittinunf.fuel.gson.jsonBody
import com.github.kittinunf.fuel.httpDelete
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.fuel.httpPut
import io.streammachine.api.cli.common.ConsentLevelMappingCreateRequest
import io.streammachine.api.cli.common.printResponse

class ConsentLevels : CliktCommand(
    name = "consent-levels",
    help = "View consent levels and their mappings."
) {
    init {
        subcommands(
            ConsentLevelsList(),
            ConsentLevelsCreate(),
            ConsentLevelsDelete()
        )
    }

    override fun run() = Unit
}

class ConsentLevelsList : CliktCommand(
    name = "list",
    help = "List your consent levels (optionally filter on a consent level by id)"
) {
    private val id by option("-i", "--id", help = "Optional filter on a consent level.").int()

    override fun run() {
        id?.let { getConsentLevel(it) } ?: getAll()
    }

    private fun getConsentLevel(levelId: Int) = "/v1/consent-levels/$levelId"
        .httpGet()
        .printResponse()

    private fun getAll() = "/v1/consent-levels"
        .httpGet()
        .printResponse()
}

class ConsentLevelsCreate : CliktCommand(
    name = "create",
    help = "Create a new consent level mapping"
) {
    private val id by argument("id", help = "The id of the consent level for which a mapping is created.").int()
    private val name by option(
        "-n", "--name",
        help = "The name / alias of the consent level for which a mapping is created."
    ).required()

    override fun run() {
        "/v1/consent-levels"
            .httpPut()
            .jsonBody(ConsentLevelMappingCreateRequest(id, name))
            .printResponse()
    }
}

class ConsentLevelsDelete : CliktCommand(
    name = "delete",
    help = "Create a consent level mapping"
) {
    private val id by argument("id", help = "The id of the consent level for which a mapping is deleted.").int()

    override fun run() {
        "/v1/consent-levels/$id"
            .httpDelete()
            .printResponse()
    }
}
