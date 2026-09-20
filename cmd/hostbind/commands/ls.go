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

		allocs, err := reg.GetActiveAllocations()
		if err != nil {
			fmt.Printf("Failed to get allocations: %v\n", err)
			os.Exit(1)
		}

		if len(allocs) == 0 {
			fmt.Println("No services registered. Run 'hostbind run --name <service>' to start one.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "PROJECT\tINSTANCE\tSERVICE\tPORT\tSTATUS\tPID")

		for _, a := range allocs {
			// Fix #16: Show "-" instead of "0" for reserved services that haven't
			// started yet (PID is NULL in DB, GetPID() returns 0).
			pid := "-"
			if a.GetPID() > 0 {
				pid = fmt.Sprintf("%d", a.GetPID())
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
				a.ProjectName, a.Instance, a.ServiceName, a.Port, a.Status, pid)
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
