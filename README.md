# Google Cloud MCP Server

A MCP server implementation for Google Cloud using Go and Cobra.
The server supports `stdio` as well as `SSE` as transport. The following
services and operations have been implemented:

- **Projects**: Projects List, Project Describe
- **Container**: Clusters List, Cluster Describe
- **Cloud Run**: Services List, Service Describe

## Build and Release

```bash
goreleaser build --snapshot --clean
goreleaser release --skip=publish --snapshot --clean
```

## Usage Instructions

If you want to use the tool locally, e.g. with Claude Desktop, use the following
configuration for the MCP server.

```json
{
    "mcpServers": {
      "gcloud": {
        "command": "/Users/mario-leander.reimer/Applications/gcp-mcp-server",
        "args": ["--transport", "stdio"],
        "env": {
        }
      }
    }
}
```

Alternatively, you can use the MCP introspector for easy local development:
```bash
# as stdio binary
npx @modelcontextprotocol/inspector go run main.go

# as SSE server using 
go run main.go --transport sse
npx @modelcontextprotocol/inspector npx mcp-remote@next http://localhost:8000/mcp
npx @modelcontextprotocol/inspector
```

## Maintainer

M.-Leander Reimer (@lreimer), <mario-leander.reimer@qaware.de>

## License

This software is provided under the MIT open source license, read 
the `LICENSE` file for details.
