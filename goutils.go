package goutils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

/*
GenericHTTPRequest ...
Function to make a generic HTTP Request
*/
func GenericHTTPRequest(requestType string, jsonPayloadStr string, url string, headers map[string]string) (responseBody string, httpResponseCode int, httpResponse *http.Response, err error) {
	var req *http.Request
	var errorMessage string
	var payload *bytes.Buffer
	httpResponseCode = -1
	debugMessage := ""
	responseBody = ""

	defer func() {
		if r := recover(); r != nil {
			debugMessage = fmt.Sprintf("MakeHTTPRequest : Recovered in f : %v", r)
			fmt.Println(debugMessage)
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("defer : unknown panic in MakeHTTPRequest")
			}
		}
	}()

	debugMessage = fmt.Sprintf("MakeHTTPRequest: requestType: (%v) , jsonPayloadStr: (%v) , url: (%v) , headers: (%v)", requestType, jsonPayloadStr, url, headers)
	fmt.Println(debugMessage)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   0,
			KeepAlive: 0,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if jsonPayloadStr != "" {
		payloadBytes := []byte(jsonPayloadStr)
		payload = bytes.NewBuffer(payloadBytes)
	} else {
		payload = nil
	}

	req, err = http.NewRequest(requestType, url, payload)
	if err != nil {
		errorMessage = fmt.Sprintf("ERROR : Could not form a new HTTP Request : %v", err.Error())
		fmt.Println(errorMessage)
		return responseBody, httpResponseCode, nil, errors.New(errorMessage)
	}

	if len(headers) != 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	req.Header.Set("Connection", "close")

	client := &http.Client{Transport: tr}

	httpResponse, err = client.Do(req)
	if err != nil {
		errorMessage = fmt.Sprintf("ERROR : Making HTTP Request : %v", err.Error())
		fmt.Println(errorMessage)
		return responseBody, httpResponseCode, nil, errors.New(errorMessage)
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		errorMessage = fmt.Sprintf("ERROR : Reading data from HTTP Request : %v", err.Error())
		fmt.Println(errorMessage)
		return responseBody, httpResponseCode, nil, errors.New(errorMessage)
	}

	httpResponseCode = httpResponse.StatusCode
	// convert body to string
	responseBody = string(body)
	return responseBody, httpResponseCode, httpResponse, nil
}

/*
Add ...
*/
func Add(num1 int, num2 int) (sum int) {
	sum = num1 + num2
	return sum
}
