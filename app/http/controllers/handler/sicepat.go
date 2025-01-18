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

func sicepatExpedition(resi string) []byte {
	url := "https://paketmu.com/kurir/wp-admin/admin-ajax.php"
	method := "POST"

	payload := strings.NewReader(`action=lacak&resi=` + resi + `&kurir=sicepat`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("cookie", "_ga_Y23ZRF3RZD=GS1.1.1736141700.1.0.1736141700.0.0.0; _ga_8PHV5B1X02=GS1.1.1736141701.1.0.1736141701.0.0.0; _gid=GA1.2.1269731236.1736141703; _gat_gtag_UA_237466645_1=1; _ga_3752GP14P2=GS1.1.1736141702.1.0.1736141702.0.0.0; _ga=GA1.1.334732771.1736141701")
	req.Header.Add("origin", "https://paketmu.com")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://paketmu.com/kurir/sicepat/")
	req.Header.Add("sec-ch-ua", "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
	req.Header.Add("x-requested-with", "XMLHttpRequest")

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

func HandleSicepat(resi string) expedition.Response {
	var response expedition.Response
	var model expedition.Sicepat
	err := json.Unmarshal([]byte(sicepatExpedition(resi)), &model)

	if err != nil {
		return expedition.Response{}
	}

	response.Resi = model.Rajaongkir.Query.Waybill
	response.Expedition = "Sicepat"

	details := model.Rajaongkir.Result.Manifest
	for i := len(details) - 1; i >= 0; i-- {
		detail := details[i]
		response.Details = append(response.Details, expedition.Details{
			Time:    utils.ParseTime(detail.ManifestDate + " " + detail.ManifestTime),
			Message: detail.ManifestDescription,
		})
	}

	return response
}
