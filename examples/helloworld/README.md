# rpc server and client

## define a proto

```proto
syntax = "proto3";

package helloworld;

option go_package="github.com/go-eagle/eagle/examples/helloworld/helloworld";

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  // Field definitions, only supported if the form transfer is used
  // int32, int64, uint32, unint64, double, float, bool, string,
  // And the corresponding repeated type, map and message types are not supported!
  // The framework automatically parses and converts parameter types
  // No limit if transferring with json or protobuf
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}

// Output parameter definition, in theory, any message can be output
// But we agreed to only include code, msg, data three fields,
// where data needs to be defined as message
// Open source versions can omit this convention
message HelloResponse {
  // Business Error Code [Machine Readable], must be greater than zero
  // The master frame less than zero is in use, pay attention to avoid it.
  int32 code = 1;
  // Business Error Messages [Human Read]
  string msg = 2;
  // business data object
  HelloReply data = 3;
}
```

## generated helper code

the protocol buffer compiler generates codes that has

- message serialization code(*.pb.go)
- remote interface stub for Client to call with the methods(*_grpc.pb.go)
- abstract interface for Server code to implement(*_grpc.pb.go)

## try it out.

enter the project root directory

```bash
cd {root_path}

go get -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

## gen proto

old-way
> use：github.com/golang/protobuf

```bash
$ protoc -I . --go_out=plugins=grpc,paths=source_relative:. examples/helloworld/protos/greeter.proto
```
> Generated `*.pb.go` contains message serialization code and `gRPC` code.

new-way

> 使用：google.golang.org/protobuf

```bash
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    examples/helloworld/protos/greeter.proto
```
> Two files `*.pb.go` and `*._grpc.pb.go` will be generated, which are message serialization code and `gRPC` code respectively.

> Official description：https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.20.0#v1.20-grpc-support

> https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code

What is the difference between `github.com/golang/protobuf` and `google.golang.org/protobuf`? https://developers.google.com/protocol-buffers/docs/reference/go/faq

## Run

1. run the server

```bash
cd examples/helloworld/server
go run main.go
```

2. run the client from another terminal

```bash
cd examples/helloworld/client
go run main.go
```

3. You’ll see the following output:

```bash
Greeting : "Hello eagle"
```

## Reference

- https://grpc.io/docs/languages/go/quickstart/
- https://developers.google.com/protocol-buffers/docs/proto3
- https://grpc.io/docs/guides/error/
- https://github.com/avinassh/grpc-errors/blob/master/go/client.go
- https://stackoverflow.com/questions/64828054/differences-between-protoc-gen-go-and-protoc-gen-go-grpc
- https://jbrandhorst.com/post/grpc-errors/
- https://godoc.org/google.golang.org/genproto/googleapis/rpc/errdetails
- https://cloud.google.com/apis/design/errors
- https://github.com/grpc/grpc/blob/master/doc/health-checking.md
- https://eddycjy.com/posts/where-is-proto/
- https://stackoverflow.com/questions/52969205/how-to-assert-grpc-error-codes-client-side-in-go
