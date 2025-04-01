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
	"strconv"
	"time"
	"os"
	"strings"
)

// ApiConfig struct holds configuration for API access
type ApiConfig struct {
	AccessKeyName   string
	AccessKeySecret string
	ApiDomain       string
	UseTrialResource bool
}

// Global variable to hold the API configuration
var apiConfig = ApiConfig{
    // Your personal data
	AccessKeyName:   os.Getenv("accessKey"), // e.g. "512345"
	AccessKeySecret: os.Getenv("secret"),

	/* "api.aidc-ai.com" for api purchased on global site
	 * 中文站购买的API请使用"cn-api.aidc-ai.com" (for api purchased on chinese site)
	 */
	ApiDomain:       "your api domain",

	/**
	 * We offer trial quota to help you familiarize and test how to use the Aidge API in your account
	 * To use trial quota, please set useTrialResource to true
	 * If you set useTrialResource to false before you purchase the API
	 * You will receive "Sorry, your calling resources have been exhausted........"
	 * 我们为您的账号提供一定数量的免费试用额度可以试用任何API。请将useTrialResource设置为true用于试用。
	 * 如设置为false，且您未购买该API，将会收到"Sorry, your calling resources have been exhausted........."的错误提示
	 */
	UseTrialResource: false / true,
}

// Function to invoke the API
func invokeAPI(apiName string, data []byte) (string, error) {
	// Timestamp in milliseconds
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	// Calculate SHA256 HMAC
	h := hmac.New(sha256.New, []byte(apiConfig.AccessKeySecret))
	h.Write([]byte(apiConfig.AccessKeySecret + timestamp))
	sign := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	url := fmt.Sprintf("https://%s/rest%s?partner_id=aidge&sign_method=sha256&sign_ver=v2&app_key=%s&timestamp=%s&sign=%s",
		apiConfig.ApiDomain, apiName, apiConfig.AccessKeyName, timestamp, sign)

	// Add "x-iop-trial": "true" for trial
	headers := map[string]string{
		"Content-Type": "application/json",
		"x-iop-trial":  strconv.FormatBool(apiConfig.UseTrialResource),
	}

	// HTTP request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	// Call api
	apiName := "/ai/image/translation_mllm/batch"

	// Constructor request Parameters
	paramJson := []map[string]interface{}{
		{
            "imageUrl": "https://img.alicdn.com/imgextra/i1/1955749012/O1CN016P3Jas2GRY7vaevsK_!!1955749012.jpg",
            "sourceLanguage": "zh",
            "targetLanguage": "en"
        },
        {
            "imageUrl": "https://img.alicdn.com/imgextra/i1/1955749012/O1CN016P3Jas2GRY7vaevsK_!!1955749012.jpg",
            "sourceLanguage": "zh",
            "targetLanguage": "ko"
        }
	}

	param := map[string]interface{}{
		"paramJson": paramJson,
	}

	// Convert parameters to JSON string
	request, err := json.Marshal(param)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	submitResult, err := invokeAPI(apiName, request)
	if err != nil {
		fmt.Println("Error invoking API:", err)
		return
	}

	var submitResultJson map[string]interface{}
	err = json.Unmarshal([]byte(submitResult), &submitResultJson)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	fmt.Println(submitResultJson)

	taskID := submitResultJson["data"].(map[string]interface{})["result"].(map[string]interface{})["taskId"].(string)

	// Query task status
	queryApiName := "/ai/image/translation_mllm/results"
	queryRequest, err := json.Marshal(map[string]string{"taskId": taskID})
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	for {
		queryResult, err := invokeAPI(queryApiName, queryRequest)
		if err != nil {
			fmt.Println("Error invoking API:", err)
			return
		}

		var queryResultJson map[string]interface{}
		err = json.Unmarshal([]byte(queryResult), &queryResultJson)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
		fmt.Println("queryResultJson=", queryResultJson)

		taskStatus := queryResultJson["data"].(map[string]interface{})["taskStatus"].(string)
		if taskStatus == "finished" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Add a small delay between requests to avoid overwhelming the API
	time.Sleep(1 * time.Second)
}
