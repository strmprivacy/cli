package io.streammachine.api.cli.common

class UnauthorizedException : IllegalStateException("No valid credentials")
class CredentialsExpiredException : IllegalStateException("Credentials have expired")
class UnauthorizedRequestAbortedException :
    IllegalStateException("Request has been aborted, due to missing credentials")
