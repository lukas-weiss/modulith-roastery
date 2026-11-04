package com.example.shop.order

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.modulith.test.ApplicationModuleTest
import org.springframework.modulith.test.Scenario
import org.junit.jupiter.api.Test

/** Starts only the order module – inventory and notification are not loaded. */
@ApplicationModuleTest
class OrderModuleTests @Autowired constructor(
    private val orders: OrderService,
) {

    @Test
    fun `placing an order publishes OrderPlaced`(scenario: Scenario) {
        scenario.stimulate { _ -> orders.place(listOf(OrderLine("crema-beans", 2), OrderLine("espresso-cups", 1))) }
            .andWaitForEventOfType(OrderPlaced::class.java)
            .matching { it.lines == listOf(OrderLine("crema-beans", 2), OrderLine("espresso-cups", 1)) }
            .toArrive()
    }
}
