package models

import "time"

type Jne struct {
	Data Data `json:"data"`
}

type Data struct {
	TrackingDirect []TrackingDirect `json:"trackingDirect"`
}

type TrackingDirect struct {
	ReferenceNo  string   `json:"referenceNo"`
	Logistic     Logistic `json:"logistic"`
	ShipmentDate string   `json:"shipmentDate"`
	Details      []Detail `json:"details"`
	Consigner    Consigne `json:"consigner"`
	Consignee    Consigne `json:"consignee"`
	Typename     string   `json:"__typename"`
}

type Consigne struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Typename string `json:"__typename"`
}

type Detail struct {
	Datetime       time.Time      `json:"datetime"`
	ShipperStatus  interface{}    `json:"shipperStatus"`
	LogisticStatus LogisticStatus `json:"logisticStatus"`
	Typename       DetailTypename `json:"__typename"`
}

type LogisticStatus struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Typename    LogisticStatusTypename `json:"__typename"`
}

type Logistic struct {
	ID       string `json:"id"`
	Typename string `json:"__typename"`
}

type LogisticStatusTypename string

const (
	TrackingsvcTrackingDirectStatus LogisticStatusTypename = "TrackingsvcTrackingDirectStatus"
)

type DetailTypename string

const (
	TrackingsvcTrackingDirectDetail DetailTypename = "TrackingsvcTrackingDirectDetail"
)
