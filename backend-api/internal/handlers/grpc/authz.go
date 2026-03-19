package grpc

import (
	"context"

	authpb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	envoyfilter "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/envoy_filter"
)

type AuthzHandler struct {
	authpb.UnimplementedAuthorizationServer
	usecase envoyfilter.Usecase
}

func NewAuthzHandler(usecase envoyfilter.Usecase) *AuthzHandler {
	return &AuthzHandler{usecase: usecase}
}

func (h *AuthzHandler) Check(ctx context.Context, req *authpb.CheckRequest) (*authpb.CheckResponse, error) {
	decision, err := h.usecase.Check(ctx, req)
	if err != nil {
		return nil, err
	}

	if decision == nil || !decision.Allowed {
		return &authpb.CheckResponse{
			Status: &statuspb.Status{Code: int32(codes.PermissionDenied)},
			HttpResponse: &authpb.CheckResponse_DeniedResponse{
				DeniedResponse: &authpb.DeniedHttpResponse{},
			},
		}, nil
	}

	return &authpb.CheckResponse{
		Status: &statuspb.Status{Code: int32(codes.OK)},
		HttpResponse: &authpb.CheckResponse_OkResponse{
			OkResponse: &authpb.OkHttpResponse{},
		},
	}, nil
}
