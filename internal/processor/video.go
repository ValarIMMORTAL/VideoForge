package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pule1234/VideoForge/config"
	"log"
	"path"
	"time"
)

const (
	TaskStateComplete = 1 // 视频生成结束
)

// 调用MoneyPrinterTurbo的/api/v1/videos（post） 及api/v1/tasks接口（get）
func GenerateVideo(ctx context.Context, params VideoParams) (string, error) {
	conf, _ := config.LoadConfig("../../")
	//todo 请求/api/v1/videos接口 将返回的taskid 传到api/v1/tasks中
	//videoUrl := "http://127.0.0.1:8080/api/v1/videos"
	videoUrl := path.Join(conf.GenerateVideoBaseUrl, conf.VideoEndpoint)

	allRequest, err := SendPostRequest(params, videoUrl, conf)
	if err != nil {
		return "", fmt.Errorf("向 %s 发送 POST 请求失败: %w", videoUrl, err)
	}
	var videoResp VideoResults
	_ = json.Unmarshal(allRequest, &videoResp)

	taskId := videoResp.Data.TaskId
	//启动协程，将轮询api/v1/tasks接口，当生成完毕时返回
	// todo 创建get api 请求
	//taskUrl := "http://127.0.0.1:8080/api/v1/tasks" + "/" + taskId
	taskUrl := path.Join(conf.GenerateVideoBaseUrl, conf.TaskEndpoint, taskId)

	go func() {
		//采用指数回避的方式轮询改接口
		backoff := 10 * time.Second
		maxBackoff := 30 * time.Second
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
				resp, err := SendGetRequest(taskUrl, conf)
				if err != nil {
					log.Println("Generate video failed: " + err.Error())
					return
				}
				var taskResp TaskResult
				_ = json.Unmarshal(resp, &taskResp)
				fmt.Println("轮询中，当前 state:", taskResp.Data.State)
				if taskResp.Data.State == TaskStateComplete {
					log.Println("任务结束，视频生成完成")
					// todo 通知视频生成完成

					return
				}
				backoff = min(backoff*2, maxBackoff)
			}
		}
	}()

	return taskId, nil
}
