package com.example.shop.order

/** One line of an order: a product and how many of it. Part of the order module's public API. */
data class OrderLine(
    val product: String,
    val quantity: Int,
)
