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

func JNEExpedition(c *gin.Context, resi string) {
	url := "https://gql-web.shipper.id/query"
	method := "POST"

	payload := strings.NewReader("{\"query\":\"query trackingDirect($input: TrackingDirectInput!) {\\n  trackingDirect(p: $input) {\\n    referenceNo\\n    logistic {\\n      id\\n      __typename\\n    }\\n    shipmentDate\\n    details {\\n      datetime\\n      shipperStatus {\\n        name\\n        description\\n        __typename\\n      }\\n      logisticStatus {\\n        name\\n        description\\n        __typename\\n      }\\n      __typename\\n    }\\n    consigner {\\n      name\\n      address\\n      __typename\\n    }\\n    consignee {\\n      name\\n      address\\n      __typename\\n    }\\n    __typename\\n  }\\n}\",\"variables\":{\"input\":{\"logisticId\":1,\"referenceNo\":[\"" + resi + "\"]}}}")

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
	req.Header.Add("sec-ch-ua", "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
	req.Header.Add("x-app-name", "shp-homepage-v5")
	req.Header.Add("x-app-version", "1.0.0")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseJNE(c, res, body)
}

func responseJNE(c *gin.Context, res *http.Response, body []byte) {
	var jne expedition.Jne
	var response expedition.Response

	err := json.Unmarshal(body, &jne)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal unmarshal data"})
		return
	}

	if len(jne.Data.TrackingDirect) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	data := jne.Data.TrackingDirect[0]
	response.Resi = data.ReferenceNo
	response.Expedition = "JNE"

	for i := len(data.Details) - 1; i >= 0; i-- {
		detail := data.Details[i]
		response.Details = append(response.Details, expedition.Details{
			Time:    detail.Datetime,
			Message: detail.LogisticStatus.Description,
		})
	}

	c.JSON(res.StatusCode, response)
}
