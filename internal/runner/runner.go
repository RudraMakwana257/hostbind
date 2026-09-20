package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/adapters"
)

type Runner struct {
	port   int
	config adapters.PortConfig
}

func New(port int, config adapters.PortConfig) *Runner {
	return &Runner{
		port:   port,
		config: config,
	}
}

// Start launches the command with injected environment variables and arguments.
// It returns the process ID (PID) and an error if it fails to start.
func (r *Runner) Start(cmdArgs []string, project, instance, service string) (int, error) {
	if len(cmdArgs) == 0 {
		return 0, fmt.Errorf("no command provided")
	}

	// Append adapter-specific arguments (e.g. --port 4300)
	finalArgs := append(cmdArgs[1:], r.config.Args...)
	
	cmd := exec.Command(cmdArgs[0], finalArgs...)
	logDir := filepath.Join(os.Getenv("HOME"), ".hostbind", "logs")
	os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, fmt.Sprintf("%s-%s-%s.log", project, instance, service))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Inherit current environment and inject HostBind vars
	env := os.Environ()
	env = append(env, fmt.Sprintf("HOSTBIND_PROJECT=%s", project))
	env = append(env, fmt.Sprintf("HOSTBIND_INSTANCE=%s", instance))
	env = append(env, fmt.Sprintf("HOSTBIND_SERVICE=%s", service))
	env = append(env, fmt.Sprintf("HOSTBIND_PORT=%d", r.port))
	env = append(env, fmt.Sprintf("HOSTBIND_URL=http://localhost:%d", r.port))

	// Inject adapter-specific env vars (e.g. PORT=4300)
	for k, v := range r.config.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start process: %w", err)
	}

	return cmd.Process.Pid, nil
}
