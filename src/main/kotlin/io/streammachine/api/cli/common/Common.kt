package io.streammachine.api.cli.common

import com.github.ajalt.clikt.output.CliktHelpFormatter
import com.github.ajalt.clikt.output.HelpFormatter
import com.github.ajalt.mordant.TermColors
import com.github.kittinunf.fuel.core.Request
import com.github.kittinunf.fuel.core.Response
import com.github.kittinunf.result.Result
import com.google.gson.GsonBuilder
import com.google.gson.JsonObject
import com.google.gson.JsonParser
import io.streammachine.api.cli.commands.AuthResponse
import io.streammachine.api.cli.commands.Login
import io.streammachine.api.cli.common.Common.GSON
import io.streammachine.api.cli.common.Common.Terminal.ORANGE
import io.streammachine.api.cli.common.Common.Terminal.TERM_COLORS
import java.io.File
import java.io.FileReader
import java.io.FileWriter
import java.nio.charset.StandardCharsets

object Common {
    object Terminal {
        internal val TERM_COLORS = TermColors()

        internal val ORANGE = TERM_COLORS.rgb("f25c03")
        internal val BLUE = TERM_COLORS.rgb("0004bf")
    }

    internal val GSON = GsonBuilder()
        .setPrettyPrinting()
        .create()

    internal var VERBOSE_LOGGING = false
    internal val VERSION = Common.javaClass.`package`.implementationVersion ?: "snapshot"

    private var CREDENTIALS: AuthResponse? = null

    internal fun getApiUrl() = System.getenv("STRM_API_HOST")?.let { "https://$it" } ?: "https://api.streammachine.io"

    internal fun getCredentialsPath() = File("${System.getProperty("user.home")}/.config/stream-machine")
    internal fun getCredentialsFile() = getCredentialsPath().resolve("credentials.json")

    internal fun writer(block: (FileWriter) -> Unit) {
        val fileWriter = FileWriter(getCredentialsFile(), StandardCharsets.UTF_8, false)
        block.invoke(fileWriter)
        fileWriter.flush()
        fileWriter.close()
    }

    internal fun getReader() = FileReader(getCredentialsFile())

    internal fun getCredentials(): AuthResponse? {
        if (CREDENTIALS == null) {
            CREDENTIALS = runCatching { GSON.fromJson(getReader(), AuthResponse::class.java) }.getOrNull()
        }

        return CREDENTIALS
    }

    internal fun getCredentialsAsJson() = runCatching { GSON.fromJson(getReader(), JsonObject::class.java) }.getOrNull()

    internal fun setCredentials(credentials: AuthResponse?) {
        CREDENTIALS = credentials
    }

    internal fun JsonObject.prettyPrint() = toString().asPrettyJson()

    internal fun File.readJsonAsString() =
        runCatching { GSON.fromJson(FileReader(this), JsonObject::class.java).toString() }
            .getOrElse {
                throw IllegalArgumentException(
                    it.cause?.message?.let { message ->
                        "Unable to read file '${this.absolutePath}', cause = $message"
                    } ?: "Unable to read file '${this.absolutePath}'"
                )
            }
}

fun Request.printResponse(noContentMessage: String? = null): Response = response { _, response, result ->
    when (result) {
        is Result.Failure -> {
            if (result.getException().exception !is UnauthorizedRequestAbortedException && response.statusCode != 401) {
                println(response.body().asString("application/json").asPrettyJson())
            } else {
                with(TERM_COLORS) {
                    println(
                        "Unauthorized request, access token might be outdated or invalid. Please re-login using: ${
                            (bold)(
                                Login.FULL_COMMAND
                            )
                        }"
                    )
                }
            }
            printVerbose(result.getException().exception.message)
        }
        is Result.Success -> {
            if (response.statusCode != 204) {
                println(response.body().asString("application/json").asPrettyJson())
            } else {
                noContentMessage?.let { println(it) } ?: println("Request has been succesfully processed.")
            }
        }
    }
}.get()

fun String.asPrettyJson(): String = GSON.toJson(JsonParser.parseString(this))

fun printVerbose(vararg messages: Any?) {
    if (Common.VERBOSE_LOGGING) {
        messages.forEach { println(it) }
    }
}

object ColorHelpFormatter : CliktHelpFormatter() {
    override fun renderTag(tag: String, value: String) = TERM_COLORS.green(super.renderTag(tag, value))
    override fun renderOptionName(name: String) = super.renderOptionName(name)
    override fun renderArgumentName(name: String) = ORANGE(super.renderArgumentName(name))
    override fun renderSubcommandName(name: String) = ORANGE(super.renderSubcommandName(name))
    override fun renderSectionTitle(title: String) =
        (TERM_COLORS.bold + TERM_COLORS.underline)(super.renderSectionTitle(title))

    override fun optionMetavar(option: HelpFormatter.ParameterHelp.Option) =
        TERM_COLORS.green(super.optionMetavar(option))
}
