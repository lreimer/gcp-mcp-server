package cmd

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var version string
var transport string

var rootCmd = &cobra.Command{
	Use:   "gcp-mcp-server",
	Short: "A MCP server implementation for Google Cloud Platform",
	Run: func(cmd *cobra.Command, args []string) {
		// This is the entry point for the command line tool.
		// You can add your logic here to start the server or perform other actions.
		// For example, you might want to initialize a server and start listening for requests.
		// Create a new MCP server
		s := server.NewMCPServer(
			"Google Cloud Platform API",
			version,
			server.WithResourceCapabilities(true, true),
			server.WithRecovery(),
			server.WithLogging(),
		)

		// Only check for "sse" since stdio is the default
		if transport == "sse" {
			sseServer := server.NewSSEServer(s, server.WithBaseURL("http://localhost:8080"))
			log.Printf("MCP Server (SSE) listening on :8080")
			if err := sseServer.Start(":8080"); err != nil {
				log.Fatalf("MCP Server (SSE) error: %v", err)
			}
		} else {
			if err := server.ServeStdio(s); err != nil {
				log.Fatalf("MCP Server (stdio) error: %v", err)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// SetVersion set the application version to be used in the MCP server
func SetVersion(v string) {
	version = v
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&transport, "transport", "t", "stdio", "Transport to use. Valid options: stdio, see")
}
