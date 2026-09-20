package allocator

import (
	"fmt"
	"net"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

// Range defines a min and max port.
type Range struct {
	Min int
	Max int
}

var DefaultRanges = map[string]Range{
	"web":    {Min: 4300, Max: 4799},
	"api":    {Min: 5300, Max: 5799},
	"worker": {Min: 6300, Max: 6799},
	"other":  {Min: 7300, Max: 7799},
}

// Allocator manages reserving and freeing ports.
type Allocator struct {
	reg *registry.Registry
}

func New(reg *registry.Registry) *Allocator {
	return &Allocator{reg: reg}
}

// isPortFree checks if a port is available on localhost.
func (a *Allocator) isPortFree(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	l.Close()
	return true
}

// Allocate finds the next available port for a service and reserves it in the registry.
func (a *Allocator) Allocate(project, instance, service, serviceType string) (int, error) {
	// TODO: implement logic
	// 1. Check if it's already allocated in registry.
	// 2. If not, pick next from DefaultRanges[serviceType].
	// 3. Verify it's free with isPortFree.
	// 4. Save to registry with status='reserved'.
	return 4300, nil
}
