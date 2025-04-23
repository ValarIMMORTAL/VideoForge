package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/cloud"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/mq"
	"log"
	"time"
)

const (
	TaskStateComplete = 1  // 视频生成结束
	TaskStateFale     = -1 // 视频生成失败
)

// 调用MoneyPrinterTurbo的/api/v1/videos（post） 及api/v1/tasks接口（get）
func GenerateVideo(
	ctx context.Context,
	params VideoParams,
	fileName string,
	userName string,
	userId int64,
	redis *cache.Redis,
	store db.Store,
	storage *cloud.QiNiu,
	msgQueue *mq.RabbitMQ,
) (string, error) {
	conf, _ := config.LoadConfig("../../")
	//videoUrl := "http://127.0.0.1:8080/api/v1/videos"
	timestamp := time.Now().Unix()
	videoUrl, err := BuildUrl(conf.GenerateVideoBaseUrl, conf.VideoEndpoint)
	if err != nil {
		return "", err
	}
	allRequest, err := SendPostRequest(params, videoUrl, conf)
	if err != nil {
		return "", fmt.Errorf("向 %s 发送 POST 请求失败: %w", videoUrl, err)
	}
	var videoResp VideoResults
	_ = json.Unmarshal(allRequest, &videoResp)

	taskId := videoResp.Data.TaskId
	taskUrl, err := BuildUrl(conf.GenerateVideoBaseUrl, conf.TaskEndpoint, taskId)
	queueName := "GenerateVideo" + userName + "_" + fmt.Sprintf("%d", userId)
	if err != nil {
		return "", err
	}
	go func() {
		//backoff := 10 * time.Second
		//maxBackoff := 30 * time.Second
		for {
			select {
			case <-ctx.Done():
				return
			//case <-time.After(backoff):
			case <-time.After(10 * time.Second):
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
					redis.SRem(ctx, userName, taskId)
					localPath := conf.VideoPath + "/" + taskId + "/combined-1.mp4"
					err = storage.UploadFile(ctx, localPath, "videofore-videos", params.VideoSubject+".mp4", fileName+userName+fmt.Sprintf("%d", timestamp)+".mp4")
					if err != nil {
						log.Println("七牛云上传出错" + err.Error())
						item := VideoMsg{
							Event:   "GenerateVideo failed" + err.Error(),
							User_id: userId,
						}
						err = msgQueue.PublishItem(item, queueName)
						if err != nil {
							log.Println("发送通知失败", err)
						}
						return
					}
					externalLink := conf.CdnDomain + params.VideoSubject + ".mp4"
					arg := db.InsertVideoParams{
						Title:     params.VideoSubject,
						Url:       externalLink,
						Duration:  taskResp.Data.AudioDuration,
						UserID:    userId,
						Subscribe: timestamp,
					}
					video, err := store.InsertVideo(ctx, arg)
					if err != nil {
						log.Println("视频和用户绑定失败", err)
						return
					}
					item := VideoMsg{
						Event:    "GenerateVideo Success",
						Video_id: video.ID,
						User_id:  video.UserID,
						Url:      video.Url,
					}
					err = msgQueue.PublishItem(item, queueName)
					if err != nil {
						log.Println("发送通知失败", err)
					}
					return
				} else if taskResp.Data.State == TaskStateFale {
					item := VideoMsg{
						Event:   "GenerateVideo failed",
						User_id: userId,
					}
					err = msgQueue.PublishItem(item, queueName)
					if err != nil {
						log.Println("发送通知失败", err)
					}
					return
				}
				//backoff = min(backoff*2, maxBackoff)
			}
		}
	}()

	return taskId, nil
}
