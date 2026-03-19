package authz

import (
	"context"

	authpb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}

func (u *Usecase) Allow(ctx context.Context, req *authpb.CheckRequest) (bool, error) {
	return true, nil
}
