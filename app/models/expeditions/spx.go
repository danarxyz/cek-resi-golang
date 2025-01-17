package models

type SpxModel struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		SlsTrackingNumber string `json:"sls_tracking_number"`
		NeedTranslate     int    `json:"need_translate"`
		DeliveryType      string `json:"delivery_type"`
		RecipientName     string `json:"recipient_name"`
		Phone             string `json:"phone"`
		CurrentStatus     string `json:"current_status"`
		TrackingList      []struct {
			Timestamp int    `json:"timestamp"`
			Status    string `json:"status"`
			Message   string `json:"message"`
		} `json:"tracking_list"`
		StatusList []struct {
			Timestamp int    `json:"timestamp"`
			Code      int    `json:"code"`
			Text      string `json:"text"`
			StateLs   string `json:"state_ls"`
			Icon      string `json:"icon"`
		} `json:"status_list"`
	} `json:"data"`
}
