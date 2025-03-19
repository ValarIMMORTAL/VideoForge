package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pule1234/VideoForge/config"
	"log"
	"time"
)

// 调用MoneyPrinterTurbo的/api/v1/videos（post） 及api/v1/tasks接口（get）
func GenerateVideo(params VideoParams) error {
	conf, _ := config.LoadConfig("../../")
	//todo 请求/api/v1/videos接口 将返回的taskid 传到api/v1/tasks中
	videoUrl := "http://127.0.0.1:8080/api/v1/videos"

	allRequest, err := SendPostRequest(params, videoUrl, conf)
	if err != nil {
		return errors.New("Generate video failed: " + err.Error())
	}

	var videoResp VideoResults
	_ = json.Unmarshal(allRequest, &videoResp)

	taskId := videoResp.Data.TaskId
	//启动协程，将轮询api/v1/tasks接口，当生成完毕时返回
	// todo 创建get api 请求
	taskUrl := "http://127.0.0.1:8080/api/v1/tasks" + "/" + taskId

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				resp, err := SendGetRequest(taskUrl, conf)
				if err != nil {
					log.Println("Generate video failed: " + err.Error())
					return
				}
				var taskResp TaskResult
				_ = json.Unmarshal(resp, &taskResp)
				fmt.Println("轮询中，当前 state:", taskResp.Data.State)
				if taskResp.Data.State == 1 {
					log.Println("任务结束，视频生成完成")
					// todo 通知视频生成完成

					return
				}
			}
		}
	}()

	return nil
}
