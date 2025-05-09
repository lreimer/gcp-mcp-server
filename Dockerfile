FROM golang:1.24.3-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gcp-mcp-server -ldflags "-s -w -X main.version=$(date +%Y-%m-%dT%H:%M:%S%z)"

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/gcp-mcp-server /gcp-mcp-server

EXPOSE 8000
CMD ["/gcp-mcp-server", "-t", "sse", "-p", "8000"]
