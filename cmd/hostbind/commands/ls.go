package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/RudraMakwana257/hostbind/internal/registry"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all running projects and services",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := registry.Open()
		if err != nil {
			fmt.Printf("Error opening registry: %v\n", err)
			os.Exit(1)
		}
		defer reg.Close()

		// TODO: Query from SQLite allocations table.
		// For v0.1 spike, we will just print the columns.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "PROJECT\tINSTANCE\tSERVICE\tPORT\tSTATUS\tPID")
		
		// Mock data for display until DB query is wired
		fmt.Fprintln(w, "hostbind\tmain\tweb\t4300\trunning\t12345")
		
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
