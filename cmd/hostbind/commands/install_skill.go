package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var installSkillCmd = &cobra.Command{
	Use:   "install-skill",
	Short: "Inject AGENTS.md instructions for AI coding assistants",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		agentsFile := filepath.Join(dir, "AGENTS.md")

		content := `> This repo uses HostBind. 
> Never choose or assume dev ports (e.g., 3000, 5173, 8000).
> 
> - Run 'hostbind context --json' for service URLs.
> - Start services with 'hostbind run --name <service>'.
> - Stop services with 'hostbind stop <service>'.
> - If something is unreachable, DO NOT guess the port. Ask HostBind.
`

		// Check if file exists to append or create new
		f, err := os.OpenFile(agentsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open AGENTS.md: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		if _, err := f.WriteString(content); err != nil {
			fmt.Printf("Failed to write to AGENTS.md: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Successfully injected HostBind rules into AGENTS.md")
		fmt.Println("AI coding agents (Claude, Cursor, etc.) will now automatically read this file and never guess ports again.")
	},
}

func init() {
	rootCmd.AddCommand(installSkillCmd)
}
