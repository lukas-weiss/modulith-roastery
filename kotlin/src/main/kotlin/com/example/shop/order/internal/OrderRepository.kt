package com.example.shop.order.internal

import org.springframework.stereotype.Repository
import java.util.UUID
import java.util.concurrent.ConcurrentHashMap

/** Kept in memory to keep the example small. */
@Repository
class OrderRepository {

    private val orders = ConcurrentHashMap<UUID, Order>()

    fun save(order: Order): Order {
        orders[order.id] = order
        return order
    }

    fun findAll(): List<Order> = orders.values.toList()
}
