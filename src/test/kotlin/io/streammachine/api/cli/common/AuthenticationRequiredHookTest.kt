package io.streammachine.api.cli.common

import com.github.kittinunf.fuel.core.Request
import com.github.kittinunf.fuel.core.extensions.authentication
import io.mockk.*
import io.streammachine.api.cli.commands.AuthResponse
import io.streammachine.api.cli.commands.Login
import org.junit.jupiter.api.*
import java.net.HttpURLConnection
import java.net.URL

internal class AuthenticationRequiredHookTest {
    private val underTest = AuthenticationRequiredHook()
    private val connection = mockk<HttpURLConnection>()

    @BeforeEach
    fun setUp() {
        mockkObject(Common)
        mockkObject(Login)
    }

    @AfterEach
    fun tearDown() {
        unmockkObject(Common)
        unmockkObject(Login)
    }

    @Test
    fun `should refresh credentials when id token is expired`() {
        // Given a Fuel request and expired credentials
        val request = spyk<Request> {
            every { url }.returns(URL("http://api.stream/an/endpoint"))
        }

        every { Common.getCredentials() } returns AuthResponse(
            "my-email@email.com",
            "some-billing-id",
            42,
            "some-arbitrary-token",
            "an-arbitrary-refresh-token"
        )
        every { Login.refreshCredentials(any(), any()) } answers {
            secondArg<(AuthResponse) -> Unit>().invoke(
                AuthResponse(
                    "my-email@email.com",
                    "some-billing-id",
                    System.currentTimeMillis().div(1000).plus(3600),
                    "a-new-id-token",
                    "a-new-refresh-token"
                )
            )
        }

        // When the preConnect hook runs
        underTest.preConnect(connection, request)

        // Then the credentials should be refreshed and the request header should be set
        verify { Login.refreshCredentials(any(), any()) }
        verify { request.authentication().bearer("a-new-id-token") }
    }

    @Test
    fun `should throw an aborted request exception when no credentials are active`() {
        // Given a Fuel request and no active credentials
        val request = spyk<Request> {
            every { url }.returns(URL("http://api.stream/an/endpoint"))
        }
        every { Common.getCredentials() } returns null

        // When the preConnect hook runs
        // Then an aborted exception should be thrown
        assertThrows<UnauthorizedRequestAbortedException> {
            underTest.preConnect(connection, request)
            verify { Login.refreshCredentials(any(), any()) }
        }

    }

    @Test
    fun `should throw an aborted request exception when refresh fails`() {
        // Given a Fuel request and no active credentials
        val request = spyk<Request> {
            every { url }.returns(URL("http://api.stream/an/endpoint"))
        }
        every { Common.getCredentials() } returns AuthResponse(
            "my-email@email.com",
            "some-billing-id",
            42,
            "some-arbitrary-token",
            "an-arbitrary-refresh-token"
        )
        every { Login.refreshCredentials(any(), any()) } throws CredentialsExpiredException()

        // When the preConnect hook runs
        // Then an aborted exception should be thrown
        assertThrows<UnauthorizedRequestAbortedException> {
            underTest.preConnect(connection, request)
            verify { Login.refreshCredentials(any(), any()) }
        }
    }

    @Test
    fun `should allow a call to authentication or refresh endpoints`() {
        // Given a Fuel request and already active credentials
        val request = spyk<Request> {
            every { url }.returns(URL("http://api.stream/v1/auth"))
        }
        every { Common.getCredentials() } returns AuthResponse(
            "my-email@email.com",
            "some-billing-id",
            42,
            "some-arbitrary-token",
            "an-arbitrary-refresh-token"
        )

        // When the preConnect hook runs
        assertDoesNotThrow {
            underTest.preConnect(connection, request)

            // And no calls are done to refresh any existing credentials
            verify(exactly = 0) { Login.refreshCredentials(any(), any()) }
        }
    }
}
