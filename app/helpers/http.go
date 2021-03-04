package helpers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HttpReqData struct {
	Method  string
	URL     string
	Auth    string
	Body    map[string]string
	Headers map[string]string
}

type HttpResponse struct {
	StatusCode   int
	Header       http.Header
	Status, Body string
}

func (httpReq *HttpReqData) MakeHTTPRequest() (HttpResponse, error) {
	if len(httpReq.Body) < 0 {
		return HttpResponse{}, errors.New("No form body found")
	}
	form := url.Values{}
	for key, value := range httpReq.Body {
		form.Add(key, value)
	}
	client := http.Client{Timeout: time.Second * 60 * 3}
	req, err := http.NewRequest(
		httpReq.Method, httpReq.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return HttpResponse{}, fmt.Errorf("makerequest: %v", err)
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(form)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	if len(httpReq.Headers) > 0 {
		for key, val := range httpReq.Headers {
			req.Header.Add(key, val)
		}
	}
	if len(httpReq.Auth) > 1 {
		user := strings.Split(httpReq.Auth, ":")
		req.SetBasicAuth(user[0], user[1])
	}
	resp, err := client.Do(req)
	if err != nil {
		return HttpResponse{}, fmt.Errorf("makerequest do: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return HttpResponse{}, fmt.Errorf("makerequest readall: %v", err)
	}
	return HttpResponse{
		Body: string(body), Header: resp.Header, Status: resp.Status, StatusCode: resp.StatusCode,
	}, nil
}
