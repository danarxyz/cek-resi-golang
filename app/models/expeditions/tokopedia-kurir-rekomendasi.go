// To parse and unparse this JSON data, add this code to your project and do:
//
//    tokopediaKurirRekomendasi, err := UnmarshalTokopediaKurirRekomendasi(bytes)
//    bytes, err = tokopediaKurirRekomendasi.Marshal()

package models

import "encoding/json"

func UnmarshalTokopediaKurirRekomendasi(data []byte) (TokopediaKurirRekomendasi, error) {
	var r TokopediaKurirRekomendasi
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TokopediaKurirRekomendasi) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TokopediaKurirRekomendasi struct {
	Status  int64   `json:"status"`
	Message string  `json:"message"`
	Data    []Datum `json:"data"`
}

type Datum struct {
	InputAwb     string            `json:"input_awb"`
	Courier      string            `json:"courier"`
	Service      string            `json:"service"`
	Airwaybill   string            `json:"airwaybill"`
	Status       string            `json:"status"`
	ErrorMessage string            `json:"error_message"`
	Seller       Buyer             `json:"seller"`
	Buyer        Buyer             `json:"buyer"`
	ShippingData [][]ShippingDatum `json:"shipping_data"`
	TrackingData []TrackingDatum   `json:"tracking_data"`
}

type Buyer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ShippingDatum struct {
	Title   string  `json:"title"`
	Content string  `json:"content"`
	RawTime *string `json:"raw_time,omitempty"`
}

type TrackingDatum struct {
	TrackingTime string `json:"tracking_time"`
	Message      string `json:"message"`
	CityName     string `json:"city_name"`
	PartnerName  string `json:"partner_name"`
}
