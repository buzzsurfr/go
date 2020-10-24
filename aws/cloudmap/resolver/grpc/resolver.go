// Package grpc implements a Cloud Map resolver. It sends the targets without scheme back to gRPC as resolved addresses.
//
// Based upon google.golang.org/grpc/resolver
package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/aws/aws-xray-sdk-go/xray"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

var logger = grpclog.Component("cloudmap")

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
		target: target,
		cc:     cc,
		sess: session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})),
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
	sess    *session.Session
	context context.Context
}

func (r *cloudMapResolver) start() {
	svc := servicediscovery.New(r.sess)
	if xray.SdkDisabled() {
		xray.AWS(svc.Client)
	}

	// Parse endpoint into service and namespace
	endpoint := parseEndpoint(r.target.Endpoint)

	// Discover instances
	result, err := svc.DiscoverInstancesWithContext(r.context, &servicediscovery.DiscoverInstancesInput{
		HealthStatus:  aws.String("ALL"),
		MaxResults:    aws.Int64(10),
		NamespaceName: aws.String(endpoint.namespace),
		ServiceName:   aws.String(endpoint.service),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case servicediscovery.ErrCodeServiceNotFound:
				logger.Errorln(servicediscovery.ErrCodeServiceNotFound, aerr.Error())
			case servicediscovery.ErrCodeNamespaceNotFound:
				logger.Errorln(servicediscovery.ErrCodeNamespaceNotFound, aerr.Error())
			case servicediscovery.ErrCodeInvalidInput:
				logger.Errorln(servicediscovery.ErrCodeInvalidInput, aerr.Error())
			case servicediscovery.ErrCodeRequestLimitExceeded:
				logger.Errorln(servicediscovery.ErrCodeRequestLimitExceeded, aerr.Error())
			default:
				logger.Errorln(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logger.Errorln(err.Error())
		}
	}

	// Format to resolver.State
	addrs := []resolver.Address{}
	for _, instance := range result.Instances {
		addrs = append(addrs, resolver.Address{
			Addr: aws.StringValue(instance.Attributes["AWS_INSTANCE_IPv4"]),
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
