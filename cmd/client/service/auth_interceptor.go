package client_service

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	// err := interceptor.scheduleRefreshToken(refreshDuration)

	// if err != nil {
	// 	return nil, err
	// }

	return interceptor, nil
}

// Unary returns a client interceptor to authenticate unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> unary interceptor: %s", method)

		log.Println(interceptor.authMethods["access_token"])
		interceptor.accessToken = "hello"
		ctx.Value("hello")

		// if interceptor.authMethods[method] {
		// 	return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		// }

		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		// return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Stream returns a client interceptor to authenticate stream RPC
func (interceptor *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		interceptor.accessToken = "hello"
		ctx.Value("hello")

		// if interceptor.authMethods[method] {
		// 	return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
		// }

		return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
		// return streamer(ctx, desc, cc, method, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "access_token", interceptor.accessToken)
}
