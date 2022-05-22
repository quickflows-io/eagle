## rpc

The RPC client interface definition layer strictly defines RPC interface routing, parameters and documents based on protobuf.

## Prepare

### Install the protoc compiler

```bash
$ PB_REL="https://github.com/protocolbuffers/protobuf/releases"
$ curl -LO $PB_REL/download/v3.12.1/protoc-3.12.1-linux-x86_64.zip

$ unzip protoc-3.12.1-linux-x86_64.zip -d /usr/local

$ export PATH="$PATH:/usr/local/bin"
```

View version

```bash
protoc --version
libprotoc 3.12.1
```

### Install the protoc-gen-go plugin

run：

```shell script
go get -u github.com/golang/protobuf/{helloworld,protoc-gen-go}
```

After compilation, `protoc-gen-go` will be installed to the `$GOBIN` directory, which defaults to `$GOPATH/bin`.   
This directory must be in the system's environment variable `$PATH` so that the `protocol` compiler can find plugins when `.proto` files are compiled.

### install grpc-go

grpc-go contains Go's grpc library

```
$ go get -u google.golang.org/grpc@v1.27.0
```

### Edit proto file

```shell script
protoc -I . --go_out=plugins=grpc,paths=source_relative:. user.proto
```

> paths options: import(default), source_relative

The server and client code will be generated in `user.proto.go`.

## Directory Structure

Usually one service one folder. There are versions under the service, one version and one folder. Internal services generally use v0 as the version.

A version can define multiple services, each with a proto file.

A typical directory structure is as follows:

```bash
rpc/user # business service
└── v0   # Service version
    ├── user.pb.go     # protobuf message Define Code [Auto Generated]
    └── user.proto     # protobuf Description File [Business Party Definition]
```

## define interface

Service interfaces are described using protobuf.

```proto
syntax = "proto3";

package user.v0; //Package name, consistent with the directory

// Service name, as long as a service can be defined
service Echo {
  // Service method, defined as needed
  rpc Hello(HelloRequest) returns (HelloResponse);
}

// input parameter definition
message HelloRequest {
  // Field definitions, only supported if the form transfer is used
  // int32, int64, uint32, unint64, double, float, bool, string,
  // And the corresponding repeated type, map and message types are not supported!
  // The framework automatically parses and converts parameter types
  // No limit if transferring with json or protobuf
  string message = 1; // This is an end-of-line comment, and business parties generally do not use it
  int32 age = 2;
  // form form format only partially supports repeated semantics
  // But the client needs to send a comma-separated string
  // If ids=1,2,3 will be resolved to []int32{1,2,3}
  repeated int32 ids = 3;
}

message HelloMessage {
  string message = 1;
}

// out parameter definition,
// Theoretically, any message can be output
// But our business requirements can only contain code, msg, data three fields,
// where data needs to be defined as message
// Open source versions can omit this convention
message HelloResponse {
  // business error code [machine readable], must be greater than zero
  / The master frame less than zero is in use, pay attention to avoid it.
  int32 code = 1;
  // business error message [human read]
  string msg = 2;
  // business data object
  HelloMessage data = 3;
}
```

## Generate code

```
# for specified services
protoc --go_out=. --twirp_out=. echo.proto

# for all services
find rpc -name '*.proto' -exec protoc --twirp_out=. --go_out=. {} \;
```

*.pb.go in the generated file is the definition code of the protobuf message, which supports both protobuf and json. *.twirp.go is the rpc routing related code.

## Implement interface

Please refer to [server/README.md](https://github.com/go-eagle/eagle/tree/master/internal/server/README.md)

## Automatic registration

The scaffolding provided by go-stone can automatically generate proto templates, server templates, and register routes. 
Run the following commands:

```bash
go run cmd/eagle/main.go rpc --server=foo --service=echo
```

will be automatically generated

```bash
rpc
└── foo
    └── v1
        ├── echo.pb.go
        ├── echo.proto
        └── echo.twirp.go
server
└── fooserver1
    └── echo.go
```

## Reference

- [Where is the Proto code?](https://eddycjy.com/posts/where-is-proto/)