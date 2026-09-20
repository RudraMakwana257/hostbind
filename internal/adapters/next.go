package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

type NextAdapter struct{}

func (a *NextAdapter) Name() string {
	return "next"
}

func (a *NextAdapter) Detect(dir string) (bool, int) {
	for _, f := range []string{"next.config.js", "next.config.mjs"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			return true, 90
		}
	}
	return false, 0
}

func (a *NextAdapter) DefaultCommand() []string {
	return []string{"npm", "run", "dev"}
}

func (a *NextAdapter) PortArgs(port int) PortConfig {
	return PortConfig{
		Env: map[string]string{
			"PORT": fmt.Sprintf("%d", port),
		},
		Args: nil,
	}
}

func (a *NextAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
	return nil
}
