[English](./README.md) | 简体中文

<p align="center">

<h1 align="center">Aidge API Go 示例</h1>

Aidge API Go 示例为您提供了示例代码，用于访问包括文本翻译在内的Aidge API。

## 环境要求

- 要运行示例，您必须拥有 Aidge API 帐户以及 `API key name` 和 `API key secret`。您可以在 Aidge 管理后台上创建并查看您的 API key信息。您可以联系您的服务
- 要使用Aidge API 示例访问产品的 API，您必须先在 [Aidge 控制台](https://www.aidge.com) 上激活该产品。

## 快速使用

以下这个代码示例向您展示了访问Aidge API的核心代码。

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

> 出于安全原因，我们不建议在源代码中硬编码凭据信息。您应该从外部配置或环境变量访问凭据。

## Changelog

每个版本的详细更改都记录在 [release notes](./ChangeLog.txt).


## References

- [Aidge官方网站](https://www.aidge.com/)

## License

This project is licensed under [Apache License Version 2](./LICENSE-2.0.txt) (SPDX-License-identifier: Apache-2.0).
