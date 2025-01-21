package grpc

import (
	"fmt"
	"goravel/app/grpc/interceptors"
	"goravel/app/utils"
	"goravel/routes"
	"log"

	"google.golang.org/grpc"
)

type Kernel struct {
}

// The application's global GRPC interceptor stack.
// These middleware are run during every request to your application.
func (kernel Kernel) UnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		interceptors.Server,
	}
}

// The application's client interceptor groups.
func (kernel Kernel) UnaryClientInterceptorGroups() map[string][]grpc.UnaryClientInterceptor {
	return map[string][]grpc.UnaryClientInterceptor{}
}

func CustomGrpcServer(host string) error {
	creds, err := utils.LoadTLSCredentials()
	if err != nil {
		return fmt.Errorf("cannot load TLS credentials: %w", err)
	}

	server := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(Kernel{}.UnaryServerInterceptors()...),
	)

	// Register services
	routes.Grpc(server)

	// Listen and serve
	log.Printf("\033[32mStarting GRPC server on %s\033[0m\n", host)
	listener, err := utils.CreateListener(host)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
