# go
Golang utilities and packages

## Packages

* `aws/cloudmap/resolver` - Cloud Map custom resolver for Go
  * `aws/cloudmap/resolver/grpc` gRPC resolver that implements [grpc/resolver.Resolver](https://pkg.go.dev/google.golang.org/grpc/resolver#Resolver)


## aws/cloudmap/resolver/grpc

Package grpc implements a Cloud Map resolver. It sends the target name without scheme back to gRPC as resolved address.

Based upon google.golang.org/grpc/resolver

### Installation

With [Go module](https://github.com/golang/go/wiki/Modules) support (Go 1.11+), simply add the following import

```go
import (
    cloudmap "github.com/buzzsurfr/go/aws/cloudmap/resolver/grpc"
)
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, to install the package, run the following command:

```
go get -u github.com/buzzsurfr/go/aws/cloudmap/resolver/grpc
```

### Usage

#### per-gRPC connection

To add the resolver to a grpc `Dial` call, ensure the hostname matches the Cloud Map service and namespace, and use the [`WithResolvers()`](https://pkg.go.dev/google.golang.org/grpc@v1.33.1#WithResolvers) function:
```go
grpc.WithResolvers(cloudmap.NewBuilder())
```

For example, the gRPC connection accepts the [`WithResolvers()`](https://pkg.go.dev/google.golang.org/grpc@v1.33.1#WithResolvers) funcdtion as a parameter. The connection string looks like:

```go
conn, err := grpc.Dial("service.namespace[:port]", grpc.WithResolvers(cloudmap.NewBuilder()))
```

#### Default resolver

The Cloud Map resolver registers itself as an available _scheme_ named `cloudmap` and can be set to the default by calling the [`SetDefaultScheme()`](https://pkg.go.dev/google.golang.org/grpc@v1.33.1/resolver#SetDefaultScheme) function in your `init` function.

```go
func init() {
	grpc.SetDefaultScheme("awscloudmap")
}
```
