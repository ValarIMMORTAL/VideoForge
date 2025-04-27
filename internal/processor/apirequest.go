package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pule1234/VideoForge/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 拼接请求URL
func BuildUrl(baseUrl, endPoint string, pathParams ...string) (string, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Invalid base URL: %s", err))
	}

	fullPath := strings.TrimRight(base.Path, "/") + "/" + strings.TrimRight(endPoint, "/")
	for _, pathParam := range pathParams {
		fullPath += "/" + strings.Trim(pathParam, "/")
	}

	base.Path = fullPath

	// 返回完整 URL
	return base.String(), nil
}

// 封装api请求的function
func SendPostRequest(
	requestData interface{},
	url string,
	conf *config.Config,
) ([]byte, error) {
	jsonRequest, _ := json.Marshal(requestData)
	payload := strings.NewReader(string(jsonRequest))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.New("failed to create copywriting request : " + err.Error())
	}

	req.Header.Add("Authorization", conf.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to push copywriting request: " + err.Error())
	}

	allResp, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed to push copywriting request: HTTP code %d", strconv.Itoa(res.StatusCode)))
	}

	if err != nil {
		return nil, err
	}

	return allResp, err
}

// get请求发送
func SendGetRequest(
	url string,
	conf *config.Config,
) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("SendGetRequest: failed to create request : " + err.Error())
	}
	req.Header.Add("Authorization", conf.ApiKey)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("SendGetRequest: failed to send request : " + err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to push copywriting request: HTTP code " + strconv.Itoa(resp.StatusCode))
	}

	allResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return allResp, err
}
