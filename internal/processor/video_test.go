package processor

import (
	"context"
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	ctx := context.Background()

	var arg = VideoParams{}
	taskid, err := GenerateVideo(ctx, arg)
	if err != nil {
		return
	}
	fmt.Println("qidong")

	//模拟主协程永不退出
	for {

	}
	fmt.Println(taskid)
}
