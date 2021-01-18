package io.streammachine.api.cli.commands

import com.github.kittinunf.fuel.core.FuelManager
import io.mockk.*
import io.streammachine.api.cli.common.Common
import io.streammachine.api.cli.common.CredentialsExpiredException
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows
import org.mockserver.integration.ClientAndServer
import org.mockserver.model.HttpRequest
import org.mockserver.model.HttpResponse
import java.io.File
import java.util.concurrent.ExecutionException

internal class LoginTest {
    private lateinit var mock: ClientAndServer

    @BeforeEach
    fun setUpEach() {
        mockkObject(Common)
        mock = ClientAndServer.startClientAndServer()
        FuelManager.instance.basePath = "http://localhost:${mock.localPort}"
    }

    @AfterEach
    fun tearDownEach() {
        unmockkObject(Common)
        Common.setCredentials(null)
    }

    @Test
    fun `credentials refresh - success`() {
        // Given a refreshToken
        val refreshToken = "a-refresh-token"
        mock.`when`(
            HttpRequest.request()
                .withBody("""{"refreshToken":"$refreshToken"}""")
                .withPath("/v1/refresh")
        ).respond(
            HttpResponse.response()
                .withBody("auth/dummy-credentials.json".getResourceStream())
                .withStatusCode(200)
        )

        val block = mockk<(AuthResponse) -> Unit>()
        val capturedAuthResponse = slot<AuthResponse>()
        justRun { block.invoke(capture(capturedAuthResponse)) }

        // When credentials are refreshed
        Login.refreshCredentials(refreshToken, block)

        // Then we should have stored new credentials
        verify(exactly = 1) { block.invoke(capturedAuthResponse.captured) }

        assertThat(capturedAuthResponse.captured).isEqualTo(
            AuthResponse(
                "my-email@email.com",
                "some-billing-id",
                42,
                "some-arbitrary-token",
                "an-arbitrary-refresh-token"
            )
        )
    }

    @Test
    fun `credentials refresh - failure`() {
        // Given a refreshToken
        val refreshToken = "bananas"
        mock.`when`(
            HttpRequest.request()
                .withBody("""{"refreshToken":"$refreshToken"}""")
                .withPath("/v1/refresh")
        ).respond(
            HttpResponse.response()
                .withBody("auth/error-refresh-response.json".getResourceStream())
                .withStatusCode(400)
        )

        val block = mockk<(AuthResponse) -> Unit>()
        justRun { block.invoke(any()) }

        // When credentials are refreshed
        val exception = assertThrows<ExecutionException> {
            Login.refreshCredentials(refreshToken, block)
            // Then we should have no interactions with the provided block
            verify(exactly = 0) { block.invoke(any()) }
        }

        assertThat(exception.cause).isOfAnyClassIn(CredentialsExpiredException::class.java)
    }

    @Test
    fun `authenticate - success`() {
        // Given an email address and password
        val email = "my@email.stream"
        val password = "my-password"

        mock.`when`(
            HttpRequest.request()
                .withMethod("POST")
                .withBody("""{"email":"$email","password":"$password"}""")
                .withPath("/v1/auth")
        ).respond(
            HttpResponse.response()
                .withBody("auth/dummy-credentials.json".getResourceStream())
                .withStatusCode(200)
        )

        Common.apply {
            val tempFile = File.createTempFile("strm-cli", "auth-success")
            every { getCredentialsPath() } returns tempFile.parentFile
            every { getCredentialsFile() } returns tempFile
        }

        // When credentials are requested
        Login().initialize(email, password)

        // Then we should have stored new credentials
        assertThat(Common.getCredentials()).isEqualTo(
            AuthResponse(
                "my-email@email.com",
                "some-billing-id",
                42,
                "some-arbitrary-token",
                "an-arbitrary-refresh-token"
            )
        )
    }

    @Test
    fun `authenticate - failure`() {
        // Given an email address and password
        val email = "my@email.stream"
        val password = "my-password"

        mock.`when`(
            HttpRequest.request()
                .withBody("""{"email":"$email","password":"$password"}""")
                .withMethod("POST")
                .withPath("/v1/auth")
        ).respond(
            HttpResponse.response()
                .withBody("auth/error-refresh-response.json".getResourceStream())
                .withStatusCode(400)
        )

        Common.apply {
            val tempFile = File.createTempFile("strm-cli", "auth-failure")
            every { getCredentialsPath() } returns tempFile.parentFile
            every { getCredentialsFile() } returns tempFile
        }

        // When credentials are requested
        Login().initialize(email, password)

        // Then we should not have stored new credentials
        assertThat(Common.getCredentials()).isNull()
    }

    private fun Login.initialize(email: String, password: String): Login {
        this.parse(
            listOf(
                "--email", email,
                "--password", password
            )
        )
        return this
    }

    private fun String.getResourceStream() =
        LoginTest::class.java.classLoader.getResourceAsStream(this).readAllBytes()
}
