package handler

import (
	"fmt"
	expedition "goravel/app/models/expeditions"
	"goravel/app/utils"
	"io"
	"net/http"
)

func tokopediaKurirRekomendasi(resi string) []byte {
	url := "https://orchestra.tokopedia.com/orc/v1/microsite/tracking?airwaybill=" + resi
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("origin", "https://www.tokopedia.com")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Microsoft Edge\";v=\"128\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0")

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

func HandleTokopedia(resi string) expedition.Response {
	var response expedition.Response
	var model expedition.TokopediaKurirRekomendasi

	model, err := expedition.UnmarshalTokopediaKurirRekomendasi(tokopediaKurirRekomendasi(resi))

	if err != nil {
		fmt.Println(err)
		return expedition.Response{}
	}

	response.Expedition = "Tokopedia"
	response.Resi = model.Data[0].Airwaybill
	for _, trackingData := range model.Data[0].TrackingData {
		parseTime := utils.ParseTime(trackingData.TrackingTime)
		response.Details = append(response.Details, expedition.Details{
			Time:    parseTime,
			Message: trackingData.Message,
		})
	}

	return response
}
