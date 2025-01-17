package models

import "time"

type JntModel struct {
	Data struct {
		TrackingDirect []struct {
			ReferenceNo string `json:"referenceNo"`
			Logistic    struct {
				ID       string `json:"id"`
				Typename string `json:"__typename"`
			} `json:"logistic"`
			ShipmentDate string `json:"shipmentDate"`
			Details      []struct {
				Datetime       time.Time   `json:"datetime"`
				ShipperStatus  interface{} `json:"shipperStatus"`
				LogisticStatus struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Typename    string `json:"__typename"`
				} `json:"logisticStatus"`
				Typename string `json:"__typename"`
			} `json:"details"`
			Consigner struct {
				Name     string `json:"name"`
				Address  string `json:"address"`
				Typename string `json:"__typename"`
			} `json:"consigner"`
			Consignee struct {
				Name     string `json:"name"`
				Address  string `json:"address"`
				Typename string `json:"__typename"`
			} `json:"consignee"`
			Typename string `json:"__typename"`
		} `json:"trackingDirect"`
	} `json:"data"`
}
