# Contributing to HostBind 🚀

Thank you for your interest in contributing! HostBind is a young project and we welcome all kinds of contributions — adapters, bug fixes, docs, tests, and new features.

---

## Development Setup

### Prerequisites
- **Go 1.23+** — [Download here](https://go.dev/dl/)
- **Node.js 16+** — for the npm wrapper scripts
- **Git**

### Clone & Build

```bash
# 1. Fork the repo on GitHub, then clone your fork:
git clone https://github.com/<your-username>/hostbind.git
cd hostbind

# 2. Build the binary:
go build -o hostbind ./cmd/hostbind

# 3. Run it locally:
./hostbind --help
```

### Run Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/registry/...
go test ./internal/adapters/...

# Verbose output
go test -v ./...
```

---

## How to Add a Framework Adapter (10 minutes!)

Framework adapters teach HostBind how to inject the allocated port into a specific dev server. The README calls these out as a high-priority need.

### 1. Create your adapter file

Create `internal/adapters/<framework>.go`:

```go
package adapters

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/RudraMakwana257/hostbind/internal/registry"
)

type MyFrameworkAdapter struct{}

func (a *MyFrameworkAdapter) Name() string { return "myframework" }

// Detect returns true if the current directory looks like a MyFramework project.
// Return a confidence score 0-100 (higher = more confident).
func (a *MyFrameworkAdapter) Detect(dir string) (bool, int) {
    // Example: check for a config file
    if _, err := os.Stat(filepath.Join(dir, "myframework.config.js")); err == nil {
        return true, 90
    }
    return false, 0
}

// DefaultCommand is what HostBind runs when no command is specified.
func (a *MyFrameworkAdapter) DefaultCommand() []string {
    return []string{"npm", "run", "dev"}
}

// PortArgs returns how to inject the port — via env vars, CLI args, or both.
func (a *MyFrameworkAdapter) PortArgs(port int) PortConfig {
    return PortConfig{
        Env:  map[string]string{"PORT": fmt.Sprintf("%d", port)},
        Args: []string{"--port", fmt.Sprintf("%d", port)},
    }
}

func (a *MyFrameworkAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
    return nil
}
```

### 2. Register it in `run.go`

Open `cmd/hostbind/commands/run.go` and add your adapter to `availableAdapters`:

```go
availableAdapters := []adapters.Adapter{
    &adapters.ViteAdapter{},
    &adapters.NextAdapter{},
    &adapters.DjangoAdapter{},
    &adapters.PythonAdapter{},
    &adapters.ExpressAdapter{},
    &adapters.MyFrameworkAdapter{}, // ← add here
}
```

### 3. Write tests

Add `internal/adapters/<framework>_test.go` or tests in `adapters_test.go`:

```go
func TestMyFrameworkAdapter_Detects(t *testing.T) {
    dir := tempDirWith(t, "myframework.config.js")
    a := &adapters.MyFrameworkAdapter{}
    matched, score := a.Detect(dir)
    if !matched { t.Error("expected detection") }
    if score < 80 { t.Errorf("expected score >= 80, got %d", score) }
}
```

### 4. Submit a PR!

Open a pull request with your adapter. We'll review it quickly.

---

## Port Injection Cheat Sheet

| Framework | Method | Example |
|-----------|--------|---------|
| Vite | CLI `--port <N>` | `npm run dev -- --port 4300` |
| Next.js | `PORT` env var | `PORT=4300 npm run dev` |
| Express | `PORT` env var | `PORT=3000 node server.js` |
| Django | Positional arg | `python manage.py runserver 0.0.0.0:8000` |
| FastAPI/uvicorn | CLI `--port <N>` | `uvicorn main:app --port 8001` |
| Flask | `PORT` env var | `PORT=5000 python app.py` |
| Laravel | CLI `--port=<N>` | `php artisan serve --port=8000` |
| Rails | CLI `-p <N>` | `rails server -p 3000` |
| Nuxt | `PORT` env var | `PORT=3000 npm run dev` |
| Astro | CLI `--port <N>` | `npx astro dev --port 4321` |

---

## Code Style

- **Go:** Follow standard `gofmt` formatting. Run `gofmt -w .` before committing.
- **Comments:** Add a comment to every exported function and type.
- **Errors:** Always wrap errors with context: `fmt.Errorf("doing X: %w", err)`.
- **No panics:** Never use `panic()` in library code; return errors instead.

---

## PR Guidelines

1. **One fix per PR** — keeps review fast and focused.
2. **Include tests** — all new adapters and bug fixes should have tests.
3. **Update README** if adding a new adapter (add it to the features list).
4. **Descriptive commit messages** — format: `fix: <what was wrong>` or `feat: add <framework> adapter`.

---

## Need Help?

Open an issue and label it `question`. We're happy to help new contributors get started!
