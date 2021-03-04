package forms

type (
	ATForm struct {
		Count   int64  `json:"count"`
		Message string `json:"message"`
		Multi   bool   `json:"multi"`
	}
)
