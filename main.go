package main

import (
	"encoding/json"
	"fmt"
	"github.com/pule1234/VideoForge/internal/processor"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	url := "https://api.chatanywhere.tech/v1/chat/completions"
	method := "POST"

	payload := strings.NewReader(`{"messages": [{"content": "人工智能兴起，帮我对这个关键字生成短视频文案","role": "system"}],"model": "gpt-3.5-turbo"}`)
	fmt.Println(payload)
	requestData := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "人工智能兴起，帮我对这个关键字生成短视频文案",
			},
		},
	}

	jsonRequest, _ := json.Marshal(requestData)
	payload = strings.NewReader(string(jsonRequest))

	fmt.Println(payload)
	//return

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer sk-E2unVpWd37R5HkoSClednWZkIKd2D3HnOS0Ewy9ydjRIYDgi")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var str processor.Response
	err = json.Unmarshal(body, &str)
	if err != nil {
		return
	}
	fmt.Println(res.StatusCode)
	fmt.Println(string(body))
	fmt.Println(str.Choices[0].Message.Content)
}
