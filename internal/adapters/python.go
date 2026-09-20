package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

type PythonAdapter struct{}

func (a *PythonAdapter) Name() string {
	return "python (fastapi/flask)"
}

func (a *PythonAdapter) Detect(dir string) (bool, int) {
	for _, f := range []string{"requirements.txt", "Pipfile", "pyproject.toml", "main.py", "app.py"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			return true, 80 // High confidence if Python files exist
		}
	}
	return false, 0
}

func (a *PythonAdapter) DefaultCommand() []string {
	// Prefer uvicorn if main.py exists (FastAPI pattern)
	if _, err := os.Stat("main.py"); err == nil {
		return []string{"python", "-m", "uvicorn", "main:app"}
	}
	// Fall back to plain python app.py (Flask pattern)
	return []string{"python", "app.py"}
}

// PortArgs returns the port configuration for the detected Python framework.
//
// Fix #8: Previously this always appended --port even for Flask apps
// (plain `python app.py`), which don't accept that CLI flag and would crash.
//
// Now we distinguish between uvicorn (accepts --port) and plain Python
// (uses PORT env var only). If a custom command is passed, the user is
// responsible for port injection.
func (a *PythonAdapter) PortArgs(port int) PortConfig {
	portStr := fmt.Sprintf("%d", port)

	// Check if the default command is uvicorn-based (main.py exists)
	_, mainExists := os.Stat("main.py")
	isUvicorn := mainExists == nil // main.py present → we use uvicorn

	if isUvicorn {
		// uvicorn accepts --port <port> as a CLI argument
		return PortConfig{
			Env:  map[string]string{"PORT": portStr},
			Args: []string{"--port", portStr},
		}
	}

	// Flask / plain python: inject via PORT env var only.
	// Adding --port here would break: `python app.py --port 4300` is invalid.
	return PortConfig{
		Env:  map[string]string{"PORT": portStr},
		Args: nil,
	}
}

func (a *PythonAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
	return nil
}
