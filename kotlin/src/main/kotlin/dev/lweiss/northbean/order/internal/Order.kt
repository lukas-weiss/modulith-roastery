package dev.lweiss.northbean.order.internal

import dev.lweiss.northbean.order.OrderLine
import java.util.UUID

/** Sub-packages are internal: other modules must not use this class. */
data class Order(
    val id: UUID,
    val lines: List<OrderLine>,
)
