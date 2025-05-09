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
npx @modelcontextprotocol/inspector npx mcp-remote@next http://localhost:8000/sse
npx @modelcontextprotocol/inspector
```

## Deployment

Currently using manual Google Cloud Run deployment. Can either be deployed
directly from source or using the Docker image built on Github.

```bash
# make sure to enable all required APIs
gcloud services enable secretmanager.googleapis.com
gcloud services enable cloudbuild.googleapis.com artifactregistry.googleapis.com
gcloud services enable run.googleapis.com 

# create Google service account with required permissions and create key file
export PROJECT_ID=$(gcloud config get-value project)
export SA_NAME=gcp-mcp-server-sa

gcloud iam service-accounts create $SA_NAME --display-name="GCP MCP Server Service Account"
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/editor"
gcloud iam service-accounts keys create $SA_NAME.json \
  --iam-account=$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com

# create the keyfile as a secret and attach IAM policy
gcloud secrets create $SA_NAME --data-file=$SA_NAME.json
gcloud secrets add-iam-policy-binding $SA_NAME \
  --member=serviceAccount:343509396461-compute@developer.gserviceaccount.com \
  --role=roles/secretmanager.secretAccessor

gcloud run deploy gcp-mcp-server --source=. \
  --region=europe-north1 \
  --port=8000 --allow-unauthenticated \
  --set-secrets=/secrets/gcp-mcp-server-sa.json=gcp-mcp-server-sa:latest \
  --set-env-vars=GOOGLE_APPLICATION_CREDENTIALS=/secrets/gcp-mcp-server-sa.json,BASE_URL=https://gcp-mcp-server-343509396461.europe-north1.run.app

gcloud run services delete gcp-mcp-server --async --region=europe-north1
```

## Maintainer

M.-Leander Reimer (@lreimer), <mario-leander.reimer@qaware.de>

## License

This software is provided under the MIT open source license, read 
the `LICENSE` file for details.
