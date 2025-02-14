/*
Copyright (C) 2024 NEURALNETICS PTE. LTD.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"strings"
)

func main() {
	// Personal data from environment variables
	accessKeyName := os.Getenv("accessKey") // e.g. "512345"
	accessKeySecret := os.Getenv("secret")
	fmt.Println(accessKeyName)
	fmt.Println(accessKeySecret)

	apiDomain := "api.aidc-ai.com" // cn-api.aidc-ai.com for cn region

	// Call api
	apiName := "/ai/text/marco/translator"
	reqeust := "{\"text\":\"[\\\"Pen for iPad, 13 mins Fast Charging Stylus with Palm Rejection, Tilt Sensitivity, Compatible with 2018-2022 iPad Air 3/4/5, iPad Mini 5/6(Black)\\\"]\",\"sourceLanguage\":\"en\",\"targetLanguage\":\"ko\",\"formatType\":\"text\"}"
	result, _ := invokeApi(accessKeyName, accessKeySecret, apiName, apiDomain, reqeust)

	// Final result for the virtual try on
	fmt.Println(result)
}

func invokeApi(accessKeyName, accessKeySecret, apiName, apiDomain, data string) (string, error) {
	// Basic URL (placeholders included)
	urlTemplate := "https://%s/rest%s?partner_id=aidge&sign_method=sha256&sign_ver=v2&app_key=%s&timestamp=%s&sign=%s"

	// Timestamp in milliseconds
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	// Calculate SHA256 HMAC
	h := hmac.New(sha256.New, []byte(accessKeySecret))
	h.Write([]byte(accessKeySecret + timestamp))
	sign := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	// Create the final URL with real values
	finalURL := fmt.Sprintf(urlTemplate, apiDomain, apiName, accessKeyName, timestamp, sign)

	// Add "x-iop-trial": "true" for trial
	headers := map[string]string{
		"Content-Type": "application/json",
// 		"x-iop-trial": "true",
	}

	// Do HTTP POST request
	response, err := makeRequest("POST", finalURL, data, headers)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return "", err
	}
	fmt.Printf("Response: %s\n", response)
	return response, err
}

// makeRequest handles the HTTP request to the specified URL with the given data and headers
func makeRequest(method, url, data string, headers map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}