package authz

import (
	"context"

	authpb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type Checker interface {
	Allow(ctx context.Context, req *authpb.CheckRequest) (bool, error)
}
