package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hostbind",
	Short: "HostBind is a local-first service registry and port allocator",
	Long: `HostBind prevents port conflicts for AI agents and developers.
It acts as the single source of truth for where local services are running.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags, e.g., --json output
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
}
