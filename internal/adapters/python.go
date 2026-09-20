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
	// Simple fallback, we assume FastAPI or uvicorn
	if _, err := os.Stat("main.py"); err == nil {
		return []string{"python", "-m", "uvicorn", "main:app"}
	}
	return []string{"python", "app.py"}
}

func (a *PythonAdapter) PortArgs(port int) PortConfig {
	return PortConfig{
		Env: map[string]string{
			"PORT": fmt.Sprintf("%d", port),
		},
		// Append --port if it's uvicorn
		Args: []string{"--port", fmt.Sprintf("%d", port)},
	}
}

func (a *PythonAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
	return nil
}
