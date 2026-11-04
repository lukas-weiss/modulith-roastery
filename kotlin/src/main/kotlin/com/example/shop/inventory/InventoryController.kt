package com.example.shop.inventory

import org.springframework.http.HttpStatus
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.ResponseStatus
import org.springframework.web.bind.annotation.RestController

@RestController
class InventoryController(private val inventory: InventoryService) {

    data class RestockRequest(val product: String, val quantity: Int)

    @GetMapping("/inventory")
    fun all(): Map<String, Int> = inventory.allStock()

    @PostMapping("/inventory/restock")
    fun restock(@RequestBody request: RestockRequest): Map<String, Int> {
        inventory.restock(request.product, request.quantity)
        return inventory.allStock()
    }

    @ExceptionHandler(IllegalArgumentException::class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    fun badRequest(e: IllegalArgumentException) = mapOf("error" to e.message)
}
