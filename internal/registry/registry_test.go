package registry_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

// openTestRegistry creates a registry backed by a temp file for isolated testing.
func openTestRegistry(t *testing.T) (*registry.Registry, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	// Point registry home to temp dir
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir) // Windows

	reg, err := registry.Open()
	if err != nil {
		t.Fatalf("failed to open test registry: %v", err)
	}
	cleanup := func() {
		reg.Close()
		os.RemoveAll(filepath.Join(tmpDir, ".hostbind"))
	}
	return reg, cleanup
}

func TestRegistryOpenAndClose(t *testing.T) {
	reg, cleanup := openTestRegistry(t)
	defer cleanup()
	if reg == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestEnsureProject(t *testing.T) {
	reg, cleanup := openTestRegistry(t)
	defer cleanup()

	id, err := reg.EnsureProject("myapp", "main", "/home/user/myapp")
	if err != nil {
		t.Fatalf("EnsureProject failed: %v", err)
	}
	if id <= 0 {
		t.Errorf("expected positive project ID, got %d", id)
	}

	// Calling again should return the same ID (idempotent)
	id2, err := reg.EnsureProject("myapp", "main", "/home/user/myapp")
	if err != nil {
		t.Fatalf("EnsureProject (second call) failed: %v", err)
	}
	if id != id2 {
		t.Errorf("expected same ID on second call: got %d and %d", id, id2)
	}
}

func TestReserveAndReleasePort(t *testing.T) {
	reg, cleanup := openTestRegistry(t)
	defer cleanup()

	projectID, _ := reg.EnsureProject("testapp", "main", "/tmp/testapp")

	// Reserve a port
	err := reg.ReservePort(projectID, "web", 4300)
	if err != nil {
		t.Fatalf("ReservePort failed: %v", err)
	}

	// Should appear in active allocations
	allocs, err := reg.GetActiveAllocations()
	if err != nil {
		t.Fatalf("GetActiveAllocations failed: %v", err)
	}
	found := false
	for _, a := range allocs {
		if a.Port == 4300 && a.ServiceName == "web" {
			found = true
			if a.Status != "reserved" {
				t.Errorf("expected status 'reserved', got %q", a.Status)
			}
		}
	}
	if !found {
		t.Error("expected port 4300 to appear in active allocations after ReservePort")
	}

	// Release the port
	err = reg.ReleasePort(4300)
	if err != nil {
		t.Fatalf("ReleasePort failed: %v", err)
	}

	// Should no longer appear
	allocs, _ = reg.GetActiveAllocations()
	for _, a := range allocs {
		if a.Port == 4300 {
			t.Error("port 4300 should have been removed from allocations after ReleasePort")
		}
	}
}

func TestUpdatePID(t *testing.T) {
	reg, cleanup := openTestRegistry(t)
	defer cleanup()

	projectID, _ := reg.EnsureProject("testapp", "main", "/tmp/testapp")
	_ = reg.ReservePort(projectID, "api", 5000)

	// Update PID → marks as running
	err := reg.UpdatePID(5000, 12345)
	if err != nil {
		t.Fatalf("UpdatePID failed: %v", err)
	}

	allocs, _ := reg.GetProjectAllocations("testapp", "main")
	for _, a := range allocs {
		if a.Port == 5000 {
			if a.Status != "running" {
				t.Errorf("expected status 'running' after UpdatePID, got %q", a.Status)
			}
			if a.GetPID() != 12345 {
				t.Errorf("expected PID 12345, got %d", a.GetPID())
			}
		}
	}
}

func TestGetProjectAllocations(t *testing.T) {
	reg, cleanup := openTestRegistry(t)
	defer cleanup()

	id1, _ := reg.EnsureProject("app-a", "main", "/tmp/app-a")
	id2, _ := reg.EnsureProject("app-b", "main", "/tmp/app-b")

	_ = reg.ReservePort(id1, "web", 3001)
	_ = reg.ReservePort(id2, "web", 3002)

	// Only app-a's ports
	allocs, err := reg.GetProjectAllocations("app-a", "main")
	if err != nil {
		t.Fatalf("GetProjectAllocations failed: %v", err)
	}
	if len(allocs) != 1 {
		t.Errorf("expected 1 allocation for app-a, got %d", len(allocs))
	}
	if allocs[0].Port != 3001 {
		t.Errorf("expected port 3001 for app-a, got %d", allocs[0].Port)
	}
}

func TestAllocationGetPID_NullSafe(t *testing.T) {
	a := &registry.Allocation{}
	// PID is unset (NULL in DB) — GetPID should return 0, not panic
	if pid := a.GetPID(); pid != 0 {
		t.Errorf("expected GetPID() == 0 for null PID, got %d", pid)
	}
}
