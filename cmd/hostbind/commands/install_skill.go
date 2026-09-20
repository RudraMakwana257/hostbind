package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var installSkillCmd = &cobra.Command{
	Use:   "install-skill",
	Short: "Inject AI agent skill instructions into your project and global agent settings",
	Run: func(cmd *cobra.Command, args []string) {
		const skillContent = `
# HostBind Agent Skill
Always use 'hostbind run --name <service>' to start servers.
Never hardcode or guess ports. Run 'hostbind context --json' to discover active ports.
Do not modify vite.config.ts or next.config.js to hardcode ports.
`
		// 1. Install locally into project's AGENTS.md
		content, err := os.ReadFile("AGENTS.md")
		if err == nil && strings.Contains(string(content), "HostBind Agent Skill") {
			fmt.Println("ℹ️  HostBind skill is already present in AGENTS.md")
		} else {
			f, err := os.OpenFile("AGENTS.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("❌ Failed to write AGENTS.md: %v\n", err)
			} else {
				defer f.Close()
				f.WriteString(skillContent)
				fmt.Println("✅ Successfully injected HostBind rules into AGENTS.md")
			}
		}

		// 2. Install globally for all known AI agent frameworks
		// Fix #10: Previously hardcoded to Antigravity path only.
		// Now installs to multiple agent skill directories for broad compatibility.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Warning: could not find home directory: %v\n", err)
			return
		}

		fullSkillContent := `---
name: hostbind-local-development
description: Use HostBind for local port management and zero-conflict multi-project setups.
---
` + skillContent

		// Agent skill directories to install into
		agentDirs := []struct {
			name string
			path string
		}{
			{
				name: "Antigravity CLI",
				path: filepath.Join(home, ".gemini", "antigravity", "skills", "hostbind"),
			},
			{
				name: "Gemini CLI (legacy)",
				path: filepath.Join(home, ".gemini", "skills", "hostbind"),
			},
		}

		// On Windows/Mac: also try Cursor rules directory
		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			agentDirs = append(agentDirs, struct {
				name string
				path string
			}{
				name: "Cursor (.cursor/rules)",
				path: filepath.Join(home, ".cursor", "rules"),
			})
		}

		for _, agentDir := range agentDirs {
			if err := os.MkdirAll(agentDir.path, 0755); err != nil {
				fmt.Printf("⚠️  Could not create dir for %s: %v\n", agentDir.name, err)
				continue
			}
			skillPath := filepath.Join(agentDir.path, "SKILL.md")
			if err := os.WriteFile(skillPath, []byte(fullSkillContent), 0644); err != nil {
				fmt.Printf("⚠️  Could not write skill for %s: %v\n", agentDir.name, err)
			} else {
				fmt.Printf("✅ Installed HostBind skill for %s\n", agentDir.name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installSkillCmd)
}
