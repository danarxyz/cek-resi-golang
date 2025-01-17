package handler

import (
	"encoding/json"
	"fmt"
	expedition "goravel/app/models/expeditions"
	"io"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

func JntExpedition(c *gin.Context, resi string) {

	url := "https://gql-web.shipper.id/query"
	method := "POST"

	payload := strings.NewReader("{\"query\":\"query trackingDirect($input: TrackingDirectInput!) {\\n  trackingDirect(p: $input) {\\n    referenceNo\\n    logistic {\\n      id\\n      __typename\\n    }\\n    shipmentDate\\n    details {\\n      datetime\\n      shipperStatus {\\n        name\\n        description\\n        __typename\\n      }\\n      logisticStatus {\\n        name\\n        description\\n        __typename\\n      }\\n      __typename\\n    }\\n    consigner {\\n      name\\n      address\\n      __typename\\n    }\\n    consignee {\\n      name\\n      address\\n      __typename\\n    }\\n    __typename\\n  }\\n}\",\"variables\":{\"input\":{\"logisticId\":9,\"referenceNo\":[\"" + resi + "\"]}}}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://shipper.id")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://shipper.id/")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Microsoft Edge\";v=\"128\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0")
	req.Header.Add("x-app-name", "shp-homepage-v5")
	req.Header.Add("x-app-version", "1.0.0")

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan respons dari API J&T"})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data respons dari API J&T"})
		return
	}

	responseJnt(c, res, body)
}

func responseJnt(c *gin.Context, res *http.Response, body []byte) {
	var response expedition.Response
	var model expedition.JntModel

	if err := json.Unmarshal(body, &model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data respons dari API J&T"})
		return
	}

	response.Resi = model.Data.TrackingDirect[0].ReferenceNo
	response.Expedition = "J&T"

	details := model.Data.TrackingDirect[0].Details
	for i := len(details) - 1; i >= 0; i-- {
		detail := details[i]
		response.Details = append(response.Details, expedition.Details{
			Time:    detail.Datetime,
			Message: detail.LogisticStatus.Description,
		})
	}

	c.JSON(res.StatusCode, response)
}
