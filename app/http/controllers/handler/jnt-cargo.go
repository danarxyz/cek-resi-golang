package handler

import (
	"encoding/json"
	"fmt"
	expedition "goravel/app/models/expeditions"
	"goravel/app/utils"
	"io"
	"net/http"

	"strings"
)

// Jnt Cargo handler
func jntCargoExpedition(resi string) []byte {
	url := "https://office.jtcargo.co.id/official/waybill/trackingCustomerByWaybillNo"
	method := "POST"

	payload := strings.NewReader(`{
    "waybillNo": "` + resi + `",
    "langType": "ID",
    "searchWaybillOrCustomerOrderId": "1"
}`)

	client := http.DefaultClient
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	return body
}

func HandleJNTCargo(resi string) expedition.Response {
	var response expedition.Response
	var model expedition.JntCargoModel
	err := json.Unmarshal([]byte(jntCargoExpedition(resi)), &model)

	if err != nil {
		return expedition.Response{}
	}

	response.Resi = model.Data[0].Keyword
	response.Expedition = "J&T Cargo"
	for _, detail := range model.Data[0].Details {
		parsedTime := utils.ParseTime(detail.ScanTime)
		response.Details = append(response.Details, expedition.Details{
			Time:    parsedTime,
			Message: detail.CustomerTracking,
		})
	}

	return response
}
