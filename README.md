# Go RSS Reader

This Go package can be used to parse multiple RSS feed URLs asynchronously.

## Development

### 1. Go RSS Reader Package (`rssreader`)

- Create `rssreader.go`

- File Structure (Go):
```
rssreader/
├── rssreader.go
└── rssreader_test.go // Unit tests
```

- Console Commands (Go Package):
```
go mod init rssreader // Initialize module | github.com/yourusername/rssreader
go test ./... // Run tests
```

### 2. Go RSS Reader Service

- Create `rssservice` folder and `main.go`

- File Structure (Go Service):
```
rssservice/
├── main.go
├── go.mod
└── go.sum
```

- Console Commands (Go Service):
```
go mod init rssservice // Initialize module | github.com/yourusername/rssservice
go mod tidy // Download dependencies
go run main.go
```

## Usage

1. Initialize the Module

Ensure you are in the `rssreader/` directory and run:
```
go mod init rssreader
go mod tidy
```

2. Install gRPC and Protobuf Compiler
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
Ensure protoc is installed on your system.

3. Generate gRPC Code from `rssreader.proto`
```
protoc --go_out=. --go-grpc_out=. rssreader.proto
```

4. Run the gRPC Server (`main.go`)
```
go run main.go
```

5. Test the Server with a gRPC Client

You can create a basic client or use tools like grpcurl:
```
grpcurl -plaintext localhost:50051 list
```
Let me know if you'd like me to generate the gRPC client code for testing!

or use the created `client.go` client for tests.


6. Add (and Run) the additional automated Tests

Execute the test file using the go test command:
```
go test ./...
```

7. Run Tests with Verbose Output

For a more detailed output, use the -v flag:
```
go test -v ./...
```

8. Expected Output

If the tests pass, you should see something like:
```
=== RUN   TestParse
--- PASS: TestParse (0.00s)
PASS
```
If they fail, you'll get an error message with details about the failure.
