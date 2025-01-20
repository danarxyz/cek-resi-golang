package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func Server(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("Intercepting call to method: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error occurred: %v", err)
	}
	return resp, err
}
