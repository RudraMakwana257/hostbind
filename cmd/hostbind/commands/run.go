package commands

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/RudraMakwana257/hostbind/internal/adapters"
	"github.com/RudraMakwana257/hostbind/internal/allocator"
	"github.com/RudraMakwana257/hostbind/internal/envgen"
	"github.com/RudraMakwana257/hostbind/internal/registry"
	"github.com/RudraMakwana257/hostbind/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [--name service] -- <command>",
	Short: "Allocate a port and run a service",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "web" // default service name
		}

		// 1. Open Registry
		reg, err := registry.Open()
		if err != nil {
			fmt.Printf("Error opening registry: %v\n", err)
			os.Exit(1)
		}
		defer reg.Close()

		// 2. Detect framework via Adapters
		dir, _ := os.Getwd()
		var activeAdapter adapters.Adapter = &adapters.GenericAdapter{} // default

		availableAdapters := []adapters.Adapter{
			&adapters.ViteAdapter{},
			&adapters.NextAdapter{},
			&adapters.DjangoAdapter{}, // Django before Python (higher specificity)
			&adapters.PythonAdapter{},
			&adapters.ExpressAdapter{},
		}

		highestScore := 0
		for _, a := range availableAdapters {
			if matched, score := a.Detect(dir); matched && score > highestScore {
				activeAdapter = a
				highestScore = score
			}
		}

		// 3. Allocate Port
		alloc := allocator.New(reg)
		project := filepath.Base(dir) // simple project extraction for v0.1
		instance := "main"

		port, err := alloc.Allocate(project, instance, name, "web")
		if err != nil {
			fmt.Printf("Error allocating port: %v\n", err)
			os.Exit(1)
		}

		// 4. Setup Runner
		portConfig := activeAdapter.PortArgs(port)
		r := runner.New(port, portConfig)

		// Command arguments & Security Check
		runArgs := args
		isCustomCommand := false
		if len(runArgs) == 0 {
			runArgs = activeAdapter.DefaultCommand()
		} else {
			isCustomCommand = true
		}

		if isCustomCommand {
			// Security: Do not allow AI agents to run arbitrary dangerous commands
			fmt.Printf("⚠️  Security Alert: An attempt is being made to run a custom command:\n")
			fmt.Printf("   > %v\n", runArgs)
			fmt.Printf("Allow this execution? [y/N]: ")

			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" && response != "yes" {
				fmt.Println("❌ Execution cancelled by user.")
				reg.ReleasePort(port)
				os.Exit(1)
			}
		}

		fmt.Printf("🚀 HostBind: Starting '%s' on port %d using %s adapter\n", name, port, activeAdapter.Name())
		fmt.Printf("💻 Command: %v\n", runArgs)

		// 5. Run process
		pid, err := r.Start(runArgs, project, instance, name)
		if err != nil {
			fmt.Printf("Error starting process: %v\n", err)
			reg.ReleasePort(port)
			os.Exit(1)
		}

		if err := reg.UpdatePID(port, pid); err != nil {
			fmt.Printf("Warning: failed to update PID in registry: %v\n", err)
		}

		// 6. Generate Env
		envVars := map[string]string{
			"HOSTBIND_PROJECT":  project,
			"HOSTBIND_INSTANCE": instance,
			"HOSTBIND_SERVICE":  name,
			"HOSTBIND_PORT":     fmt.Sprintf("%d", port),
			"HOSTBIND_URL":      fmt.Sprintf("http://localhost:%d", port),
		}

		// Optional: add adapter specific envs
		for k, v := range portConfig.Env {
			envVars[k] = v
		}

		if err := envgen.GenerateSync(dir, envVars); err != nil {
			fmt.Printf("Warning: failed to sync .env.hostbind: %v\n", err)
		}

		home, _ := os.UserHomeDir()
		fmt.Printf("✅ Running with PID %d\n", pid)
		fmt.Printf("📄 Logs: %s/.hostbind/logs/%s-%s-%s.log\n", home, project, instance, name)
		fmt.Printf("ℹ️  Press Ctrl+C to stop and release port %d\n", port)

		// Fix #5: Wait for process to exit and auto-release port on exit.
		// This prevents stale "running" entries in the registry when the process
		// crashes or is stopped via Ctrl+C.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Wait for signal or process to exit
		done := make(chan struct{})
		go func() {
			// We can't directly wait on a PID we didn't exec ourselves via cmd.Wait,
			// so we poll via signal(0) every second until it dies.
			// In future: refactor Runner to return *exec.Cmd for proper Wait().
			for {
				select {
				case <-done:
					return
				default:
					if !isPIDAlive(pid) {
						fmt.Printf("\n🛑 Service '%s' (PID %d) exited. Releasing port %d...\n", name, pid, port)
						reg.ReleasePort(port)
						close(done)
						return
					}
					// Check every second
					select {
					case <-done:
						return
					}
				}
			}
		}()

		select {
		case <-quit:
			fmt.Printf("\n🛑 Stopping '%s' (PID %d)...\n", name, pid)
			p, err := os.FindProcess(pid)
			if err == nil {
				// Send SIGTERM first for graceful shutdown, then SIGKILL
				gracefulStop(p)
			}
			reg.ReleasePort(port)
			fmt.Printf("✅ Port %d released.\n", port)
		case <-done:
			// Process already exited (handled in goroutine above)
		}
	},
}

func init() {
	runCmd.Flags().String("name", "", "Name of the service (e.g. api, web)")
	rootCmd.AddCommand(runCmd)
}
