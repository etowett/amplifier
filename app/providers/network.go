package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	BasicAuth string
	Body      io.Reader
	Header    http.Header
	Method    string
	timeout   time.Duration
	URL       string
}

type Response struct {
	Body       string
	Header     http.Header
	Status     string
	StatusCode int
}

func NewJSONRequest(
	reqUrl string,
	body map[string]interface{},
	basicAuth string,
	headers map[string]string,
) (*Request, error) {

	jsonStr, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("body json marshal: %v", err)
	}

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")

	if len(headers) > 0 {
		for key, val := range headers {
			header.Add(key, val)
		}
	}

	req := &Request{
		BasicAuth: basicAuth,
		Body:      bytes.NewBuffer(jsonStr),
		Header:    header,
		Method:    http.MethodPost,
		timeout:   time.Second * 60 * 2,
		URL:       reqUrl,
	}

	return req, nil
}

func Do(req *Request) (*Response, error) {

	httpClient := http.Client{
		Timeout: req.timeout,
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("new request error: %v", err)
	}

	httpReq.Close = true
	httpReq.Header = req.Header

	if len(req.BasicAuth) > 1 {
		user := strings.Split(req.BasicAuth, ":")
		httpReq.SetBasicAuth(user[0], user[1])
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("client do error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read error: %v", err)
	}

	response := &Response{
		Body:       string(body),
		Header:     resp.Header,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}

	return response, nil
}
