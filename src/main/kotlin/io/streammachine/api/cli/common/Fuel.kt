package io.streammachine.api.cli.common

import com.github.kittinunf.fuel.core.*
import com.github.kittinunf.fuel.core.extensions.authentication
import com.github.kittinunf.fuel.core.interceptors.LogRequestInterceptor
import io.streammachine.api.cli.commands.Login
import java.io.IOException
import java.io.InputStream
import java.net.HttpURLConnection

fun initializeFuel() {
    FuelManager.instance.basePath = Common.getApiUrl()
    FuelManager.instance.apply {
        hook = AuthenticationRequiredHook()
        addRequestInterceptor(AuthenticateRequestInterceptor())

        if (Common.VERBOSE_LOGGING) {
            addRequestInterceptor(LogRequestInterceptor)
        }
    }
}

class AuthenticationRequiredHook : Client.Hook {
    override fun preConnect(connection: HttpURLConnection, request: Request) {
        runCatching {
            if (!Login.LOGIN_PATHS.contains(request.url.path)) {
                Common.getCredentials()?.let { originalCredentials ->
                    if (System.currentTimeMillis().div(1000) >= originalCredentials.expiresAt) {
                        Login.refreshCredentials(originalCredentials.refreshToken) {
                            printVerbose("Refreshed credentials using refresh token. Verbose logging header with Authorization Bearer token outdated.")
                            request.authentication().bearer(it.idToken)
                        }
                    }
                } ?: throw UnauthorizedException()
            }
        }.getOrElse {
            with(Common.Terminal.TERM_COLORS) {
                when {
                    it is UnauthorizedException || it.cause is UnauthorizedException -> println(
                        "Currently there are ${(underline)("no credentials")} active. Please re-login using: ${
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
    }

    override fun interpretResponseStream(request: Request, inputStream: InputStream?): InputStream? = inputStream

    override fun postConnect(request: Request) {
        // no-op
    }

    override fun httpExchangeFailed(request: Request, exception: IOException) {
        // no-op
    }
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
