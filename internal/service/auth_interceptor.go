package service

import (
	"context"
	// "fmt"
	"log"

	"github.com/MyyPo/grpc-chat/internal/util"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/metadata"
	// "google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	tokenManager *util.TokenManager
}

func NewAuthInterceptor(tokenManager *util.TokenManager) *AuthInterceptor {
	return &AuthInterceptor{tokenManager: tokenManager}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> Unary interceptor: ", info.FullMethod)

		// if err := interceptor.authorize(ctx, info.FullMethod); err != nil {
		// 	return nil, err
		// }

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> Stream interceptor: ", info.FullMethod)

		// if err := interceptor.authorize(stream.Context(), info.FullMethod); err != nil {
		// 	return err
		// }
		return handler(srv, stream)
	}
}

// func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	fmt.Println(md)
// 	if !ok {
// 		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
// 	}

// 	values := md["access_token"]
// 	if len(values) == 0 {
// 		return status.Errorf(codes.Unauthenticated, "access token not provided")
// 	}
// 	accessToken := values[0]

// 	err := interceptor.tokenManager.ValidateToken(accessToken, true)
// 	if err != nil {
// 		return status.Errorf(codes.Unauthenticated, "invalid access token: %v", err)
// 	}

// 	return nil
// }