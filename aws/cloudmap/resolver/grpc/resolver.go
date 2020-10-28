// Package grpc implements a Cloud Map resolver. It sends the targets without scheme back to gRPC as resolved addresses.
//
// Based upon google.golang.org/grpc/resolver
package grpc

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

var logger = grpclog.Component("awscloudmap")

const scheme = "awscloudmap"

// BuilderOption for passing options to your builder
type BuilderOption func(*cloudMapBuilder) error

// WithContext adds a context to your builder
func WithContext(ctx context.Context) BuilderOption {
	return func(b *cloudMapBuilder) error {
		b.context = ctx
		return nil
	}
}

// NewBuilder builds a new Cloud Map resolver builder. NewBuilder can be used inline with a grpc.Dial call.
//
// Example:
//     conn, err := grpc.Dial("service.namespace:50051", grpc.WithResolvers(cloudmap.NewBuilder())
func NewBuilder(opts ...BuilderOption) resolver.Builder {
	b := &cloudMapBuilder{
		context: context.Background(),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

type cloudMapBuilder struct {
	context context.Context
}

// Build builds the resolver.
//   target can be in the format "service.namespace[:port]" or "awscloudmap:///service.namespace[:port]"
func (b *cloudMapBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &cloudMapResolver{
		target:  target,
		cc:      cc,
		context: b.context,
	}
	r.start()
	return r, nil
}

func (*cloudMapBuilder) Scheme() string {
	return scheme
}

type cloudMapResolver struct {
	target  resolver.Target
	cc      resolver.ClientConn
	context context.Context
}

func (r *cloudMapResolver) start() {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		logger.Errorf("unable to load SDK config, %v", err)
	}
	svc := servicediscovery.NewFromConfig(cfg)

	// Parse endpoint into service and namespace
	endpoint := parseEndpoint(r.target.Endpoint)

	// Discover instances
	result, err := svc.DiscoverInstances(r.context, &servicediscovery.DiscoverInstancesInput{
		HealthStatus:  types.HealthStatusFilterAll,
		MaxResults:    aws.Int32(10),
		NamespaceName: aws.String(endpoint.namespace),
		ServiceName:   aws.String(endpoint.service),
	})
	var serviceNotFoundErr *types.ServiceNotFound
	if errors.As(err, &serviceNotFoundErr) {
		logger.Errorln(serviceNotFoundErr.ErrorMessage())
	}
	var namespaceNotFoundErr *types.NamespaceNotFound
	if errors.As(err, &namespaceNotFoundErr) {
		logger.Errorln(namespaceNotFoundErr.ErrorMessage())
	}
	var invalidInputErr *types.InvalidInput
	if errors.As(err, &invalidInputErr) {
		logger.Errorln(invalidInputErr.ErrorMessage())
	}
	var requestLimitExceededErr *types.RequestLimitExceeded
	if errors.As(err, &requestLimitExceededErr) {
		logger.Errorln(requestLimitExceededErr.ErrorMessage())
	}

	// Format to resolver.State
	addrs := []resolver.Address{}
	for _, instance := range result.Instances {
		addrs = append(addrs, resolver.Address{
			Addr: aws.ToString(instance.Attributes["AWS_INSTANCE_IPv4"]),
		})
	}

	r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint}}})
}

func (*cloudMapResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*cloudMapResolver) Close() {}

type cloudMapEndpoint struct {
	service   string
	namespace string
}

func parseEndpoint(endpoint string) cloudMapEndpoint {
	host, _, _ := net.SplitHostPort(endpoint)
	split := strings.SplitN(host, ".", 2)
	return cloudMapEndpoint{
		service:   split[0],
		namespace: split[1],
	}
}

func init() {
	resolver.Register(NewBuilder())
}
