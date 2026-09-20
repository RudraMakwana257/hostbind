package main

import (
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type TestArgs struct {
	Name string `json:"name"`
}

func main() {
	transport := stdio.NewStdioServerTransport()
	server := mcp.NewServer(transport)
	err := server.RegisterTool("test", "test desc", func(args TestArgs) (*mcp.ToolResponse, error) {
		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Hello %s", args.Name))), nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Success")
}
