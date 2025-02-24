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

	apiDomain := "api.aidc-ai.com" // for api purchased on global site
	// apiDomain := "cn-api.aidc-ai.com" // 中文站购买的API请使用此域名 (for api purchased on chinese site)

	// Call api
	apiName := "/ai/image/removal"
	reqeust := "{\"image_url\":\"https://ae01.alicdn.com/kf/Sa78257f1d9a34dad8ee494178db12ec8l.jpg\",\"non_object_remove_elements\":\"[1,2,3,4]\",\"object_remove_elements\":\"[1,2,3,4]\",\"mask\":\"474556 160 475356 160 476156 160 476956 160 477756 160 478556 160 479356 160 480156 160 480956 160 481756 160 482556 160 483356 160 484156 160 484956 160 485756 160 486556 160 487356 160 488156 160 488956 160 489756 160 490556 160 491356 160 492156  160\"}"
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