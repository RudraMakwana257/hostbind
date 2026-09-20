package commands

import (
	"fmt"
	"net/http"
	"os"

	"github.com/RudraMakwana257/hostbind/internal/proxy"
	"github.com/RudraMakwana257/hostbind/internal/registry"
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start the reverse proxy to enable <service>.<project>.localhost routing",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")

		reg, err := registry.Open()
		if err != nil {
			fmt.Printf("Error opening registry: %v\n", err)
			os.Exit(1)
		}
		defer reg.Close()

		p := proxy.New(reg)
		addr := fmt.Sprintf("127.0.0.1:%d", port)

		fmt.Printf("🌐 HostBind Reverse Proxy started on http://%s\n", addr)
		fmt.Printf("   Try accessing: http://<service>.<project>.localhost:%d\n", port)
		fmt.Printf("   Example: http://web.my-app.localhost:%d\n\n", port)

		if err := http.ListenAndServe(addr, p); err != nil {
			fmt.Fprintf(os.Stderr, "Proxy server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	proxyCmd.Flags().Int("port", 8080, "Port for the proxy server to listen on")
	rootCmd.AddCommand(proxyCmd)
}
