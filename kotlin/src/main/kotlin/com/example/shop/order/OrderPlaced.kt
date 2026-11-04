package com.example.shop.order

import java.util.UUID

/**
 * Domain event published by the order module.
 * It lives in the module's base package, so it is part of the module's public API
 * and other modules may use it.
 */
data class OrderPlaced(
    val orderId: UUID,
    val lines: List<OrderLine>,
)
