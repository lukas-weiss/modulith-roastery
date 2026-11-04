package com.example.shop.inventory

import com.example.shop.order.OrderLine
import com.example.shop.order.OrderPlaced
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.modulith.test.ApplicationModuleTest
import org.springframework.modulith.test.Scenario
import org.junit.jupiter.api.Test
import java.util.UUID
import kotlin.test.assertEquals

/** Starts only the inventory module and feeds it an event, as if the order module had sent it. */
@ApplicationModuleTest
class InventoryModuleTests @Autowired constructor(
    private val inventory: InventoryService,
) {

    @Test
    fun `stock is reduced when an order is placed`(scenario: Scenario) {
        scenario.publish(OrderPlaced(UUID.randomUUID(), listOf(OrderLine("espresso-cups", 2))))
            .andWaitForStateChange({ inventory.stockOf("espresso-cups") }) { it == 4 }
            .andVerify { stock -> assertEquals(4, stock) }
    }
}
