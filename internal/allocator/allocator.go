package allocator

import (
	"fmt"
	"net"
	"os"

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

type Allocator struct {
	reg *registry.Registry
}

func New(reg *registry.Registry) *Allocator {
	return &Allocator{reg: reg}
}

func (a *Allocator) isPortFree(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	l.Close()
	return true
}

func (a *Allocator) Allocate(project, instance, service, serviceType string) (int, error) {
	wd, _ := os.Getwd()
	projID, err := a.reg.EnsureProject(project, instance, wd)
	if err != nil {
		return 0, fmt.Errorf("failed to ensure project: %w", err)
	}

	allocs, err := a.reg.GetActiveAllocations()
	if err != nil {
		return 0, fmt.Errorf("failed to get active allocations: %w", err)
	}

	usedPorts := make(map[int]bool)
	for _, alloc := range allocs {
		// If this exact service already has an allocation, reuse it if it's not stale/running by another PID
		if alloc.ProjectName == project && alloc.Instance == instance && alloc.ServiceName == service {
			if a.isPortFree(alloc.Port) {
				a.reg.ReservePort(projID, service, alloc.Port)
				return alloc.Port, nil
			}
		}
		usedPorts[alloc.Port] = true
	}

	r, ok := DefaultRanges[serviceType]
	if !ok {
		r = DefaultRanges["other"]
	}

	for p := r.Min; p <= r.Max; p++ {
		if !usedPorts[p] && a.isPortFree(p) {
			if err := a.reg.ReservePort(projID, service, p); err != nil {
				return 0, err
			}
			return p, nil
		}
	}

	return 0, fmt.Errorf("no available ports in range %d-%d", r.Min, r.Max)
}
