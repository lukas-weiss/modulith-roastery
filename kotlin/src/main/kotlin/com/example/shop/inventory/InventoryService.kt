package com.example.shop.inventory

import com.example.shop.inventory.internal.StockRepository
import org.springframework.stereotype.Service

/** Public API of the inventory module. */
@Service
class InventoryService(private val stock: StockRepository) {

    fun stockOf(product: String): Int = stock.get(product)

    fun allStock(): Map<String, Int> = stock.all()

    fun restock(product: String, quantity: Int) {
        require(stock.contains(product)) { "Unknown product: $product" }
        require(quantity in 1..100) { "Quantity must be between 1 and 100" }
        stock.add(product, quantity)
    }
}
