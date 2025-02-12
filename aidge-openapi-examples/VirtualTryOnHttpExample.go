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
	"encoding/json"
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

	// Call virtual try on submit
	apiName := "/ai/virtual/tryon"
	submitRequest := `{"requestParams":"[{\"clothesList\":[{\"imageUrl\":\"https://ae-pic-a1.aliexpress-media.com/kf/H7588ee37b7674fea814b55f2f516fda1z.jpg\",\"type\":\"tops\"}],\"model\":{\"base\":\"General\",\"gender\":\"female\",\"style\":\"universal_1\",\"body\":\"slim\"},\"viewType\":\"mixed\",\"inputQualityDetect\":0,\"generateCount\":4}]"}`
	submitResult, err := invokeApi(accessKeyName, accessKeySecret, apiName, apiDomain, submitRequest)
	if err != nil {
		fmt.Println("Error invoking API:", err)
		return
	}

	var submitResultJSON map[string]interface{}
	err = json.Unmarshal([]byte(submitResult), &submitResultJSON)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	data, _ := submitResultJSON["data"].(map[string]interface{})
	result, _ := data["result"].(map[string]interface{})
	taskID, _ := result["taskId"].(string)

	// Query task status
	queryApiName := "/ai/virtual/tryon-results"
	queryRequest := `{"task_id":"` + taskID + `"}`
	queryResult := ""
	for {
		queryResult, err = invokeApi(accessKeyName, accessKeySecret, queryApiName, apiDomain, queryRequest)
		if err != nil {
			fmt.Println("Error querying API:", err)
			return
		}

		var queryResultJSON map[string]interface{}
		err = json.Unmarshal([]byte(queryResult), &queryResultJSON)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		data, _ := queryResultJSON["data"].(map[string]interface{})
		taskStatus, _ := data["taskStatus"].(string)
		if taskStatus == "finished" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Final result for the virtual try on
	fmt.Println(queryResult)
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