package routes

import (
	"goravel/app/grpc/controllers"
	"goravel/app/protos"

	"github.com/goravel/framework/facades"
)

func Grpc() {
	protos.RegisterResiServiceServer(facades.Grpc().Server(), controllers.NewResiController())
}
