package io.streammachine.api.cli.common

import com.github.kittinunf.fuel.gson.responseObject
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.result.Result
import com.google.gson.FieldNamingPolicy
import com.google.gson.GsonBuilder

object UpdateCheck {
    private const val REPOSITORY = "streammachineio/cli"
    private val GSON = GsonBuilder()
        .setFieldNamingPolicy(FieldNamingPolicy.LOWER_CASE_WITH_UNDERSCORES)
        .create()

    fun printUpdateMessageIfAvailable() {
        latestVersion()?.let {
            with(Common.Terminal.TERM_COLORS) {
                println(
                    "Version ${(bold)(it)} of the Stream Machine CLI is available. Get it at ${
                        (bold)("https://github.com/streammachineio/cli/releases/tag/v$it")
                    }"
                )
                println()
            }
        }
    }

    private fun latestVersion(): String? {
        return getLatestRelease()?.let {
            val version = it.tagName.replace("v", "")

            if (version != Common.VERSION) {
                version
            } else {
                null
            }
        }
    }

    private fun getLatestRelease() = runCatching {
        "https://api.github.com/repos/$REPOSITORY/releases/latest"
            .httpGet()
            .responseObject<GitHubRelease>(GSON)
    }.mapCatching {
        if (it.third is Result.Success) it.third.get() else null
    }.getOrNull()

}

internal data class GitHubRelease(
    val tagName: String
)
