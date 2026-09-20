package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

// DjangoAdapter handles Python Django web framework applications.
// Django uses `python manage.py runserver 0.0.0.0:<port>` to bind a specific port.
type DjangoAdapter struct{}

func (a *DjangoAdapter) Name() string {
	return "django"
}

// Detect checks for Django project markers.
// A Django project always has a manage.py file at the root.
// We also check for django in requirements.txt / pyproject.toml for extra confidence.
func (a *DjangoAdapter) Detect(dir string) (bool, int) {
	// manage.py is the most reliable Django indicator
	if _, err := os.Stat(filepath.Join(dir, "manage.py")); err != nil {
		return false, 0
	}

	// Check for django in requirements for higher confidence
	for _, reqFile := range []string{"requirements.txt", "Pipfile", "pyproject.toml"} {
		data, err := os.ReadFile(filepath.Join(dir, reqFile))
		if err != nil {
			continue
		}
		if containsAny(string(data), []string{"Django", "django"}) {
			return true, 95 // manage.py + django in requirements = very high confidence
		}
	}

	// manage.py alone is enough (could be Django or another framework using it)
	return true, 80
}

// DefaultCommand returns the Django development server command.
// We bind to 0.0.0.0 to expose on all interfaces (standard for dev).
// The port placeholder will be replaced by PortArgs.
func (a *DjangoAdapter) DefaultCommand() []string {
	return []string{"python", "manage.py", "runserver"}
}

// PortArgs for Django appends the address:port to the runserver command.
// Django runserver syntax: python manage.py runserver [addr:]port
func (a *DjangoAdapter) PortArgs(port int) PortConfig {
	portStr := fmt.Sprintf("%d", port)
	return PortConfig{
		Env: map[string]string{
			"PORT":     portStr,
			"DJANGO_DEV_PORT": portStr,
		},
		// Django runserver takes "addr:port" or just "port" as a positional arg
		Args: []string{fmt.Sprintf("0.0.0.0:%d", port)},
	}
}

func (a *DjangoAdapter) EnvFor(serviceName string, reg *registry.Registry) map[string]string {
	return nil
}

// containsAny checks if s contains any of the given substrings.
func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if contains(s, sub) {
			return true
		}
	}
	return false
}
