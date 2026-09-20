package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

// ExpressAdapter handles Node.js Express.js applications.
// Express reads the PORT environment variable at startup.
type ExpressAdapter struct{}

func (a *ExpressAdapter) Name() string {
	return "express"
}

// Detect checks for Express.js project markers.
// Looks for express in package.json dependencies or an app.js/server.js entry point.
func (a *ExpressAdapter) Detect(dir string) (bool, int) {
	pkgPath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false, 0
	}

	content := string(data)
	// Check for express dependency in package.json
	if contains(content, `"express"`) {
		// Check for typical Express entry files for higher confidence
		for _, f := range []string{"app.js", "server.js", "index.js"} {
			if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
				return true, 85
			}
		}
		return true, 70 // express in package.json but no typical entry file
	}
	return false, 0
}

func (a *ExpressAdapter) DefaultCommand() []string {
	// Try common entry points in priority order
	for _, f := range []string{"server.js", "app.js", "index.js"} {
		if _, err := os.Stat(f); err == nil {
			return []string{"node", f}
		}
	}
	return []string{"node", "index.js"}
}

// PortArgs injects the PORT environment variable.
// Express apps typically read process.env.PORT with:
//
//	const port = process.env.PORT || 3000;
//	app.listen(port, ...);
func (a *ExpressAdapter) PortArgs(port int) PortConfig {
	return PortConfig{
		Env:  map[string]string{"PORT": fmt.Sprintf("%d", port)},
		Args: nil, // Express doesn't accept --port as a CLI flag
	}
}

func (a *ExpressAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
	return nil
}

// contains is a simple substring check helper.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
