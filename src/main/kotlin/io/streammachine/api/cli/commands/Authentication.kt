package io.streammachine.api.cli.commands

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.options.flag
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.prompt
import com.github.kittinunf.fuel.gson.jsonBody
import com.github.kittinunf.fuel.gson.responseObject
import com.github.kittinunf.fuel.httpPost
import com.github.kittinunf.result.Result.Failure
import com.github.kittinunf.result.Result.Success
import io.streammachine.api.cli.Strm
import io.streammachine.api.cli.common.Common
import io.streammachine.api.cli.common.Common.GSON
import io.streammachine.api.cli.common.Common.getCredentialsFile
import io.streammachine.api.cli.common.Common.getCredentialsPath
import io.streammachine.api.cli.common.Common.prettyPrint
import io.streammachine.api.cli.common.Common.writer
import io.streammachine.api.cli.common.CredentialsExpiredException
import io.streammachine.api.cli.common.asPrettyJson
import io.streammachine.api.cli.common.printVerbose
import java.time.LocalDateTime
import java.time.ZoneOffset

class Authentication : CliktCommand(
    name = COMMAND,
    help = "Authenticate against Stream Machine.",
    printHelpOnEmptyArgs = true
) {
    companion object {
        internal const val COMMAND = "auth"
    }

    init {
        subcommands(
            Login(),
            Revoke(),
            Show(),
            Token()
        )
    }

    override fun run() = Unit
}

class Login : CliktCommand(
    name = COMMAND,
    help = "Login with your Stream Machine Portal credentials."
) {
    private val email by option("-e", "--email", help = "The email address of the user you want to login as.").prompt(
        requireConfirmation = false,
        hideInput = false,
        default = Common.getCredentials()?.email,
        text = "Enter your Stream Machine portal email address",
    )
    private val password by option(
        "-p", "--password", help = "Password for this login"
    ).prompt(
        requireConfirmation = false,
        hideInput = true,
        text = "Please enter your Stream Machine portal password"
    )

    companion object {
        internal const val COMMAND = "login"
        internal const val FULL_COMMAND = "${Strm.COMMAND} ${Authentication.COMMAND} $COMMAND"

        internal const val AUTH_PATH = "/v1/auth"
        internal const val REFRESH_PATH = "/v1/refresh"
        internal val LOGIN_PATHS = listOf(AUTH_PATH, REFRESH_PATH)

        internal fun refreshCredentials(refreshToken: String, block: (AuthResponse) -> Unit) {
            REFRESH_PATH.httpPost()
                .jsonBody(RefreshRequest(refreshToken))
                .responseObject<AuthResponse> { _, _, result ->
                    when (result) {
                        is Failure -> throw CredentialsExpiredException()
                        is Success -> {
                            storeCredentials(null, result.get())
                            block(result.get())
                        }
                    }
                }
                .get()
        }

        private fun storeCredentials(email: String?, result: AuthResponse) =
            runCatching {
                getCredentialsPath().mkdirs()
                writer { GSON.toJson(result, it) }
                Common.setCredentials(result)
            }.fold({
                email?.let { println("Succesfully logged in as $it") }
            }) {
                email?.let { _ -> println("Failed to store credentials for $it at ${getCredentialsFile()}.") }
            }
    }

    override fun run() {
        AUTH_PATH.httpPost()
            .jsonBody(AuthRequest(email, password))
            .responseObject<AuthResponse> { _, response, result ->
                when (result) {
                    is Failure -> println(response.body().asString("application/json").asPrettyJson())
                    is Success -> storeCredentials(email, result.get())
                }
            }
            .get()
    }
}

class Revoke : CliktCommand(
    name = "revoke",
    help = "Revoke your currently active credentials."
) {
    override fun run() =
        runCatching {
            if (getCredentialsFile().exists()) getCredentialsFile().delete() else false
        }.fold({
            if (it) {
                println("Revoked currently active credentials.")
            } else {
                println("Currently there are no active credentials.")
            }
        }) {
            println("Failed to revoke currently active credentials.")
            println("reason = ${it.message}")
        }
}

class Show : CliktCommand(
    name = "show",
    help = "Show your currently active credentials and details."
) {
    private val asJson by option("-j", "--json", help = "Show the current credentials as JSON").flag(default = false)

    override fun run() {
        val credentials =
            if (asJson) Common.getCredentialsAsJson()?.prettyPrint() else Common.getCredentials()
                ?.toString()

        credentials?.let { println(it) } ?: println("Currently there are no credentials active.")
    }
}

class Token : CliktCommand(
    name = "print-access-token",
    help = "Print your current access token"
) {
    override fun run() {
        Common.getCredentials()?.let {
            if (System.currentTimeMillis().div(1000) >= it.expiresAt) {
                Login.refreshCredentials(it.refreshToken) {
                    printVerbose("Refreshed credentials using refresh token. Verbose logging header with Authorization Bearer token outdated.")
                }
            }

            println(it.idToken)
        } ?: println("Currently there are no credentials active.")
    }
}

data class AuthRequest(val email: String, val password: String)
data class RefreshRequest(val refreshToken: String)
data class AuthResponse(
    val email: String,
    val billingId: String,
    val expiresAt: Long,
    val idToken: String,
    val refreshToken: String
) {
    override fun toString(): String {
        return """
            Credentials for $email
            Billing id = $billingId
            Current token valid until = ${LocalDateTime.ofEpochSecond(expiresAt, 0, ZoneOffset.UTC)} UTC
        """.trimIndent()
    }
}
