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
		To           []*SMSRecipient
		Multi        bool
		Message      string
		SenderID     string
		StatusUrl    string
		StatusSecret string
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

	KomsnerResponse struct {
		Success    bool                `json:"success"`
		Status     int                 `json:"status"`
		Message    string              `json:"message"`
		Recipients []*KomsnerRecipient `json:"recipients"`
	}

	KomsnerRecipient struct {
		Number string  `json:"number"`
		Status string  `json:"status"`
		Cost   float64 `json:"cost"`
		ID     string  `json:"id"`
	}
)
