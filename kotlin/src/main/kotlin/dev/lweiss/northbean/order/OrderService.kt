package dev.lweiss.northbean.order

import dev.lweiss.northbean.order.internal.Order
import dev.lweiss.northbean.order.internal.OrderRepository
import org.springframework.context.ApplicationEventPublisher
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional
import java.util.UUID

/** Public API of the order module. */
@Service
class OrderService(
    private val repository: OrderRepository,
    private val events: ApplicationEventPublisher,
) {

    @Transactional
    fun place(lines: List<OrderLine>): UUID {
        require(lines.isNotEmpty()) { "An order needs at least one line" }
        require(lines.all { it.quantity in 1..100 }) { "Quantity must be between 1 and 100" }
        require(lines.distinctBy { it.product }.size == lines.size) { "Each product may appear only once" }

        val order = repository.save(Order(UUID.randomUUID(), lines))

        // No call to inventory or notification here – just tell the world what happened.
        // Listeners run after this transaction commits.
        events.publishEvent(OrderPlaced(order.id, order.lines))

        return order.id
    }
}
