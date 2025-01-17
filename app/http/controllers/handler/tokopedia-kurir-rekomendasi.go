package handler

import (
	"fmt"
	expedition "goravel/app/models/expeditions"
	"io"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
)

func TokopediaKurirRekomendasi(c *gin.Context, resi string) {
	url := "https://orchestra.tokopedia.com/orc/v1/microsite/tracking?airwaybill=" + resi
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get data from Tokopedia"})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read response from Tokopedia"})
		return
	}

	responseTokopedia(c, res, body)
}

func responseTokopedia(c *gin.Context, res *http.Response, body []byte) {
	var response expedition.Response
	var model expedition.TokopediaKurirRekomendasi

	model, err := expedition.UnmarshalTokopediaKurirRekomendasi(body)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal Tokopedia response"})
		return
	}

	response.Expedition = "Tokopedia"
	response.Resi = model.Data[0].Airwaybill
	timeFormat := "02 Jan 15:04 WIB"
	for _, trackingData := range model.Data[0].TrackingData {
		parseTime, err := time.Parse(timeFormat, trackingData.TrackingTime)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse time from Tokopedia"})
			return
		}
		response.Details = append(response.Details, expedition.Details{
			Time:    parseTime,
			Message: trackingData.Message,
		})
	}

	c.JSON(res.StatusCode, response)
}
