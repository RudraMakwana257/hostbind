package adapters

import "github.com/RudraMakwana257/hostbind/internal/registry"

// PortConfig represents the environment variables and CLI arguments
// needed to force a service to run on a specific port.
type PortConfig struct {
	Env  map[string]string
	Args []string
}

// Adapter defines how HostBind forces a specific framework to use an allocated port.
type Adapter interface {
	// Name returns the name of the framework (e.g., "vite", "next", "fastapi")
	Name() string
	
	// Detect checks if the current directory matches this framework.
	// Returns a boolean indicating match, and a confidence score (0-100).
	Detect(dir string) (bool, int)
	
	// DefaultCommand returns the default start command for this framework.
	DefaultCommand() []string
	
	// PortArgs returns the Env and Args required to run on the given port.
	PortArgs(port int) PortConfig
	
	// EnvFor returns framework-specific environment variables for other services
	// (e.g., VITE_API_URL) based on the registry's known services.
	EnvFor(serviceName string, reg *registry.Registry) map[string]string
}
