// EXPERIMENTAL - DO NOT USE
// Package net provides net.Resolver instances implementing a wrapper using
// Cloud Map for service discovery.
//
// To replace the net.DefaultResolver with a Cloud Map resolver:
//
//     net.DefaultResolver = net.NewResolver()
//
package net

import (
	"context"
	"net"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

type DialerFunc func(ctx context.Context, network, address string) (net.Conn, error)

func NewResolver(parent *net.Resolver) *net.Resolver {
	if parent == nil {
		parent = &net.Resolver{}
	}

	return &net.Resolver{
		PreferGo:     true,
		StrictErrors: parent.StrictErrors,
		Dial:         NewDialer(),
	}
}

func NewDialer(parent DialerFunc) DialerFunc {
	return func(ctx context.Context, network, address string) (net.Conn, error) {

	}
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cmSvc := servicediscovery.New(sess)
	_ = cmSvc
}
