package dev.lweiss.northbean.order

import org.springframework.http.HttpStatus
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.ResponseStatus
import org.springframework.web.bind.annotation.RestController
import java.util.UUID

@RestController
@RequestMapping("/orders")
class OrderController(private val orders: OrderService) {

    data class PlaceOrderRequest(val lines: List<OrderLine>)
    data class PlaceOrderResponse(val orderId: UUID)

    @PostMapping
    fun place(@RequestBody request: PlaceOrderRequest): PlaceOrderResponse =
        PlaceOrderResponse(orders.place(request.lines))

    @ExceptionHandler(IllegalArgumentException::class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    fun badRequest(e: IllegalArgumentException) = mapOf("error" to e.message)
}
