package commands

import (
	"fmt"
	"os"

	"path/filepath"

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
			&adapters.PythonAdapter{},
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
				os.Exit(1)
			}
		}

		fmt.Printf("🚀 HostBind: Starting '%s' on port %d using %s adapter\n", name, port, activeAdapter.Name())
		fmt.Printf("💻 Command: %v\n", runArgs)

		// 5. Run process
		pid, err := r.Start(runArgs, project, instance, name)
		if err != nil {
			fmt.Printf("Error starting process: %v\n", err)
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

		fmt.Printf("✅ Running with PID %d\n", pid)
		// TODO: Wait for process to exit and handle graceful shutdown
	},
}

func init() {
	runCmd.Flags().String("name", "", "Name of the service (e.g. api, web)")
	rootCmd.AddCommand(runCmd)
}
