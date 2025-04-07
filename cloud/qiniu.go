package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

type QiNiu struct {
	UploadManager  *uploader.UploadManager
	ObjectsManager *objects.ObjectsManager
}

func NewQiNiu(accessKey, secretKey string) *QiNiu {
	mac := credentials.NewCredentials(accessKey, secretKey)
	upload := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})

	object := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: mac},
	})
	return &QiNiu{
		UploadManager:  upload,
		ObjectsManager: object,
	}
}

// 上传单个文件
func (q *QiNiu) UploadFile(ctx context.Context, localFile, bucketName, fileName, objectsName string) error {
	err := q.UploadManager.UploadFile(ctx, localFile, &uploader.ObjectOptions{
		BucketName: bucketName,
		ObjectName: &objectsName,
		FileName:   fileName,
		CustomVars: map[string]string{
			"name": fileName,
		},
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// 获取用户的所有视频
func (q *QiNiu) GetAllFileByUser(ctx context.Context, user, bucketName string) ([]string, error) {
	bucket := q.ObjectsManager.Bucket(bucketName)

	// 回调函数
	onResponse := func(od *objects.ObjectDetails) {
		marshal, err := json.Marshal(od)
		if err != nil {
			return
		}
		fmt.Println(string(marshal))
		fmt.Printf("%s: %d bytes\n", od.Name, od.Size)
	}

	onError := func(err error) {
		fmt.Printf("Error: %v\n", err)
	}

	//获取文件列表
	fileNames := make([]string, 0)
	fileNames = append(fileNames, "为什么要运动.mp4")

	ops := make([]objects.Operation, len(fileNames))
	for i, name := range fileNames {
		ops[i] = bucket.Object(name).Stat().
			OnResponse(onResponse).
			OnError(onError)
	}
	if err := q.ObjectsManager.Batch(context.Background(), ops, nil); err != nil {
		return nil, err
	}
	return nil, nil
}
