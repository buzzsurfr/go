# go

Golang utilities and packages

## Packages

* `aws/cloudmap/resolver` - Cloud Map custom resolver for Go
  * `aws/cloudmap/resolver/grpc` gRPC resolver that implements [grpc/resolver.Resolver](https://pkg.go.dev/google.golang.org/grpc/resolver#Resolver)
  * `awsutil/detector` - Detects the compute type

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

```sh
go get -u github.com/buzzsurfr/go/aws/cloudmap/resolver/grpc
```

### Usage

#### per-gRPC connection

To add the resolver to a grpc `Dial` call, ensure the hostname specifies the scheme `awscloudmap:///`, matches the Cloud Map service and namespace, and use the [`WithResolvers()`](https://pkg.go.dev/google.golang.org/grpc@v1.33.1#WithResolvers) function:

```go
grpc.WithResolvers(cloudmap.NewBuilder())
```

For example, the gRPC connection accepts the [`WithResolvers()`](https://pkg.go.dev/google.golang.org/grpc@v1.33.1#WithResolvers) funcdtion as a parameter. The connection string looks like:

```go
conn, err := grpc.Dial("awscloudmap:///service.namespace[:port]", grpc.WithResolvers(cloudmap.NewBuilder()))
```

#### Default resolver

The Cloud Map resolver registers itself as an available _scheme_ named `awscloudmap` and can be set to the default by calling the [`SetDefaultScheme()`](https://pkg.go.dev/google.golang.org/grpc@v1.33.1/resolver#SetDefaultScheme) function in your `init` function. If you set the default scheme to `awscloudmap`, then you do not need to include the scheme in your hostname, but _ALL_ requests for that application will go to Cloud Map by default.

```go
func init() {
  grpc.SetDefaultScheme("awscloudmap")
}
```

### AWS X-Ray integration

The Cloud Map resolver supports AWS X-Ray for sending traces used by the resolver.

The Builder supports a functional option `WithContext()` where the context can be passed into the package from the main program. The context is necessary in order for the xray package to identify the TraceID.

X-Ray needs to be integrated into your application. To do so, start a segment at the beginning of your program.

```go
xrayCtx, seg := xray.BeginSegment(context.Background(), "Segment Name")
defer seg.Close(nil)
```

Since the resolver is loaded at init time, adding X-Ray integration requires creating and overriding the existing resolver. Create the custom resolver with the context, then call `grpc.Dial` using `WithResolvers()` option.

```go
customResolver := cloudmap.NewBuilder(cloudmap.WithContext(xrayCtx))
conn, err := grpc.Dial("awscloudmap:///service.namespace[:port]", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithResolvers(customResolver))
```

## awsutil/detector

Package detector detects the compute type, supporting the following compute types (in this order):

* AWS Lambda
* Amazon ECS
* ~~Amazon EKS~~
* Kubernetes
* Docker
* Amazon EC2

The `Detect` function returns a `DetectOutput` type which does support casting the type to a string as well as boolean check functions (either as functions or methods):

* `IsLambda`
* `IsECS`
* `IsEKS`
* `IsKubernetes`
* `IsDocker`
* `IsEC2`

### Example

```go
import "detector"

func main() {
  fmt.Println(detector.Detect())
}
```
