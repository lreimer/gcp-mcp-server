package gcp

import (
	"context"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/api/iterator"
)

// Adds all the supported GCP project tools to the server
func AddProjectTools(s *server.MCPServer) {
	projectsList(s)
	projectDescribe(s)
}

func projectDescribe(s *server.MCPServer) {
	// create a new MCP tool for describing a Google Cloud Project
	projectTool := mcp.NewTool("project_describe",
		mcp.WithDescription("Get and describe a Google Cloud Project."),
		mcp.WithString("name",
			mcp.Description("The name of the project to describe (either project ID or project number)."),
			mcp.DefaultString(Project),
			mcp.Required(),
		),
	)

	// add the tool to the server
	s.AddTool(projectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the parameter from the request
		name := request.Params.Arguments["name"].(string)

		// Create a new Project Client
		c, err := resourcemanager.NewProjectsClient(ctx)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to create resource manager client.", err), err
		}
		defer c.Close()

		req := &resourcemanagerpb.GetProjectRequest{
			Name: "projects/" + name,
		}
		resp, err := c.GetProject(ctx, req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to describe project.", err), err
		}

		return newToolResult(resp)
	})
}

func projectsList(s *server.MCPServer) {
	// create a new MCP tool for listing Google Cloud Projects
	projectsTool := mcp.NewTool("projects_list",
		mcp.WithDescription("List existing Google Cloud Projects."),
		mcp.WithString("organization",
			mcp.Description("The GCP organization ID."),
			mcp.DefaultString(Organization),
			mcp.Required(),
		),
	)

	// add the tool to the server
	s.AddTool(projectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the optional filter parameter from the request
		orgId := request.Params.Arguments["organization"].(string)

		// Create a new Project Client
		c, err := resourcemanager.NewProjectsClient(ctx)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to create resource manager client.", err), err
		}
		defer c.Close()

		req := &resourcemanagerpb.ListProjectsRequest{
			Parent: "organizations/" + orgId,
		}

		// Use iterator to collect all projects
		it := c.ListProjects(ctx, req)
		var projects []*resourcemanagerpb.Project

		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return mcp.NewToolResultErrorFromErr("unable to list projects.", err), err
			}
			projects = append(projects, resp)
		}

		// Check if the response is empty
		if len(projects) == 0 {
			return mcp.NewToolResultText("No projects found."), nil
		}

		return newToolResult(projects)
	})
}
