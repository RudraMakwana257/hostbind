package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

type ViteAdapter struct{}

func (a *ViteAdapter) Name() string {
	return "vite"
}

func (a *ViteAdapter) Detect(dir string) (bool, int) {
	// Look for vite.config.js, vite.config.ts, or "vite" in package.json
	for _, f := range []string{"vite.config.js", "vite.config.ts"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			return true, 90
		}
	}
	return false, 0
}

func (a *ViteAdapter) DefaultCommand() []string {
	return []string{"npm", "run", "dev"}
}

func (a *ViteAdapter) PortArgs(port int) PortConfig {
	return PortConfig{
		Env: nil,
		Args: []string{"--port", fmt.Sprintf("%d", port)},
	}
}

func (a *ViteAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
	// If Vite needs to talk to an API, it usually uses VITE_API_URL
	// In a real implementation, we would query the registry for the 'api' service URL.
	// For now, we return empty or hardcode a lookup.
	return nil
}
