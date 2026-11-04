# Modulith Roastery (Kotlin)

The online shop of a small coffee roastery (crema beans, espresso beans, espresso cups). It shows the core ideas of [Spring Modulith](https://spring.io/projects/spring-modulith):
one Spring Boot application, split into **modules** that talk to each other through **events**.

Kotlin 2.4 · Spring Boot 4.1 · Spring Modulith 2.1 · Java 21

The same shop without a framework, in Go: [`../go`](../go). Overview of both: [`../README.md`](../README.md).

## The idea in one picture

```
POST /orders
     │
     ▼
┌──────────┐   OrderPlaced event   ┌─────────────┐
│  order   │ ────────────────────► │  inventory  │  reduces stock   (GET /inventory)
└──────────┘          │            └─────────────┘
                      │            ┌──────────────┐
                      └──────────► │ notification │  logs a confirmation
                                   └──────────────┘
```

`order` does not know `inventory` or `notification`. It only publishes `OrderPlaced`.
New features (e.g. invoicing) can listen to the same event without changing `order`.

## Project layout

```
src/main/kotlin/com/example/shop/
├── ShopApplication.kt
├── order/                      ← module "order"
│   ├── OrderPlaced.kt          public API: the event
│   ├── OrderLine.kt            public API: one product and quantity of an order
│   ├── OrderService.kt         public API: place an order
│   ├── OrderController.kt
│   └── internal/               hidden from other modules
│       ├── Order.kt
│       └── OrderRepository.kt
├── inventory/                  ← module "inventory"
│   ├── InventoryService.kt
│   ├── InventoryController.kt
│   └── internal/
│       ├── OrderPlacedListener.kt
│       └── StockRepository.kt
├── notification/               ← module "notification"
│   └── NotificationListener.kt
└── eventlog/                   ← module "eventlog": reads EVENT_PUBLICATION for the UI
    └── EventLogController.kt

src/main/resources/static/      ← the web UI (plain HTML, CSS, JS, no build step)
```

## The four rules to know

1. **Each direct sub-package of the main package is a module.** Here: `order`, `inventory`, `notification`.
   No extra configuration is needed.
2. **The base package of a module is its public API, sub-packages are internal.** `inventory` may use
   `order.OrderPlaced`, but not `order.internal.Order`.
3. **Modules talk through events.** `OrderService` calls `publishEvent(OrderPlaced(...))`. Listeners are
   annotated with `@ApplicationModuleListener`. They run asynchronously in their own transaction, after the
   order's transaction has committed.
4. **Events are stored until they are handled.** `spring-modulith-starter-jdbc` writes each event to the
   `EVENT_PUBLICATION` table (in-memory H2 here). If a listener fails, the event stays marked as incomplete
   and is not lost.

## Tests show what Modulith checks

| Test | What it shows |
|---|---|
| `ModularityTests` | `ApplicationModules.verify()` fails the build if a module uses another module's `internal` classes or if modules depend on each other in a cycle. It also writes diagrams to `build/spring-modulith-docs/`. |
| `OrderModuleTests` | `@ApplicationModuleTest` starts **only** the order module. `Scenario` checks that placing an order publishes `OrderPlaced`. |
| `InventoryModuleTests` | Starts only the inventory module, publishes an `OrderPlaced` event by hand, then waits until the stock has changed. |

Try breaking a rule: in `notification`, import `com.example.shop.order.internal.Order` and run the tests.
`ModularityTests` fails with:

```
Module 'notification' depends on non-exposed type com.example.shop.order.internal.Order within module 'order'!
```

## Run it

You need JDK 21 (`JAVA_HOME` set).

```bash
./gradlew test        # run all tests
./gradlew bootRun     # start on http://localhost:8080
```

Open http://localhost:8080 in a browser. Place an order and watch the event travel from `order` to
`inventory` and `notification`. The event log on the right shows each row of the `EVENT_PUBLICATION` table.

An order can contain several products. To see a failing listener, order more than is in stock. `inventory`
fails for the whole order (it takes all products out of stock or none), `notification` still succeeds, and the
event publication for `inventory` stays stored as `FAILED`. Restock the product, then click "Resubmit" to hand
the event to the listener again.

The page is in English by default. The switch in the top right corner changes it to German.

Changes to the files in `src/main/resources/static/` need a restart of `bootRun`. The Go version has the same files in `go/static/`; keep both copies identical (CI checks it).

The same with curl:

```bash
curl localhost:8080/inventory
# {"crema-beans":12,"espresso-beans":8,"espresso-cups":6}

curl -X POST localhost:8080/orders -H 'Content-Type: application/json' \
     -d '{"lines":[{"product":"espresso-beans","quantity":3},{"product":"espresso-cups","quantity":2}]}'
# {"orderId":"..."}

curl localhost:8080/inventory
# {"crema-beans":12,"espresso-beans":5,"espresso-cups":4}
```

The log shows both listeners reacting to the same event on their own threads:

```
[task-1] ...OrderPlacedListener   : Inventory: reduced stock for order ...: 3 x espresso-beans, 2 x espresso-cups
[task-2] ...NotificationListener  : Notification: order ... confirmed (3 x espresso-beans, 2 x espresso-cups)
```
