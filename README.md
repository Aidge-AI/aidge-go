English | [简体中文](./README-CN.md)

<p align="center">

<h1 align="center">Aidge API Examples for Go</h1>

The Aidge API examples for Go provide you  to access Aidge services such as Text Translation.

## Requirements

- To run the examples, you must have an Aidge API account as well as an `API Key Name` and an `API Key Secret`. Create and view your AccessKey on Aidge dashboard.
- To use the Aidge API examples for Go to access the APIs of a product, you must first activate the product on the [Aidge console](https://www.aidge.com) if required.

## Quick Examples

The following code example:

```go
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Your personal data
	accessKeyName := "your access key name" // e.g. 512345
	accessKeySecret := "your access key secret"
	apiName := "api name"     // e.g. /ai/text/translation/and/polishment
	apiDomain := "api domain" // e.g. api.aidc-ai.com or cn-api.aidc-ai.com
	data := "{\"requestParams\":\"your api request params\"}"

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

	// Headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Do HTTP POST request
	response, err := makeRequest("POST", finalURL, data, headers)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return
	}
	fmt.Printf("Response: %s\n", response)
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

```

> For security reason, we don't recommend to hard code credentials information in source code. You should access
> credentials from external configurations or environment variables.

## Changelog

Detailed changes for each release are documented in the [release notes](./ChangeLog.txt).


## References

- [Aidge Home Page](https://www.aidge.com/)

## License

This project is licensed under [Apache License Version 2](./LICENSE-2.0.txt) (SPDX-License-identifier: Apache-2.0).
