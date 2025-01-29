# Sample Adapter

This repository contains a sample adapter for ingesting data from a SCIM 2.0 server.

## Code Structure

- `pkg/scim`: Contains the implementation of the adapter.
- Other modules in `pkg` contain utility functions.
- `cmd/adapter/main.go`: Responsible for running all adapters defined within `pkg`. Ensure to call `RegisterAdapter` for any new adapter added.

## Build

### Building a Docker Image

To build the Docker image for `sample-adapter`, run the following command from the root of the repository:

```bash
docker build -t sample-adapter:your-tag-here .
```

To run the container after building, use:

```bash
docker run -d --name sample-adapter sample-adapter:your-tag-here
```

### Building a Binary

To build and run the adapter as a binary, use the following commands:

```bash
go build -o sample-adapter ./cmd/adapter
./sample-adapter
```

If you encounter a permission error, make the binary executable:

```bash
chmod +x sample-adapter
```

## Run

**Note:**
The adapter server requires an auth token on startup to initialize successfully. This is provided by the environment variable `AUTH_TOKENS_PATH`. The file must contain a JSON array of strings representing auth tokens. For example:
```json
["this-is-an-auth-token"]
```

To run the adapter locally for development and testing purposes, execute:

```bash
export AUTH_TOKENS_PATH={absolute-path-to-file-containing-tokens}
go run cmd/adapter/main.go 
# OR if you have a previously built binary, you can run
./sample-adapter
```
