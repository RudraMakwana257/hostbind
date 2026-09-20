package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RudraMakwana257/hostbind/internal/registry"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/spf13/cobra"
)

type EmptyArgs struct{}

type StopArgs struct {
	ServiceName string `json:"service_name"`
}

var mcpServerCmd = &cobra.Command{
	Use:    "mcp",
	Short:  "Run the HostBind MCP server over stdio for AI integration",
	Hidden: true, // Internal use for agents
	Run: func(cmd *cobra.Command, args []string) {
		transport := stdio.NewStdioServerTransport()
		server := mcp.NewServer(transport, mcp.WithName("hostbind"), mcp.WithVersion("0.3.0"))

		// Tool: hostbind_context
		server.RegisterTool("hostbind_context", "Get all running projects, services, and their URLs.", func(args EmptyArgs) (*mcp.ToolResponse, error) {
			reg, err := registry.Open()
			if err != nil {
				return nil, err
			}
			defer reg.Close()

			dir, _ := os.Getwd()
			project := filepath.Base(dir)
			allocs, err := reg.GetProjectAllocations(project, "main")
			if err != nil {
				return nil, err
			}

			services := make(map[string]string)
			for _, a := range allocs {
				services[a.ServiceName] = fmt.Sprintf("http://localhost:%d", a.Port)
			}

			out := ContextOutput{
				Project:  project,
				Instance: "main",
				Services: services,
			}

			b, _ := json.MarshalIndent(out, "", "  ")
			return mcp.NewToolResponse(mcp.NewTextContent(string(b))), nil
		})

		// Tool: hostbind_stop
		server.RegisterTool("hostbind_stop", "Stop a running service by name in the current project.", func(args StopArgs) (*mcp.ToolResponse, error) {
			reg, err := registry.Open()
			if err != nil {
				return nil, err
			}
			defer reg.Close()

			dir, _ := os.Getwd()
			project := filepath.Base(dir)

			allocs, err := reg.GetProjectAllocations(project, "main")
			if err != nil {
				return nil, err
			}

			var target *registry.Allocation
			for _, a := range allocs {
				if a.ServiceName == args.ServiceName {
					target = &a
					break
				}
			}

			if target == nil {
				return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Service '%s' is not running.", args.ServiceName))), nil
			}

			msg := ""
			if target.GetPID() > 0 {
				if process, err := os.FindProcess(target.GetPID()); err == nil {
					_ = process.Kill()
					msg += fmt.Sprintf("Stopped process PID %d. ", target.GetPID())
				}
			}

			reg.ReleasePort(target.Port)
			msg += fmt.Sprintf("Released port %d.", target.Port)

			return mcp.NewToolResponse(mcp.NewTextContent(msg)), nil
		})

		// TODO: hostbind_start can be complex because it needs to detach the process properly,
		// or AI can just run `hostbind run` directly via bash.
		
		if err := server.Serve(); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mcpServerCmd)
}
