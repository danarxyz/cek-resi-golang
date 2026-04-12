package routes

import (
	"goravel/app/grpc/controllers"
	"goravel/app/protos"

	"google.golang.org/grpc"
)

func Grpc(server *grpc.Server) {
	protos.RegisterResiServiceServer(server, controllers.NewResiController())
}
