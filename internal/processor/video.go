package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/cloud"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"log"
	"time"
)

const (
	TaskStateComplete = 1 // 视频生成结束
)

// 调用MoneyPrinterTurbo的/api/v1/videos（post） 及api/v1/tasks接口（get）
func GenerateVideo(ctx context.Context, params VideoParams, userName string, userId int64, redis *cache.Redis, store db.Store, storage *cloud.QiNiu) (string, error) {
	conf, _ := config.LoadConfig("../../")
	//videoUrl := "http://127.0.0.1:8080/api/v1/videos"

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
	//启动协程，将轮询api/v1/tasks接口，当生成完毕时返回
	//taskUrl := "http://127.0.0.1:8080/api/v1/tasks" + "/" + taskId
	//taskId := "dcd22394-3476-4c15-8206-eca05d60121c" //若需要测试异步轮询，需要将改行放开注释，将改行上面的代码注释（conf不需要注释）
	taskUrl, err := BuildUrl(conf.GenerateVideoBaseUrl, conf.TaskEndpoint, taskId)
	if err != nil {
		return "", err
	}
	go func() {
		//采用指数回避的方式轮询改接口
		//backoff := 10 * time.Second
		//maxBackoff := 30 * time.Second
		for {
			select {
			case <-ctx.Done():
				fmt.Println("协程退出")
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
					// todo 通知视频生成完成， 需要传递当前用户的信息
					redis.SRem(ctx, userName, taskId) //视频生成完成、将当前user的task从redis中删除
					//视频存储到七牛云
					localPath := conf.VideoPath + "/" + taskId + "/combined-1.mp4"
					err = storage.UploadFile(ctx, localPath, "videofore-videos", params.VideoSubject+".mp4", params.VideoSubject+".mp4")
					if err != nil {
						log.Println("七牛云上传出错" + err.Error())
						return
					}
					externalLink := conf.CdnDomain + params.VideoSubject + ".mp4"
					arg := db.InsertVideoParams{
						Title:    params.VideoSubject,
						Url:      externalLink,
						Duration: taskResp.Data.AudioDuration,
						UserID:   userId,
					}
					_, err := store.InsertVideo(ctx, arg)
					if err != nil {
						log.Println("视频和用户绑定失败", err)
						return
					}
					return
				}
				//backoff = min(backoff*2, maxBackoff)
			}
		}
	}()

	return taskId, nil
}
