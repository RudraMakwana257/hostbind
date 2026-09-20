package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port <service>",
	Short: "Print the allocated port for a given service (useful for scripting)",
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

		for _, a := range allocs {
			if a.ServiceName == serviceName {
				fmt.Printf("%d\n", a.Port)
				return
			}
		}

		fmt.Fprintf(os.Stderr, "Service '%s' not found or not running\n", serviceName)
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
}
