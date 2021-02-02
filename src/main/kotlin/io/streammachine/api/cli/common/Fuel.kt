package io.streammachine.api.cli.common

import com.github.kittinunf.fuel.core.FoldableRequestInterceptor
import com.github.kittinunf.fuel.core.FuelManager
import com.github.kittinunf.fuel.core.Request
import com.github.kittinunf.fuel.core.RequestTransformer
import com.github.kittinunf.fuel.core.extensions.authentication
import com.github.kittinunf.fuel.core.interceptors.LogRequestInterceptor
import io.streammachine.api.cli.commands.Login

fun initializeFuel() {
    FuelManager.instance.basePath = Common.getApiUrl()
    FuelManager.instance.apply {
        addRequestInterceptor(AuthenticateRequestInterceptor())

        if (Common.VERBOSE_LOGGING) {
            addRequestInterceptor(LogRequestInterceptor)
        }
    }
}

fun Request.authenticated(): Request {
    runCatching {
        if (!Login.LOGIN_PATHS.contains(url.path)) {
            Common.getCredentials()?.let { originalCredentials ->
                if (System.currentTimeMillis().div(1000) >= originalCredentials.expiresAt) {
                    Login.refreshCredentials(originalCredentials.refreshToken) {
                        printVerbose("Refreshed credentials using refresh token.")
                    }
                }
            } ?: throw UnauthorizedException()
        }
    }.onFailure {
        with(Common.Terminal.TERM_COLORS) {
            when {
                it is UnauthorizedException || it.cause is UnauthorizedException -> println(
                    "Currently there are ${(underline)("no credentials")} active. Please login using: ${
                        (bold)(
                            Login.FULL_COMMAND
                        )
                    }"
                )
                it is CredentialsExpiredException || it.cause is CredentialsExpiredException -> println(
                    "Currently there are ${(underline)("no credentials")} active that could be refreshed. Please re-login using: ${
                        (bold)(
                            Login.FULL_COMMAND
                        )
                    }"
                )
                else -> println(
                    "Unfortunately something went wrong while refreshing your login credentials. Please re-login using: ${
                        (bold)(
                            Login.FULL_COMMAND
                        )
                    }"
                )
            }
        }

        throw UnauthorizedRequestAbortedException()
    }

    return this
}

class AuthenticateRequestInterceptor : FoldableRequestInterceptor {
    override fun invoke(next: RequestTransformer): RequestTransformer = { request ->
        if (!Login.LOGIN_PATHS.contains(request.url.path)) {
            Common.getCredentials()?.let {
                request.authentication().bearer(it.idToken)
            }
        }

        next(request)
    }
}
