package adapters

import (
	"os"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

type GenericAdapter struct{}

func (a *GenericAdapter) Name() string {
	return "generic"
}

func (a *GenericAdapter) Detect(dir string) (bool, int) {
	// Generic adapter is the fallback.
	return true, 10
}

func (a *GenericAdapter) DefaultCommand() []string {
	// Need a package.json check for npm start, otherwise python, etc.
	// We'll keep it simple for now.
	if _, err := os.Stat("package.json"); err == nil {
		return []string{"npm", "start"}
	}
	return []string{"./start.sh"}
}

func (a *GenericAdapter) PortArgs(port int) PortConfig {
	return PortConfig{
		Env: map[string]string{
			"PORT": string(port),
		},
		Args: nil,
	}
}

func (a *GenericAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
	return nil
}
