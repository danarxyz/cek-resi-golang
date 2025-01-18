package handler

import (
	"encoding/json"
	"fmt"
	expedition "goravel/app/models/expeditions"
	"io"
	"net/http"

	"strings"
)

func jneExpedition(resi string) []byte {
	url := "https://gql-web.shipper.id/query"
	method := "POST"

	payload := strings.NewReader("{\"query\":\"query trackingDirect($input: TrackingDirectInput!) {\\n  trackingDirect(p: $input) {\\n    referenceNo\\n    logistic {\\n      id\\n      __typename\\n    }\\n    shipmentDate\\n    details {\\n      datetime\\n      shipperStatus {\\n        name\\n        description\\n        __typename\\n      }\\n      logisticStatus {\\n        name\\n        description\\n        __typename\\n      }\\n      __typename\\n    }\\n    consigner {\\n      name\\n      address\\n      __typename\\n    }\\n    consignee {\\n      name\\n      address\\n      __typename\\n    }\\n    __typename\\n  }\\n}\",\"variables\":{\"input\":{\"logisticId\":1,\"referenceNo\":[\"" + resi + "\"]}}}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
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
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return body
}

func HandleJNE(resi string) expedition.Response {
	var jne expedition.Jne
	var response expedition.Response

	err := json.Unmarshal(jneExpedition(resi), &jne)
	if err != nil {
		return expedition.Response{}
	}

	if len(jne.Data.TrackingDirect) == 0 {
		return expedition.Response{}
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

	return response
}
