# Modulith Roastery

A learning project: the online shop of Northbean, a small coffee roastery, built as a **modular monolith** twice.

| Folder | Stack | How the modules are kept apart |
|---|---|---|
| [`kotlin/`](kotlin) | Kotlin 2.4, Spring Boot 4.1, Spring Modulith 2.1, Java 21 | Spring Modulith checks the module rules in a test and stores events in `EVENT_PUBLICATION` |
| [`go/`](go) | Go 1.27, standard library only | The compiler enforces `internal/` packages, a small hand-written registry stores events |

Both have the same modules, the same events, the same HTTP API and the same web UI. Read one, then the other,
to see what Spring Modulith does for you and what it takes to do it by hand.

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

`order` does not know `inventory` or `notification`. It only publishes `OrderPlaced`. A fourth module,
`eventlog`, shows the stored event publications in the UI. Order more than is in stock and `inventory` fails
while `notification` succeeds. The failed publication stays stored until you restock and resubmit it.

## Run it

```bash
cd kotlin && ./gradlew bootRun     # needs JDK 21
cd go && go run .                  # needs Go 1.22 or newer, it fetches 1.27 itself
```

Both start on http://localhost:8080, so run one at a time (or set `PORT` for Go).
Tests: `./gradlew test` in `kotlin/`, `go test -race ./...` in `go/`.

## Side by side

| Idea | Kotlin | Go |
|---|---|---|
| A module | direct sub-package of `dev.lweiss.northbean` | direct sub-directory, except `internal/` |
| Hidden internals | sub-packages, checked by `ApplicationModules.verify()` | `internal/` packages, enforced by the compiler |
| Wiring | Spring component scan | `newShop()` in `main.go` |
| Publish an event | `ApplicationEventPublisher.publishEvent` | `events.Publisher.Publish` |
| Listen to an event | `@ApplicationModuleListener` | `events.Listen` with a method `On(ctx, event) error` |
| Event storage | `EVENT_PUBLICATION` table in H2 (in memory) | `internal/events.Registry` (in memory) |
| Listener runs | asynchronously, after the order's transaction has committed | on its own goroutine, after the order is stored |
| Resubmit failed events | `FailedEventPublications.resubmit` | `Registry.ResubmitFailed` |
| Module tests | `@ApplicationModuleTest` with `Scenario` | the module's constructor with a fresh registry |

The web UI in `kotlin/src/main/resources/static/` and `go/static/` is the same. CI fails if the two copies differ.
