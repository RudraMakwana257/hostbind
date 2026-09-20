package commands

import (
	"os"
	"runtime"
	"syscall"
	"time"
)

// isPIDAlive checks whether a process with the given PID is currently alive.
// Fix #6: On Unix, os.FindProcess() ALWAYS returns success (even for dead PIDs),
// so we must use signal(0) to actually verify liveness.
func isPIDAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		return true // FindProcess failure means dead on Windows
	}
	// Unix: signal(0) returns error if process doesn't exist
	return process.Signal(syscall.Signal(0)) == nil
}

// gracefulStop tries to stop a process gracefully:
//   - On Unix: sends SIGTERM, waits up to 3 seconds, then SIGKILL
//   - On Windows: calls Kill() directly (no SIGTERM equivalent)
//
// Fix #7: Previously code called Kill() immediately, bypassing cleanup handlers
// in dev servers (Next.js, Vite, etc.).
func gracefulStop(process *os.Process) {
	if runtime.GOOS == "windows" {
		// Windows has no SIGTERM, kill directly
		process.Kill()
		return
	}
	// Send SIGTERM for graceful shutdown
	process.Signal(syscall.SIGTERM)
	// Wait up to 3 seconds for the process to exit
	done := make(chan struct{})
	go func() {
		process.Wait()
		close(done)
	}()
	select {
	case <-done:
		// Process exited cleanly
	case <-time.After(3 * time.Second):
		// Force kill if it didn't exit
		process.Kill()
	}
}
