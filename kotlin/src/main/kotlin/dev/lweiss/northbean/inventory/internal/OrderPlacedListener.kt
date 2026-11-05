package dev.lweiss.northbean.inventory.internal

import dev.lweiss.northbean.order.OrderPlaced
import org.slf4j.LoggerFactory
import org.springframework.modulith.events.ApplicationModuleListener
import org.springframework.stereotype.Component

@Component
class OrderPlacedListener(private val stock: StockRepository) {

    private val log = LoggerFactory.getLogger(javaClass)

    /**
     * @ApplicationModuleListener = runs asynchronously, in its own transaction,
     * after the publishing transaction has committed.
     */
    @ApplicationModuleListener
    fun on(event: OrderPlaced) {
        stock.reduce(event.lines.associate { it.product to it.quantity })
        log.info("Inventory: reduced stock for order {}: {}", event.orderId, describe(event))
    }

    private fun describe(event: OrderPlaced) = event.lines.joinToString { "${it.quantity} x ${it.product}" }
}
