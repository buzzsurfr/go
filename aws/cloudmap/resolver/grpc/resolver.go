// Package grpc implements a Cloud Map resolver. It sends the target
// name without scheme back to gRPC as resolved address.
//
// Based upon google.golang.org/grpc/resolver
package grpc

import (
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

var logger = grpclog.Component("cloudmap")

const scheme = "awscloudmap"

// NewBuilder builds a new resolver builder
func NewBuilder() resolver.Builder {
	return &cloudMapBuilder{}
}

type cloudMapBuilder struct{}

// Build builds the resolver.
//   target can be in the format "service.namespace[:port]" or "awscloudmap:///service.namespace[:port]"
func (*cloudMapBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &cloudMapResolver{
		target: target,
		cc:     cc,
		sess: session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})),
	}
	r.start()
	return r, nil
}

func (*cloudMapBuilder) Scheme() string {
	return scheme
}

type cloudMapResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	sess   *session.Session
}

func (r *cloudMapResolver) start() {
	svc := servicediscovery.New(r.sess)

	// Parse endpoint into service and namespace
	endpoint := parseEndpoint(r.target.Endpoint)

	// Discover instances
	result, err := svc.DiscoverInstances(&servicediscovery.DiscoverInstancesInput{
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
