package com.example.shop.inventory.internal

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows
import kotlin.test.assertEquals

class StockRepositoryTests {

    private val stock = StockRepository()

    @Test
    fun `reduces every product of an order`() {
        stock.reduce(mapOf("crema-beans" to 2, "espresso-cups" to 1))

        assertEquals(10, stock.get("crema-beans"))
        assertEquals(5, stock.get("espresso-cups"))
    }

    @Test
    fun `changes nothing if one product is short`() {
        assertThrows<IllegalStateException> {
            stock.reduce(mapOf("crema-beans" to 2, "espresso-cups" to 7))
        }

        assertEquals(12, stock.get("crema-beans"))
        assertEquals(6, stock.get("espresso-cups"))
    }
}
