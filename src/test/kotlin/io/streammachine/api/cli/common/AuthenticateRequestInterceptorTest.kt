package io.streammachine.api.cli.common

import com.github.kittinunf.fuel.core.Headers.Companion.AUTHORIZATION
import com.github.kittinunf.fuel.core.Method
import com.github.kittinunf.fuel.core.requests.DefaultRequest
import io.mockk.every
import io.mockk.mockkObject
import io.mockk.unmockkObject
import io.streammachine.api.cli.commands.AuthResponse
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import java.net.URL

internal class AuthenticateRequestInterceptorTest {
    private val underTest = AuthenticateRequestInterceptor().invoke { it }

    @BeforeEach
    fun setUp() {
        mockkObject(Common)
        every { Common.getCredentials() } returns AuthResponse(
            "my-email@email.com",
            "some-billing-id",
            42,
            "some-arbitrary-token",
            "an-arbitrary-refresh-token"
        )
    }

    @AfterEach
    fun tearDown() {
        unmockkObject(Common)
    }

    @Test
    fun `non-authentication path should get an authorization header`() {
        // Given a request
        val request = DefaultRequest(Method.GET, URL("http://api.stream/an/endpoint"))

        // When it is modified
        val modified = underTest.invoke(request)

        // Then it should contain a bearer token
        assertThat(modified[AUTHORIZATION]).contains("Bearer some-arbitrary-token").hasSize(1)
    }

    @Test
    fun `authentication path should not get an authorization header`() {
        // Given a request
        val request = DefaultRequest(Method.GET, URL("http://api.stream/v1/auth"))

        // When it is modified
        val modified = underTest.invoke(request)

        // Then it should contain a bearer token
        assertThat(modified[AUTHORIZATION]).hasSize(0)
    }
}
