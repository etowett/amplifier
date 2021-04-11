package providers

import (
	"amplifier/app/entities"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/revel/revel"
)

type (
	KomsnerSender interface {
		Send(*entities.SendRequest) (*entities.KomsnerResponse, error)
	}

	AppKomsnerSender struct {
		apiURL       string
		apiBasicAuth string
	}
)

func NewKomsnerSender() KomsnerSender {
	return NewKomsnerSenderWithParameters(
		os.Getenv("KOMSNER_URL"),
		os.Getenv("KOMSNER_AUTH"),
	)
}

func NewKomsnerSenderWithParameters(
	apiURL string,
	apiBasicAuth string,
) *AppKomsnerSender {
	return &AppKomsnerSender{
		apiURL:       apiURL,
		apiBasicAuth: apiBasicAuth,
	}
}

func (p *AppKomsnerSender) Send(
	req *entities.SendRequest,
) (*entities.KomsnerResponse, error) {
	recs := make([]map[string]string, 0)

	for _, rec := range req.To {
		recs = append(recs, map[string]string{
			"number":  rec.Phone,
			"message": rec.Message,
		})
	}

	requestBody := map[string]interface{}{
		"message":       req.Message,
		"multi":         req.Multi,
		"source":        req.SenderID,
		"destination":   recs,
		"status_url":    req.StatusUrl,
		"status_secret": req.StatusSecret,
	}

	jsonRequest, err := NewJSONRequest(
		p.apiURL,
		requestBody,
		p.apiBasicAuth,
		nil,
	)
	if err != nil {
		return &entities.KomsnerResponse{
			Message: "komsnerRequestFailed",
		}, fmt.Errorf("could not create new komsner json request: %v", err)
	}
	revel.AppLog.Infof("Komsner Request: %+v", jsonRequest)

	resp, err := Do(jsonRequest)
	if err != nil {
		return &entities.KomsnerResponse{
			Message: "komsnerRequestFailed",
		}, fmt.Errorf("komsner make request error: %v", err)
	}

	log.Printf("Komsner response: %+v", resp)

	atNumbers := make([]string, len(req.To))
	for i, rec := range req.To {
		atNumbers[i] = rec.Phone
	}

	var retData entities.KomsnerResponse
	err = json.Unmarshal([]byte(resp.Body), &retData)
	if err != nil {
		return &entities.KomsnerResponse{
			Message: "komsnerRequestUnmarshalFailed",
		}, err
	}
	return &retData, nil
}
