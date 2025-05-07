package main

import "github.com/lreimer/gcp-mcp-server/cmd"

var version string

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
