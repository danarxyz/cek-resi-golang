package controllers

import (
	"goravel/app/http/controllers/handler"
	"goravel/app/mails"
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ResiController struct {
	// Dependent services
}

func NewResiController() *ResiController {
	return &ResiController{
		// Inject services
	}
}

func (r *ResiController) Index(ctx http.Context) http.Response {
	resi := ctx.Request().Query("sls_tracking_number")
	type_expedition := ctx.Request().Query("type")

	switch type_expedition {
	case "spx":
		return ctx.Response().Success().Json(handler.HandleSpx(resi))
	case "jnt-cargo":
		return ctx.Response().Success().Json(handler.HandleJNTCargo(resi))
	case "jnt":
		return ctx.Response().Success().Json(handler.HandleJNT(resi))
	case "tokopedia":
		return ctx.Response().Success().Json(handler.HandleTokopedia(resi))
	case "sicepat":
		return ctx.Response().Success().Json(handler.HandleSicepat(resi))
	case "jne":
		return ctx.Response().Success().Json(handler.HandleJNE(resi))
	default:
		ctx.Response().Status(400).Json(http.Json{
			"message": "Ekspedisi tidak ditemukan",
		})
	}
	return ctx.Response().Success().Json(http.Json{
		"message": "Ekspedisi tidak ditemukan",
	})
}

func (r *ResiController) AddExpedition(ctx http.Context) http.Response {
	resi := ctx.Request().Input("sls_tracking_number")
	type_expedition := ctx.Request().Input("type")
	email := ctx.Request().Input("email")
	validation, err := facades.Validation().Make(ctx.Context(), map[string]any{
		"sls_tracking_number": resi,
		"type":                type_expedition,
		"email":               email,
	}, map[string]string{
		"sls_tracking_number": "required|string",
		"type":                "required|string",
		"email":               "required|email",
	})

	// check if an error occured, might not be validation error
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Validation setup failed",
			"error":   err.Error(),
		})
	}

	// check for validation errors
	if validation.Fails() {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"success": false,
			"message": "Validation failed",
			"errors":  validation.Errors().All(),
		})
	}
	// get details
	details := handler.HandleExpediton(type_expedition, resi).GetNewStatus().Message
	newData := &models.Resi{
		TrackingNum: resi,
		Expedition:  type_expedition,
		Status:      "tracking",
		Details:     details,
		Email:       email,
	}
	var data models.Resi
	// check if the data already exists
	if err := facades.Orm().Query().Where("tracking_num", resi).FirstOrFail(&data); err == nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"success": false,
			"message": "Expedition already exists",
		})
	}

	if err := facades.Orm().Query().Create(newData); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to add new expedition",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"success": true,
		"message": "Expedition added successfully",
		"data":    newData,
	})
}

func (r *ResiController) UpdateStatus(ctx http.Context) http.Response {
	resi := ctx.Request().Query("sls_tracking_number")
	status := ctx.Request().Query("status")
	// check if the data exists
	var data models.Resi
	if err := facades.Orm().Query().Where("tracking_num", resi).FirstOrFail(&data); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"success": false,
			"message": "Expedition not found",
		})
	}
	data.Status = status
	if err := facades.Orm().Query().Save(&data); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to update expedition",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"success": true,
		"message": "Expedition updated successfully",
		"data":    data,
	})
}

func (r *ResiController) DeleteExpedition(ctx http.Context) http.Response {
	resi := ctx.Request().Query("sls_tracking_number")
	// check if the data exists
	var data models.Resi
	if err := facades.Orm().Query().Where("tracking_num", resi).FirstOrFail(&data); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"success": false,
			"message": "Expedition not found",
		})
	}
	if _, err := facades.Orm().Query().Delete(&data); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to delete expedition",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"success": true,
		"message": "Expedition deleted successfully",
	})
}

func (r *ResiController) CheckExpedition(ctx http.Context) http.Response {
	CheckExpedition()
	return ctx.Response().Success().Json(http.Json{
		"success": true,
		"message": "Expedition checked successfully",
	})
}

func CheckExpedition() {
	// get all data with status tracking
	var data []models.Resi
	if err := facades.Orm().Query().Where("status", "tracking").Get(&data); err != nil {
		// write log
		return
	}
	// check each data, if detail status is different, update the status and send email
	for _, d := range data {
		details := handler.HandleExpediton(d.Expedition, d.TrackingNum)
		if d.Details != details.GetNewStatus().Message {
			d.Details = details.GetNewStatus().Message
			if err := facades.Orm().Query().Save(&d); err != nil {
				// write log
				continue
			}
			if err := sendMail(d); err != nil {
				// write log
				continue
			}
		}
	}
}

func sendMail(resi models.Resi) error {
	mail := mails.NewNotify(resi)
	err := facades.Mail().Queue(mail)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResiController) Show(ctx http.Context) http.Response {
	// get all data
	var data []models.Resi
	if err := facades.Orm().Query().OrderBy("created_at", "asc").Get(&data); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to get expeditions",
			"error":   err.Error(),
		})
	}
	return ctx.Response().Success().Json(http.Json{
		"success": true,
		"message": "Expeditions retrieved successfully",
		"data":    data,
	})
}
