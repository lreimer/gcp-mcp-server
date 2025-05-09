package gcp

import (
	"context"

	run "cloud.google.com/go/run/apiv2"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/api/iterator"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Adds all the supported GCP Cloud Run services as tool to the server
func AddCloudRunTools(s *server.MCPServer) {
	servicesList(s)
	serviceDescribe(s)
}

func serviceDescribe(s *server.MCPServer) {
	// create a new MCP tool for describing Cloud Run services
	serviceTool := mcp.NewTool("run_service_describe",
		mcp.WithDescription("Get and describe a Google Cloud Run service."),
		mcp.WithString("name",
			mcp.Description("The name of the Cloud Run service to describe."),
			mcp.Required(),
		),
		mcp.WithString("project",
			mcp.Description("The GCP project name."),
			mcp.Required(),
		),
		mcp.WithString("location",
			mcp.Description("Region (e.g. europe-west1) for the service."),
			mcp.Required(),
		),
	)

	// add the tool to the server
	s.AddTool(serviceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the parameters from the request
		project := request.Params.Arguments["project"].(string)
		location := request.Params.Arguments["location"].(string)
		name := request.Params.Arguments["name"].(string)

		// Create a new Cloud Run client
		c, err := run.NewServicesClient(ctx)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to create Cloud Run services client.", err), err
		}
		defer c.Close()

		// Create the request to get the service
		req := &runpb.GetServiceRequest{
			Name: "projects/" + project + "/locations/" + location + "/services/" + name,
		}

		// Call the API
		resp, err := c.GetService(ctx, req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to describe service.", err), err
		}

		// Return the result
		return newToolResult(resp)
	})
}

func servicesList(s *server.MCPServer) {
	// create a new MCP tool for listing Cloud Run services
	serviceTool := mcp.NewTool("run_services_list",
		mcp.WithDescription("List existing Google Cloud Run services."),
		mcp.WithString("project",
			mcp.Description("The GCP project name."),
			mcp.DefaultString(Project),
			mcp.Required(),
		),
		mcp.WithString("location",
			mcp.Description("Region (e.g. europe-west1) for the services. Use a specific region, as Cloud Run is regional."),
			mcp.DefaultString(Location),
			mcp.Required(),
		),
	)

	// add the tool to the server
	s.AddTool(serviceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the parameters from the request
		project := request.Params.Arguments["project"].(string)
		location := request.Params.Arguments["location"].(string)

		// Create a new Cloud Run client
		c, err := run.NewServicesClient(ctx)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to create Cloud Run services client.", err), err
		}
		defer c.Close()

		// Create the request to list services
		req := &runpb.ListServicesRequest{
			Parent: "projects/" + project + "/locations/" + location,
		}

		// Use iterator to collect all projects
		it := c.ListServices(ctx, req)
		var services []*runpb.Service

		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}

			// add service to the list
			services = append(services, resp)
		}

		// Check if the response is empty
		if len(services) == 0 {
			return mcp.NewToolResultText("No Cloud Run services found."), nil
		}

		// Return the result
		return newToolResult(services)
	})
}
