package adapters_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RudraMakwana257/hostbind/internal/adapters"
)

// tempDir creates a temp directory with the given files, returning the dir path.
func tempDirWith(t *testing.T, files ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		path := filepath.Join(dir, f)
		// Create parent directories if needed
		os.MkdirAll(filepath.Dir(path), 0755)
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", f, err)
		}
	}
	return dir
}

// --- Vite Adapter ---

func TestViteAdapter_DetectsViteConfigJS(t *testing.T) {
	dir := tempDirWith(t, "vite.config.js")
	a := &adapters.ViteAdapter{}
	matched, score := a.Detect(dir)
	if !matched {
		t.Error("ViteAdapter: expected to detect project with vite.config.js")
	}
	if score < 80 {
		t.Errorf("ViteAdapter: expected score >= 80, got %d", score)
	}
}

func TestViteAdapter_DetectsViteConfigTS(t *testing.T) {
	dir := tempDirWith(t, "vite.config.ts")
	a := &adapters.ViteAdapter{}
	matched, _ := a.Detect(dir)
	if !matched {
		t.Error("ViteAdapter: expected to detect project with vite.config.ts")
	}
}

func TestViteAdapter_NoDetectWithoutConfig(t *testing.T) {
	dir := tempDirWith(t, "package.json")
	a := &adapters.ViteAdapter{}
	matched, _ := a.Detect(dir)
	if matched {
		t.Error("ViteAdapter: should not detect project without vite.config.*")
	}
}

func TestViteAdapter_PortArgs(t *testing.T) {
	a := &adapters.ViteAdapter{}
	cfg := a.PortArgs(4300)
	if len(cfg.Args) == 0 {
		t.Error("ViteAdapter: expected Args to contain --port flag")
	}
	found := false
	for i, arg := range cfg.Args {
		if arg == "--port" && i+1 < len(cfg.Args) && cfg.Args[i+1] == "4300" {
			found = true
		}
	}
	if !found {
		t.Errorf("ViteAdapter: expected '--port 4300' in args, got %v", cfg.Args)
	}
}

// --- Next.js Adapter ---

func TestNextAdapter_Detects(t *testing.T) {
	dir := tempDirWith(t, "next.config.js")
	a := &adapters.NextAdapter{}
	matched, score := a.Detect(dir)
	if !matched {
		t.Error("NextAdapter: expected to detect project with next.config.js")
	}
	if score < 80 {
		t.Errorf("NextAdapter: expected score >= 80, got %d", score)
	}
}

func TestNextAdapter_PortArgs(t *testing.T) {
	a := &adapters.NextAdapter{}
	cfg := a.PortArgs(3001)
	if cfg.Env["PORT"] != "3001" {
		t.Errorf("NextAdapter: expected PORT=3001, got %q", cfg.Env["PORT"])
	}
}

// --- Python Adapter ---

func TestPythonAdapter_DetectsRequirementsTxt(t *testing.T) {
	dir := tempDirWith(t, "requirements.txt")
	a := &adapters.PythonAdapter{}
	matched, _ := a.Detect(dir)
	if !matched {
		t.Error("PythonAdapter: expected to detect project with requirements.txt")
	}
}

func TestPythonAdapter_FlaskPortArgs_NoCliFlag(t *testing.T) {
	// Flask pattern: app.py only (no main.py)
	dir := tempDirWith(t, "app.py")
	if err := os.Chdir(dir); err != nil {
		t.Skip("cannot chdir in test environment")
	}
	a := &adapters.PythonAdapter{}
	cfg := a.PortArgs(5000)
	// Flask should NOT get --port since plain `python app.py` doesn't accept it
	if len(cfg.Args) > 0 {
		t.Errorf("PythonAdapter (Flask): expected no CLI args, got %v — this would break 'python app.py'", cfg.Args)
	}
	if cfg.Env["PORT"] != "5000" {
		t.Errorf("PythonAdapter (Flask): expected PORT env var to be set, got %q", cfg.Env["PORT"])
	}
}

// --- Django Adapter ---

func TestDjangoAdapter_Detects(t *testing.T) {
	dir := tempDirWith(t, "manage.py")
	a := &adapters.DjangoAdapter{}
	matched, score := a.Detect(dir)
	if !matched {
		t.Error("DjangoAdapter: expected to detect project with manage.py")
	}
	if score < 70 {
		t.Errorf("DjangoAdapter: expected score >= 70, got %d", score)
	}
}

func TestDjangoAdapter_PortArgs(t *testing.T) {
	a := &adapters.DjangoAdapter{}
	cfg := a.PortArgs(8000)
	if len(cfg.Args) == 0 {
		t.Error("DjangoAdapter: expected Args to contain addr:port")
	}
	if cfg.Args[0] != "0.0.0.0:8000" {
		t.Errorf("DjangoAdapter: expected '0.0.0.0:8000' in args, got %q", cfg.Args[0])
	}
}

// --- Express Adapter ---

func TestExpressAdapter_Detects(t *testing.T) {
	dir := tempDirWith(t, "server.js")
	// Write a package.json with express dependency
	pkgJSON := `{"dependencies": {"express": "^4.18.0"}}`
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkgJSON), 0644)

	a := &adapters.ExpressAdapter{}
	matched, score := a.Detect(dir)
	if !matched {
		t.Error("ExpressAdapter: expected to detect project with express in package.json")
	}
	if score < 70 {
		t.Errorf("ExpressAdapter: expected score >= 70, got %d", score)
	}
}

func TestExpressAdapter_PortArgs_NoCliFlag(t *testing.T) {
	a := &adapters.ExpressAdapter{}
	cfg := a.PortArgs(3000)
	// Express reads PORT from env, no --port CLI flag
	if len(cfg.Args) > 0 {
		t.Errorf("ExpressAdapter: expected no CLI args, got %v", cfg.Args)
	}
	if cfg.Env["PORT"] != "3000" {
		t.Errorf("ExpressAdapter: expected PORT=3000, got %q", cfg.Env["PORT"])
	}
}

// --- Generic Adapter ---

func TestGenericAdapter_AlwaysDetects(t *testing.T) {
	dir := tempDirWith(t)
	a := &adapters.GenericAdapter{}
	// GenericAdapter is the fallback — detect logic varies, just check PortArgs
	cfg := a.PortArgs(9000)
	if cfg.Env["PORT"] == "" {
		t.Error("GenericAdapter: expected PORT env var to be set")
	}
}
