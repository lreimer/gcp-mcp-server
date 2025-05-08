package gcp

import (
	"context"

	container "cloud.google.com/go/container/apiv1"
	containerpb "cloud.google.com/go/container/apiv1/containerpb"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Adds all the supported GCP container services as tool to the server
func AddContainerTools(s *server.MCPServer) {
	clusterList(s)
	clusterDescribe(s)
}

func clusterDescribe(s *server.MCPServer) {
	// create a new MCP tool for describing Sonar projects
	projectsTool := mcp.NewTool("cluster_describe",
		mcp.WithDescription("Get and describe a GKE Kubernetes cluster."),
		mcp.WithString("name",
			mcp.Description("The name of the cluster to describe."),
			mcp.Required(),
		),
		mcp.WithString("project",
			mcp.Description("The GCP project name."),
			mcp.Required(),
		),
		mcp.WithString("location",
			mcp.Description("Compute zone or region (e.g. europe-west4 or europe-north1) for the clusters."),
			mcp.Required(),
		),
	)

	// add the tool to the server
	s.AddTool(projectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the location parameter from the request
		project := request.Params.Arguments["project"].(string)
		location := request.Params.Arguments["location"].(string)
		name := request.Params.Arguments["name"].(string)

		c, err := container.NewClusterManagerClient(context.Background())
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to create cluster manager client.", err), err
		}
		defer c.Close()

		req := &containerpb.GetClusterRequest{
			Name: "projects/" + project + "/locations/" + location + "/clusters/" + name,
		}
		resp, err := c.GetCluster(ctx, req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to describe cluster.", err), err
		}

		return newToolResult(resp)
	})
}

func clusterList(s *server.MCPServer) {
	// create a new MCP tool for listing Sonar projects
	projectsTool := mcp.NewTool("cluster_list",
		mcp.WithDescription("List existing GKE Kubernetes clusters with running containers."),
		mcp.WithString("project",
			mcp.Description("The GCP project name. '*' matches all projects."),
			mcp.Required(),
		),
		mcp.WithString("location",
			mcp.Description("Compute zone or region (e.g. europe-west4 or europe-north1) for the clusters. '*' matches all locations."),
			mcp.Required(),
		),
	)

	// add the tool to the server
	s.AddTool(projectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the location parameter from the request
		project := request.Params.Arguments["project"].(string)
		location := request.Params.Arguments["location"].(string)

		c, err := container.NewClusterManagerClient(context.Background())
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to create cluster manager client.", err), err
		}
		defer c.Close()

		req := &containerpb.ListClustersRequest{
			// see https://pkg.go.dev/cloud.google.com/go/container/apiv1/containerpb#ListClustersRequest.
			Parent: "projects/" + project + "/locations/" + location,
		}
		resp, err := c.ListClusters(ctx, req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to list clusters.", err), err
		}

		// Check if the response is empty
		if len(resp.Clusters) == 0 {
			return mcp.NewToolResultText("No clusters found."), nil
		}

		// marshal the response to JSON
		return newToolResult(resp.Clusters)
	})
}
