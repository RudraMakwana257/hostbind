package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/RudraMakwana257/hostbind/internal/registry"
	"github.com/spf13/cobra"
)

type ContextOutput struct {
	Project  string            `json:"project"`
	Instance string            `json:"instance"`
	Services map[string]string `json:"services"`
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Print context for AI agents (use --json for agent consumption)",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")

		reg, err := registry.Open()
		if err != nil {
			fmt.Printf("Error opening registry: %v\n", err)
			os.Exit(1)
		}
		defer reg.Close()

		// TODO: Query the current directory's project/instance from DB.
		out := ContextOutput{
			Project:  "hostbind",
			Instance: "main",
			Services: map[string]string{
				"web": "http://localhost:4300",
			},
		}

		if jsonFlag {
			b, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Printf("Project: %s (Instance: %s)\n\n", out.Project, out.Instance)
			for svc, url := range out.Services {
				fmt.Printf("  %s -> %s\n", svc, url)
			}
		}
	},
}

func init() {
	contextCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(contextCmd)
}
