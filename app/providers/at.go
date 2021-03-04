package providers

import (
	"amplifier/app/entities"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/revel/revel"
)

type (
	AfricasTalkingSender interface {
		Send(*entities.SendRequest) (*entities.ATResponse, error)
	}

	AppAfricasTalkingSender struct {
		apiURL  string
		apiUser string
		apiKey  string
	}
)

func NewAfricasTalkingSender() AfricasTalkingSender {
	return NewAfricasTalkingSenderWithParameters(
		os.Getenv("AT_URL"),
		os.Getenv("AT_USER"),
		os.Getenv("AT_KEY"),
	)
}

func NewAfricasTalkingSenderWithParameters(
	apiURL string,
	apiUser string,
	apiKey string,
) *AppAfricasTalkingSender {
	return &AppAfricasTalkingSender{
		apiURL:  apiURL,
		apiUser: apiUser,
		apiKey:  apiKey,
	}
}

func (p *AppAfricasTalkingSender) Send(
	req *entities.SendRequest,
) (*entities.ATResponse, error) {
	atNumbers := make([]string, len(req.To))
	for i, rec := range req.To {
		atNumbers[i] = rec.Phone
	}

	reqBody := map[string]string{
		"username": p.apiUser,
		"message":  req.Message,
		"to":       strings.Join(atNumbers, ","),
		"from":     req.SenderID,
	}
	revel.AppLog.Infof("AT_REQ: %+v", reqBody)

	respStatus, respBody, err := makeFormRequest(
		p.apiURL,
		reqBody,
		"",
		"POST",
		map[string]string{"apikey": p.apiKey},
	)

	if err != nil {
		return &entities.ATResponse{
			Message: "atRequestFailed",
		}, err
	}

	if respStatus < 200 || respStatus > 299 {
		return &entities.ATResponse{
			Message: "atBadStatus",
		}, fmt.Errorf("AT not Ok status =[%v], body =[%v]", respStatus, string(respBody))
	}
	revel.AppLog.Infof("AT_RESP: %v", string(respBody))

	var retData entities.ATResponse
	err = json.Unmarshal(respBody, &retData)
	if err != nil {
		return &entities.ATResponse{
			Message: "atRequestUnmarshalFailed",
		}, err
	}
	return &retData, nil
}

func makeJsonRequest(
	targetUrl string,
	reqBody []byte,
	basicAuth string,
	requestType string,
	reqHeaders map[string]string,
) (int, []byte, error) {

	req, err := http.NewRequest(requestType, targetUrl, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if len(reqHeaders) > 0 {
		for key, val := range reqHeaders {
			req.Header.Add(key, val)
		}
	}
	authData := strings.Split(basicAuth, ":")
	if len(authData) > 1 {
		req.SetBasicAuth(authData[0], authData[1])
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, []byte(""), fmt.Errorf("client do: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, []byte(""), fmt.Errorf("error reading body: %v", err)
	}

	respStatus, err := strconv.Atoi(resp.Status)
	if err != nil {
		respStatus = 0
	}

	return respStatus, body, nil
}

func makeFormRequest(
	targetUrl string,
	reqBody map[string]string,
	basicAuth string,
	requestType string,
	reqHeaders map[string]string,
) (int, []byte, error) {

	form := url.Values{}
	for key, value := range reqBody {
		form.Add(key, value)
	}

	req, err := http.NewRequest(requestType, targetUrl, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Length", strconv.Itoa(len(form)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	if len(reqHeaders) > 0 {
		for key, val := range reqHeaders {
			req.Header.Add(key, val)
		}
	}
	authData := strings.Split(basicAuth, ":")
	if len(authData) > 1 {
		req.SetBasicAuth(authData[0], authData[1])
	}
	client := http.Client{Timeout: time.Second * 60 * 2}
	resp, err := client.Do(req)
	if err != nil {
		return 0, []byte(""), fmt.Errorf("client do: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, []byte(""), fmt.Errorf("error reading body: %v", err)
	}

	return resp.StatusCode, body, nil
}
