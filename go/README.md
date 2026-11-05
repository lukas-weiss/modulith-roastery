# Modulith Roastery (Go)

The Go twin of [`../kotlin`](../kotlin): the same online shop of Northbean, a small coffee roastery (crema
beans, espresso beans, espresso cups), the same modules, events, HTTP API and web UI. There is no Spring
Modulith for Go, so the parts it provides are written out by hand, with nothing but the standard library.

Go 1.27 · no dependencies

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
main.go                     wires the modules together (Spring does this in Kotlin)
order/                      ← module "order"
├── events.go               public API: OrderPlaced and OrderLine
├── service.go              public API: place an order
├── controller.go
└── internal/               hidden from other modules
    └── order.go            Order and OrderRepository
inventory/                  ← module "inventory"
├── service.go
├── controller.go
└── internal/
    ├── listener.go         OrderPlacedListener
    └── stock.go            StockRepository
notification/               ← module "notification"
└── listener.go
eventlog/                   ← module "eventlog": shows the stored publications in the UI
└── controller.go
internal/                   shared plumbing, not a module
├── events/                 the event publication registry (what Spring Modulith does in Kotlin)
├── check/                  Require(), Go's counterpart of Kotlin's require()
├── web/                    JSON helpers for the handlers
└── uuid/
static/                     ← the web UI, embedded into the binary
```

## The four rules, in Go

| Kotlin with Spring Modulith | Go |
|---|---|
| Each direct sub-package of the main package is a module. | Each direct sub-directory is a module, except `internal/`. |
| Sub-packages of a module are internal. `ApplicationModules.verify()` fails the build otherwise. | The same, but the **compiler** enforces it: only code inside `order/` may import `order/internal`. Import cycles do not compile either. |
| Listeners are annotated with `@ApplicationModuleListener` and run asynchronously. | A listener has a method `On(ctx, event) error` and is registered with `events.Listen`. Each runs on its own goroutine. |
| `spring-modulith-starter-jdbc` stores each event per listener in `EVENT_PUBLICATION` (H2 in memory). | `internal/events.Registry` stores each event per listener, in memory. A listener that returns an error (or panics) leaves its publication `FAILED`; `ResubmitFailed` hands it over again. |

One difference: there are no transactions here. In Kotlin the listeners run after the order's transaction has
committed. In Go, `order` publishes the event after the order is stored.

The registry stores the event as JSON and listeners get it back from that JSON, as Spring Modulith does.
So an event must survive a round trip through `encoding/json`.

## Tests

| Test | What it shows |
|---|---|
| `modularity_test.go` | Lists the modules and what they depend on, fails on a use of another module's sub-packages, on a cycle, or if `order` uses another module. Writes a PlantUML diagram to `build/modulith-docs/`. |
| `order/service_test.go` | Starts only the order module and checks that placing an order publishes `OrderPlaced`. |
| `inventory/service_test.go` | Starts only the inventory module, publishes an `OrderPlaced` event by hand, then waits until the stock has changed. |
| `inventory/internal/stock_test.go` | The stock changes for all products of an order, or for none. |
| `internal/events/registry_test.go` | A failed publication stays stored and completes on resubmit, a panicking listener fails. |
| `eventlog/controller_test.go` | `moduleOf` finds the module of a listener ID, also when the ID has no separator after the prefix. |
| `shop_test.go` | The whole shop over HTTP: order too much, inventory fails, restock, resubmit. |

Try breaking a rule: in `notification`, import `lweiss.dev/northbean/order/internal`. Nothing compiles:

```
notification/listener.go:12:2: use of internal package lweiss.dev/northbean/order/internal not allowed
```

## Run it

```bash
go test -race ./...   # run all tests
go run .              # start on http://localhost:8080 (PORT=… to change it)
```

Open http://localhost:8080 in a browser. Place an order and watch the event travel from `order` to
`inventory` and `notification`. The event log on the right shows each stored publication.

An order can contain several products. To see a failing listener, order more than is in stock. `inventory`
fails for the whole order (it takes all products out of stock or none), `notification` still succeeds, and the
publication for `inventory` stays stored as `FAILED`. Restock the product, then click "Resubmit" to hand the
event to the listener again.

The page is in English by default. The switch in the top right corner changes it to German. The files in
`static/` are embedded into the binary, so changes to them need a restart. The Kotlin version has the same files
in `kotlin/src/main/resources/static/`; keep both copies identical (CI checks it).

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

The log shows both listeners reacting to the same event:

```
INFO Inventory: reduced stock for order ...: 3 x espresso-beans, 2 x espresso-cups
INFO Notification: order ... confirmed (3 x espresso-beans, 2 x espresso-cups)
```
