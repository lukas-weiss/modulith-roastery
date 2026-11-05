package dev.lweiss.northbean.eventlog

import dev.lweiss.northbean.NorthbeanApplication
import org.springframework.jdbc.core.simple.JdbcClient
import org.springframework.modulith.events.FailedEventPublications
import org.springframework.modulith.events.ResubmissionOptions
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RestController
import tools.jackson.databind.JsonNode
import tools.jackson.databind.ObjectMapper
import java.time.OffsetDateTime
import java.util.UUID

/**
 * Module "eventlog": shows what Spring Modulith stored in the EVENT_PUBLICATION table.
 *
 * Each published event gets one row per listener. The row is marked COMPLETED when that listener
 * has finished, or FAILED if it threw an exception. Failed rows stay and can be resubmitted.
 * This module only reads the table and depends on no other module.
 */
@RestController
class EventLogController(
    private val jdbc: JdbcClient,
    private val json: ObjectMapper,
    private val failed: FailedEventPublications,
) {

    data class Publication(
        val id: UUID,
        val eventType: String,
        val module: String,
        val listener: String,
        val event: JsonNode,
        val status: String,
        val publishedAt: OffsetDateTime,
        val completedAt: OffsetDateTime?,
        val attempts: Int,
    )

    @GetMapping("/events")
    fun recent(): List<Publication> =
        jdbc.sql(
            """
            SELECT ID, EVENT_TYPE, LISTENER_ID, SERIALIZED_EVENT, STATUS,
                   PUBLICATION_DATE, COMPLETION_DATE, COMPLETION_ATTEMPTS
            FROM EVENT_PUBLICATION
            ORDER BY PUBLICATION_DATE DESC
            LIMIT 60
            """.trimIndent(),
        ).query { rs, _ ->
            val listenerId = rs.getString("LISTENER_ID")
            Publication(
                id = rs.getObject("ID", UUID::class.java),
                eventType = rs.getString("EVENT_TYPE").substringAfterLast('.'),
                module = moduleOf(listenerId),
                listener = listenerId.substringBefore('(').split('.').takeLast(2).joinToString("."),
                event = json.readTree(rs.getString("SERIALIZED_EVENT")),
                status = rs.getString("STATUS") ?: if (rs.getObject("COMPLETION_DATE") != null) "COMPLETED" else "PUBLISHED",
                publishedAt = rs.getObject("PUBLICATION_DATE", OffsetDateTime::class.java),
                completedAt = rs.getObject("COMPLETION_DATE", OffsetDateTime::class.java),
                attempts = rs.getInt("COMPLETION_ATTEMPTS"),
            )
        }.list()

    /** Hands all failed publications to their listeners again. */
    @PostMapping("/events/resubmit")
    fun resubmit() {
        failed.resubmit(ResubmissionOptions.defaults())
    }

    /** "dev.lweiss.northbean.inventory.internal.OrderPlacedListener.on(...)" -> "inventory" */
    private fun moduleOf(listenerId: String): String =
        listenerId.removePrefix(NorthbeanApplication::class.java.packageName + ".").substringBefore('.')
}
