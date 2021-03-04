package entities

type (

	// ATResponse struct for at response
	ATResponse struct {
		Message    string         `json:"Message"`
		Recipients []*ATRecipient `json:"Recipients"`
	}

	ATRecipient struct {
		Number    string `json:"number"`
		Status    string `json:"status"`
		Cost      string `json:"cost"`
		MessageID string `json:"messageId"`
	}

	SendRequest struct {
		To       []*SMSRecipient
		Message  string
		SenderID string
	}

	SMSRecipient struct {
		Phone      string  `json:"phone"`
		Message    string  `json:"message"`
		Cost       float64 `json:"cost"`
		Correlator string  `json:"correlator"`
		Status     string  `json:"status"`
		Reason     string  `json:"reason"`
		SheetRow   int64   `json:"sheet_row"`
		Route      string  `json:"route"`
		IsValid    bool    `json:"is_valid"`
	}
)
