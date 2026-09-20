package commands

import (
	"fmt"
	"os"
	"path/filepath"
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

		// 2. Install globally for Antigravity/Cursor/Agents
		home, err := os.UserHomeDir()
		if err == nil {
			skillDir := filepath.Join(home, ".gemini", "antigravity-cli", "skills", "hostbind-local-development")
			if err := os.MkdirAll(skillDir, 0755); err == nil {
				skillPath := filepath.Join(skillDir, "SKILL.md")
				
				// Full skill content
				fullSkill := `---
name: hostbind-local-development
description: Use HostBind for local port management and zero-conflict multi-project setups.
---
` + skillContent

				os.WriteFile(skillPath, []byte(fullSkill), 0644)
				fmt.Println("✅ Successfully installed global HostBind skill for AI agents")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installSkillCmd)
}
