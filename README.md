# Google Cloud MCP Server

A MCP server implementation for Google Cloud using Go and Cobra.
The server supports `stdio` as well as `SSE` as transport.

## Build and Deploy

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

## Maintainer

M.-Leander Reimer (@lreimer), <mario-leander.reimer@qaware.de>

## License

This software is provided under the MIT open source license, read 
the `LICENSE` file for details.
