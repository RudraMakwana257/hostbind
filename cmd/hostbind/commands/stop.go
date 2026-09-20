package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <service>",
	Short: "Stop a running service and release its port",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		dir, _ := os.Getwd()
		projectName := filepath.Base(dir)
		instanceName := "main"

		reg, err := registry.Open()
		if err != nil {
			fmt.Printf("Error opening registry: %v\n", err)
			os.Exit(1)
		}
		defer reg.Close()

		allocs, err := reg.GetProjectAllocations(projectName, instanceName)
		if err != nil {
			fmt.Printf("Failed to get allocations: %v\n", err)
			os.Exit(1)
		}

		var target *registry.Allocation
		for _, a := range allocs {
			if a.ServiceName == serviceName {
				target = &a
				break
			}
		}

		if target == nil {
			fmt.Printf("Service '%s' is not running in this project.\n", serviceName)
			return
		}

		if target.GetPID() > 0 {
			process, err := os.FindProcess(target.GetPID())
			if err == nil {
				// Kill the process
				if err := process.Kill(); err != nil {
					fmt.Printf("Warning: failed to kill process %d: %v\n", target.GetPID(), err)
				} else {
					fmt.Printf("🛑 Stopped service '%s' (PID: %d)\n", serviceName, target.GetPID())
				}
			}
		}

		// Release the port in the registry
		if err := reg.ReleasePort(target.Port); err != nil {
			fmt.Printf("Error releasing port: %v\n", err)
		} else {
			fmt.Printf("✅ Released port %d\n", target.Port)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
