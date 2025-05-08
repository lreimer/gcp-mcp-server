FROM gcr.io/distroless/static-debian12

COPY dist/gcp-mcp-server_linux_amd64_v1/gcp-mcp-server /
COPY gcp-mcp-server-sa.json /

EXPOSE 8080
ENV GOOGLE_APPLICATION_CREDENTIALS=gcp-mcp-server-sa.json

CMD ["/gcp-mcp-server", "-t", "sse", "-p", "8080"]