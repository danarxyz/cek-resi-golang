package handler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	expedition "goravel/app/models/expeditions"
	"io"
	"net/http"

	"time"

	"github.com/goravel/framework/facades"
)

func spxTrackingNumber(resi string) string {
	config := facades.Config()
	k := config.Get("SPX_TOKEN").(string)
	r := float64(time.Now().UnixNano() / int64(time.Millisecond) / 1e3)
	h := sha256.New()
	rs := fmt.Sprintf("%d", int64(r))
	h.Write([]byte(resi + rs + k))
	return fmt.Sprintf(resi+"|"+rs+"%x", h.Sum(nil))
}

// Spx handler
func spxExpedition(resi string) []byte {
	trackingNum := spxTrackingNumber(resi)
	url := fmt.Sprintf("https://spx.co.id/api/v2/fleet_order/tracking/search?sls_tracking_number=%s", trackingNum)
	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Baca data dari respons API dan kirimkan sebagai respons ke klien (client)
	var responseData []byte
	responseData, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return responseData
}

func HandleSpx(resi string) expedition.Response {
	var response expedition.Response
	var model expedition.SpxModel

	err := json.Unmarshal([]byte(spxExpedition(resi)), &model)
	if err != nil {
		return response
	}

	response.Resi = model.Data.SlsTrackingNumber
	response.Expedition = "SPX"
	for _, detail := range model.Data.TrackingList {
		response.Details = append(response.Details, expedition.Details{
			Time:    time.Unix(int64(detail.Timestamp), 0),
			Message: detail.Message,
		})
	}

	return response
}
