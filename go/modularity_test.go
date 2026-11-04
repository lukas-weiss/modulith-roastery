package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// The Go compiler already enforces two of Spring Modulith's rules: only packages inside order/ may
// import order/internal, and packages may not import each other in a cycle. This test checks the same
// on the level of modules, prints the module structure and writes a PlantUML diagram.

const modulePath = "example.com/shop/"

// modules maps each module (direct sub-directory with Go code, except internal/) to the modules it uses.
func modules(t *testing.T) map[string][]string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	result := map[string][]string{}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "internal" || !hasGoFiles(entry.Name()) {
			continue
		}
		result[entry.Name()] = dependencies(t, entry.Name())
	}
	return result
}

// dependencies lists the other modules that the non-test code of module imports. It fails if an import
// reaches into another module's sub-packages.
func dependencies(t *testing.T, module string) []string {
	t.Helper()
	var uses []string
	err := filepath.WalkDir(module, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			imported := strings.Trim(spec.Path.Value, `"`)
			rest, ok := strings.CutPrefix(imported, modulePath)
			if !ok {
				continue
			}
			other, sub, _ := strings.Cut(rest, "/")
			if other == module || other == "internal" {
				continue
			}
			if sub != "" {
				t.Errorf("Module '%s' depends on non-exposed package %s within module '%s'!", module, imported, other)
			}
			if !slices.Contains(uses, other) {
				uses = append(uses, other)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(uses)
	return uses
}

func hasGoFiles(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(path string, _ os.DirEntry, _ error) error {
		found = found || strings.HasSuffix(path, ".go")
		return nil
	})
	return found
}

// Fails if a module uses another module's internals or if modules depend on each other in a cycle.
func TestModuleStructureIsValid(t *testing.T) {
	all := modules(t)
	for _, name := range slices.Sorted(maps.Keys(all)) {
		t.Logf("# %s\n  depends on: %v", name, all[name])
	}

	var visit func(module string, path []string)
	visit = func(module string, path []string) {
		if i := slices.Index(path, module); i >= 0 {
			t.Fatalf("Cycle between modules: %s", strings.Join(append(path[i:], module), " -> "))
		}
		for _, next := range all[module] {
			visit(next, append(slices.Clone(path), module))
		}
	}
	for module := range all {
		visit(module, nil)
	}

	if uses := all["order"]; len(uses) > 0 {
		t.Errorf("order must not know other modules, but uses %v", uses)
	}
}

// Writes a PlantUML component diagram of the modules to build/modulith-docs.
func TestWriteDocumentation(t *testing.T) {
	all := modules(t)
	var b strings.Builder
	b.WriteString("@startuml\n")
	for _, name := range slices.Sorted(maps.Keys(all)) {
		fmt.Fprintf(&b, "component [%s]\n", name)
	}
	for _, name := range slices.Sorted(maps.Keys(all)) {
		for _, uses := range all[name] {
			fmt.Fprintf(&b, "[%s] --> [%s] : uses\n", name, uses)
		}
	}
	b.WriteString("@enduml\n")

	dir := filepath.Join("build", "modulith-docs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "components.puml"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}
