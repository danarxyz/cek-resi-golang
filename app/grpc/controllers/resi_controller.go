package controllers

import (
	"context"
	"goravel/app/http/controllers"
	"goravel/app/http/controllers/handler"
	"goravel/app/models"
	"goravel/app/protos"

	"github.com/goravel/framework/facades"
)

type ResiController struct {
	protos.UnimplementedResiServiceServer
}

func NewResiController() *ResiController {
	return &ResiController{}
}

func (r *ResiController) CreateResi(ctx context.Context, req *protos.CreateResiRequest) (*protos.ResiResponse, error) {
	// get details from request
	resi_details := handler.HandleExpediton(req.Expedition, req.TrackingNum).GetNewStatus().Message
	resi := models.Resi{
		PackageName: req.PackageName,
		TrackingNum: req.TrackingNum,
		Expedition:  req.Expedition,
		Status:      "tracking",
		Details:     resi_details,
		Email:       req.Email,
	}

	if err := facades.Orm().Query().Create(&resi); err != nil {
		return nil, err
	}

	return &protos.ResiResponse{Resi: &protos.Resi{
		Id:          int32(resi.ID),
		PackageName: resi.PackageName,
		TrackingNum: resi.TrackingNum,
		Expedition:  resi.Expedition,
		Status:      resi.Status,
		Details:     resi.Details,
		Email:       resi.Email,
	}}, nil
}

func (r *ResiController) GetResi(ctx context.Context, req *protos.GetResiRequest) (*protos.ResiResponse, error) {
	var resi models.Resi
	if err := facades.Orm().Query().Where("tracking_num", req.TrackingNum).FirstOrFail(&resi); err != nil {
		return nil, err
	}

	return &protos.ResiResponse{Resi: &protos.Resi{
		Id:          int32(resi.ID),
		PackageName: resi.PackageName,
		TrackingNum: resi.TrackingNum,
		Expedition:  resi.Expedition,
		Status:      resi.Status,
		Details:     resi.Details,
		Email:       resi.Email,
	}}, nil
}

func (r *ResiController) UpdateResi(ctx context.Context, req *protos.UpdateResiRequest) (*protos.ResiResponse, error) {
	var resi models.Resi
	if err := facades.Orm().Query().Where("tracking_num", req.TrackingNum).FirstOrFail(&resi); err != nil {
		return nil, err
	}

	resi.Status = req.Status

	if err := facades.Orm().Query().Save(&resi); err != nil {
		return nil, err
	}

	return &protos.ResiResponse{Resi: &protos.Resi{
		Id:          int32(resi.ID),
		PackageName: resi.PackageName,
		TrackingNum: resi.TrackingNum,
		Expedition:  resi.Expedition,
		Status:      resi.Status,
		Details:     resi.Details,
		Email:       resi.Email,
	}}, nil
}

func (r *ResiController) DeleteResi(ctx context.Context, req *protos.DeleteResiRequest) (*protos.DeleteResiResponse, error) {
	var resi models.Resi
	if err := facades.Orm().Query().Where("tracking_num", req.TrackingNum).FirstOrFail(&resi); err != nil {
		return nil, err
	}

	if _, err := facades.Orm().Query().Delete(&resi); err != nil {
		return nil, err
	}

	return &protos.DeleteResiResponse{
		Success: true,
	}, nil
}
func (r *ResiController) GetAllResi(ctx context.Context, req *protos.Empty) (*protos.ResiListResponse, error) {
	var resis []models.Resi
	if err := facades.Orm().Query().Find(&resis); err != nil {
		return nil, err
	}

	var resiResponses []*protos.Resi
	for _, resi := range resis {
		resiResponses = append(resiResponses, &protos.Resi{
			Id:          int32(resi.ID),
			PackageName: resi.PackageName,
			TrackingNum: resi.TrackingNum,
			Expedition:  resi.Expedition,
			Status:      resi.Status,
			Details:     resi.Details,
			Email:       resi.Email,
		})
	}

	return &protos.ResiListResponse{Resis: resiResponses}, nil

}

func (r *ResiController) CheckAllResi(ctx context.Context, req *protos.Empty) (*protos.CheckAllResiResponse, error) {
	controllers.CheckExpedition()
	return &protos.CheckAllResiResponse{Success: true}, nil
}
