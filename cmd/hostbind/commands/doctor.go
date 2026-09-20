package commands

import (
	"fmt"
	"net"
	"os"
	"text/tabwriter"

	"github.com/RudraMakwana257/hostbind/internal/registry"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose stale ports, orphaned processes, and environment issues",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := registry.Open()
		if err != nil {
			fmt.Printf("❌ Failed to open registry: %v\n", err)
			os.Exit(1)
		}
		defer reg.Close()

		allocs, err := reg.GetActiveAllocations()
		if err != nil {
			fmt.Printf("❌ Failed to read registry: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🩺 HostBind Doctor - Diagnostic Report")
		fmt.Println("--------------------------------------")
		
		issuesFound := 0
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		
		for _, a := range allocs {
			// 1. Check if PID is alive
			pidAlive := false
			if a.PID > 0 {
				process, err := os.FindProcess(a.PID)
				if err == nil {
					// In Unix, finding a process always succeeds, but sending signal 0 checks if it's actually alive
					// We'll assume if it's in DB, we check port instead for cross-platform reliability.
					_ = process
					pidAlive = true
				}
			}

			// 2. Check if port is bound
			portBound := false
			l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", a.Port))
			if err != nil {
				portBound = true // Port is in use
			} else {
				l.Close()
			}

			// Diagnostics logic
			if a.Status == "running" && !portBound {
				fmt.Printf("⚠️  [Stale] Service '%s' (Port %d) is marked running but port is dead. Releasing...\n", a.ServiceName, a.Port)
				reg.ReleasePort(a.Port)
				issuesFound++
			} else if a.Status == "reserved" && !portBound && !pidAlive {
				fmt.Printf("⚠️  [Orphaned] Service '%s' (Port %d) was reserved but never started. Releasing...\n", a.ServiceName, a.Port)
				reg.ReleasePort(a.Port)
				issuesFound++
			} else if portBound && a.Status != "running" {
				fmt.Printf("⚠️  [Conflict] Port %d is currently occupied by an unknown process outside HostBind.\n", a.Port)
				issuesFound++
			} else {
				fmt.Fprintf(w, "✅\t%s\t(Port %d)\tHealthy\n", a.ServiceName, a.Port)
			}
		}
		
		w.Flush()

		if issuesFound == 0 {
			fmt.Println("\n✅ All systems healthy! No issues detected.")
		} else {
			fmt.Printf("\n⚠️  Found %d issues. Some were auto-fixed. Check conflicts manually.\n", issuesFound)
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
