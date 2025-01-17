package handler

import (
	"encoding/json"
	"fmt"
	expedition "goravel/app/models/expeditions"
	"io"
	"net/http"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Jnt Cargo handler
func JntCargoExpedition(c *gin.Context, resi string) {
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
		return
	}

	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan respons dari API J&T Cargo"})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data respons dari API J&T Cargo"})
		return
	}

	responseJntCargo(c, res, body)
}

func responseJntCargo(c *gin.Context, res *http.Response, body []byte) {
	var response expedition.Response
	var model expedition.JntCargoModel
	err := json.Unmarshal([]byte(body), &model)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memparsing data respons dari API J&T Cargo"})
		return
	}

	response.Resi = model.Data[0].Keyword
	response.Expedition = "J&T Cargo"
	timeFormat := "2006-01-02 15:04:05"
	for _, detail := range model.Data[0].Details {
		parsedTime, err := time.Parse(timeFormat, detail.ScanTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memparsing waktu dari API J&T Cargo"})
			return
		}
		response.Details = append(response.Details, expedition.Details{
			Time:    parsedTime,
			Message: detail.CustomerTracking,
		})
	}

	c.JSON(res.StatusCode, response)
}
