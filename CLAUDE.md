# Modulith Roastery – Regeln für Claude

Lernprojekt: derselbe Shop einer Kaffeerösterei als modularer Monolith in `kotlin/` (Spring Modulith) und `go/`
(nur Standardbibliothek). Details zu jeder Version stehen in ihrem README.

- Beide Versionen verhalten sich gleich: Module, Events, HTTP-API, Validierung, Log-Texte. Eine Änderung an der
  Logik kommt immer in beide, im selben Commit.
- Die UI-Dateien in `kotlin/src/main/resources/static/` und `go/static/` sind identisch (die CI prüft das mit
  `diff -r`). Nur eine ändern und dann kopieren. Texte so formulieren, dass sie für beide Versionen passen.
- Ein Ordner `shared/` für die UI geht nicht: `go:embed` kann nichts außerhalb von `go/` einbinden.
- Vor einem Commit: `./gradlew test` in `kotlin/`, `gofmt -l .`, `go vet ./...` und `go test -race ./...` in `go/`.
- Selbst gestartete Server (Port 8080) nach dem Test sofort stoppen.
