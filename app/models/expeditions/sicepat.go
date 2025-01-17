package models

type Sicepat struct {
	Rajaongkir Rajaongkir `json:"rajaongkir"`
	Bg         string     `json:"bg"`
}

type Rajaongkir struct {
	Query  Query  `json:"query"`
	Status Status `json:"status"`
	Result Result `json:"result"`
}

type Query struct {
	Waybill string `json:"waybill"`
	Courier string `json:"courier"`
}

type Result struct {
	Delivered      bool           `json:"delivered"`
	Summary        Summary        `json:"summary"`
	Details        details        `json:"details"`
	Manifest       []Manifest     `json:"manifest"`
	DeliveryStatus DeliveryStatus `json:"delivery_status"`
}

type DeliveryStatus struct {
	Status      string `json:"status"`
	PodReceiver string `json:"pod_receiver"`
	PodDate     string `json:"pod_date"`
	PodTime     string `json:"pod_time"`
}

type details struct {
	WaybillNumber    string `json:"waybill_number"`
	WaybillDate      string `json:"waybill_date"`
	WaybillTime      string `json:"waybill_time"`
	Weight           string `json:"weight"`
	Origin           string `json:"origin"`
	Destination      string `json:"destination"`
	ShippperName     string `json:"shippper_name"`
	ShipperAddress1  string `json:"shipper_address1"`
	ShipperAddress2  string `json:"shipper_address2"`
	ShipperAddress3  string `json:"shipper_address3"`
	ShipperCity      string `json:"shipper_city"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverAddress1 string `json:"receiver_address1"`
	ReceiverAddress2 string `json:"receiver_address2"`
	ReceiverAddress3 string `json:"receiver_address3"`
	ReceiverCity     string `json:"receiver_city"`
}

type Manifest struct {
	ManifestCode        string `json:"manifest_code"`
	ManifestDescription string `json:"manifest_description"`
	ManifestDate        string `json:"manifest_date"`
	ManifestTime        string `json:"manifest_time"`
	CityName            string `json:"city_name"`
}

type Summary struct {
	CourierCode   string `json:"courier_code"`
	CourierName   string `json:"courier_name"`
	WaybillNumber string `json:"waybill_number"`
	ServiceCode   string `json:"service_code"`
	WaybillDate   string `json:"waybill_date"`
	ShipperName   string `json:"shipper_name"`
	ReceiverName  string `json:"receiver_name"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Status        string `json:"status"`
}

type Status struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
}
