package cmd

import (
	"log"
	"os"

	"github.com/lreimer/gcp-mcp-server/gcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var version string
var capabilities []string
var transport string
var baseURL string
var port string

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

		// Add the capabilities to the server
		for _, cap := range capabilities {
			if cap == "container" || cap == "all" {
				gcp.AddContainerTools(s)
			}
			if cap == "run" || cap == "all" {
				gcp.AddCloudRunTools(s)
			}
			if cap == "project" || cap == "all" {
				gcp.AddProjectTools(s)
			}
		}

		// Only check for "sse" since stdio is the default
		if transport == "sse" {
			sseServer := server.NewSSEServer(s, server.WithBaseURL(baseURL))
			ssePort := "0.0.0.0:" + port
			log.Printf("MCP Server (SSE) listening on %s", ssePort)
			if err := sseServer.Start(ssePort); err != nil {
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
	rootCmd.Flags().StringArrayVarP(&capabilities, "capabilities", "c", []string{"all"}, "The capabilities to use. Valid options: all, container, run, project")
	rootCmd.Flags().StringVarP(&transport, "transport", "t", "stdio", "Transport to use. Valid options: stdio, sse")
	rootCmd.Flags().StringVarP(&baseURL, "url", "u", "http://localhost:8000", "The public SSE base URL to use.")
	rootCmd.Flags().StringVarP(&port, "port", "p", "8000", "The local SSE server port to use.")
}
