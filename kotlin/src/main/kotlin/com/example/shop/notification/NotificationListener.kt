package com.example.shop.notification

import com.example.shop.order.OrderPlaced
import org.slf4j.LoggerFactory
import org.springframework.modulith.events.ApplicationModuleListener
import org.springframework.stereotype.Component

@Component
class NotificationListener {

    private val log = LoggerFactory.getLogger(javaClass)

    @ApplicationModuleListener
    fun on(event: OrderPlaced) {
        // A real app would send an email here.
        log.info(
            "Notification: order {} confirmed ({})",
            event.orderId,
            event.lines.joinToString { "${it.quantity} x ${it.product}" },
        )
    }
}
