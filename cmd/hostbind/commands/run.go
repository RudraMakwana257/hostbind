package commands

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/adapters"
	"github.com/RudraMakwana257/hostbind/internal/allocator"
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

		// Command arguments
		runArgs := args
		if len(runArgs) == 0 {
			runArgs = activeAdapter.DefaultCommand()
		}

		fmt.Printf("🚀 HostBind: Starting '%s' on port %d using %s adapter\n", name, port, activeAdapter.Name())
		fmt.Printf("💻 Command: %v\n", runArgs)

		// 5. Run process
		pid, err := r.Start(runArgs, project, instance, name)
		if err != nil {
			fmt.Printf("Error starting process: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Running with PID %d\n", pid)
		// TODO: Wait for process to exit and handle graceful shutdown
	},
}

func init() {
	runCmd.Flags().String("name", "", "Name of the service (e.g. api, web)")
	rootCmd.AddCommand(runCmd)
}
