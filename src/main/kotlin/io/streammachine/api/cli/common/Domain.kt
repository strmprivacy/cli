package io.streammachine.api.cli.common

enum class ConsentLevelType {
    GRANULAR,
    CUMULATIVE;

    companion object {
        fun getTypes() = values().map { it.name.toLowerCase() }.toTypedArray()
    }
}

data class StreamCreateRequest(
    val name: String,
    val description: String?,
    val tags: List<String>?
)

data class OutputStreamCreateRequest(
    val name: String,
    val description: String?,
    val tags: List<String>?,
    val consentLevelType: ConsentLevelType? = ConsentLevelType.CUMULATIVE,
    val consentLevels: List<Int> = emptyList(),
    val decrypterName: String? = null
)

data class SinkCreateRequest(
    val sinkType: SinkType,
    val sinkName: String,
    val bucketName: String,
    val credentials: String? = null
)

data class ExporterCreateRequest(
    val name: String,
    val sinkName: String,
    val sinkType: SinkType,
    val intervalSecs: Int,
    val type: ExporterType = ExporterType.BATCH,
    val pathPrefix: String?,
    val extraConfig: String?,
    val exportKeys: Boolean
)

enum class ExporterType {
    BATCH
}

enum class SinkType {
    GCLOUD,
    S3;

    companion object {
        fun getTypes() = values().map { it.name.toLowerCase() }.toTypedArray()
    }
}

data class ConsentLevelMappingCreateRequest(
    val id: Int,
    val name: String
)