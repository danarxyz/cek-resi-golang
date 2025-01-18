package routes

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Api() {
	resiController := controllers.NewResiController()
	facades.Route().Get("/cekresi", resiController.Index)
	facades.Route().Post("/resi", resiController.AddExpedition)
	facades.Route().Put("/resi", resiController.UpdateStatus)
	facades.Route().Delete("/resi", resiController.DeleteExpedition)
	facades.Route().Get("/check", resiController.CheckExpedition)
	facades.Route().Get("/resi", resiController.Show)
}
