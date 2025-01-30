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


# Run main.go
go run cmd/adapter/main.go 

# OR if you have a previously built binary, you can run
./sample-adapter

# OR if you choose to run the docker image, you'll need to provide the file and env variable.
# In a typical deployment, it is done by mounting a volume containing the file. For example:
docker run -d --name sample-adapter \
    -v /local/path/to/file:/container/secrets \
    -e ADAPTER_TOKENS_PATH=/container/secrets/auth.json \
    sample-adapter:latest
```

### Fetch Data from the System of Record

By default, the adapter listens on port 8080. You can use Postman to send a gRPC request to the adapter by following these steps:

1. Define the [`GetPage` Protobuf definition](https://github.com/SGNL-ai/adapter-framework/blob/f2cafb0d963b54c350350967906ce59776d720a1/api/adapter/v1/adapter.proto).

2. In the sidebar, click on **Collections** and create a new collection with the type set to **gRPC**.

3. Within this new collection, create a new gRPC request. Enter the URL of the adapter (e.g., `http://localhost:8080`) and select the `GetPage` method from the dropdown.

4. In the **Metadata** tab, add a `token` key and set its value to one of the tokens in the `AUTH_TOKENS_PATH` file.

5. In the **Message** tab, enter the `GetPage` request following the schema defined in step 1.

An example gRPC request:

```json
{
    "cursor": "",
    "datasource": {
        "type": "SCIM2.0-1.0.0", // The type here should match the adapter type defined in `cmd/adapter/main.go`.
        "address": "{{address}}",
        "auth": {
            "http_authorization": "Bearer {{token}}"
        },
        "config": "{{b64_encoded_string}}"
    },
    "entity": {
        "attributes": [
            {
                "external_id": "id",
                "type": "ATTRIBUTE_TYPE_STRING",
                "id": "id"
            }
        ],
        "external_id": "Users",
        "id": "Users",
        "ordered": false
    },
    "page_size": "100"
}
```

The `config` should be a base64 encoded string of the `Config` struct defined in `config.go`. For example, if the `Config` struct is:

```go
type Config struct {
    APIVersion string `json:"apiVersion,omitempty"`
}
```

then the `config` field should be:

```json
{
    "apiVersion": "v1"
}
```

which is base64 encoded to `eyJhcGlWZXJzaW9uIjoidjEifQ==`.

Hit Send!