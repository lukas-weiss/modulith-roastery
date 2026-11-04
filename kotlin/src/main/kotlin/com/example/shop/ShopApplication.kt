package com.example.shop

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.scheduling.annotation.EnableAsync

/**
 * Every direct sub-package of `com.example.shop` is an application module:
 *
 * - `order`        – accepts orders and publishes [com.example.shop.order.OrderPlaced]
 * - `inventory`    – reacts to OrderPlaced and reduces the stock
 * - `notification` – reacts to OrderPlaced and "sends" a confirmation
 *
 * `order` knows nothing about the other two. They only talk through events.
 */
@SpringBootApplication
@EnableAsync // @ApplicationModuleListener runs listeners asynchronously
class ShopApplication

fun main(args: Array<String>) {
    runApplication<ShopApplication>(*args)
}
