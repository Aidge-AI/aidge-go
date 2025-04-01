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
	"strings"
	"time"
	"os"
)

// ApiConfig holds the API configuration
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

func invokeAPI(apiName string, data interface{}, urlParam string, isGet bool) (string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	// Calculate sha256 sign
	signString := apiConfig.AccessKeySecret + timestamp
	h := hmac.New(sha256.New, []byte(apiConfig.AccessKeySecret))
	h.Write([]byte(signString))
	sign := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	url := fmt.Sprintf("https://%s/rest%s?partner_id=aidge&sign_method=sha256&sign_ver=v2&app_key=%s&timestamp=%s&sign=%s",
		apiConfig.ApiDomain, apiName, apiConfig.AccessKeyName, timestamp, sign)

	// Add "x-iop-trial": "true" for trial
	headers := map[string]string{
		"Content-Type": "application/json",
		"x-iop-trial":  strings.ToLower(strconv.FormatBool(apiConfig.UseTrialResource)),
	}

	var jsonBody []byte
	var err error

	if data != nil {
		jsonBody, err = json.Marshal(data)
		if err != nil {
			return "", err
		}
	}

	client := &http.Client{}
	var req *http.Request
	if isGet {
		req, err = http.NewRequest("GET", url + "&" + urlParam, bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	}

	if err != nil {
		return "", err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
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

func main() {

	// Call submit api
	apiName := "/ai/image/translation_mllm/batch"

	// Construct request parameters
	requestParams := []map[string]string{
		{
			"imageUrl":       "https://img.alicdn.com/imgextra/i1/1955749012/O1CN016P3Jas2GRY7vaevsK_!!1955749012.jpg",
			"sourceLanguage": "zh",
			"targetLanguage": "en",
		},
		{
			"imageUrl":       "https://img.alicdn.com/imgextra/i1/1955749012/O1CN016P3Jas2GRY7vaevsK_!!1955749012.jpg",
			"sourceLanguage": "zh",
			"targetLanguage": "ko",
		},
	}

	// Convert parameters to JSON string
	submitRequest := map[string]string{
		"paramJson": string(func() []byte {
			data, _ := json.Marshal(requestParams)
			return data
		}()),
	}

	submitResult, err := invokeAPI(apiName, submitRequest, "", false)
	if err != nil {
		fmt.Println("Error invoking API:", err)
		return
	}

	var submitResultJson map[string]interface{}
	if err := json.Unmarshal([]byte(submitResult), &submitResultJson); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	fmt.Println("submitResultJson:", submitResultJson)

	taskID := submitResultJson["data"].(map[string]interface{})["result"].(map[string]interface{})["taskId"].(string)

	// Query task status
	queryAPIName := "/ai/image/translation_mllm/results"
	var queryResult string
	for {
		queryResult, err = invokeAPI(queryAPIName, nil, "taskId=" + taskID, true)
		if err != nil {
			fmt.Println("Error querying API:", err)
			return
		}

		var queryResultJson map[string]interface{}
		if err := json.Unmarshal([]byte(queryResult), &queryResultJson); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}
		fmt.Println("queryResultJson:", queryResultJson)

		taskStatus := queryResultJson["data"].(map[string]interface{})["taskStatus"].(string)
		if taskStatus == "finished" {
			break
		}

		time.Sleep(5 * time.Second)
	}

	// Final result
	fmt.Println(queryResult)
}