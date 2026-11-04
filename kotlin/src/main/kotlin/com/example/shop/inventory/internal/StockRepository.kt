package com.example.shop.inventory.internal

import org.springframework.stereotype.Repository

/** Kept in memory to keep the example small. */
@Repository
class StockRepository {

    private val stock = mutableMapOf(
        "crema-beans" to 12,
        "espresso-beans" to 8,
        "espresso-cups" to 6,
    )

    @Synchronized
    fun get(product: String): Int = stock[product] ?: 0

    /**
     * Takes all quantities out of stock, or none: if one product is short, nothing changes and this fails.
     * The listener then fails and its event publication stays incomplete.
     */
    @Synchronized
    fun reduce(quantities: Map<String, Int>) {
        val short = quantities.filter { (product, quantity) -> get(product) < quantity }
        check(short.isEmpty()) {
            "Not enough in stock: " + short.entries.joinToString { (product, quantity) -> "$product (${get(product)} in stock, $quantity ordered)" }
        }
        quantities.forEach { (product, quantity) -> stock[product] = get(product) - quantity }
    }

    @Synchronized
    fun add(product: String, quantity: Int) {
        stock[product] = get(product) + quantity
    }

    @Synchronized
    fun contains(product: String): Boolean = product in stock

    @Synchronized
    fun all(): Map<String, Int> = stock.toSortedMap()
}
