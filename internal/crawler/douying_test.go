package crawler

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"testing"
)

func TestStart(t *testing.T) {
	loadConfig, err := config.LoadConfig("../../")
	if err != nil {
		fmt.Println(err)
	}

	conn, err := pgx.Connect(context.Background(), loadConfig.DBSource)
	if err != nil {
		t.Error("connect postgres err : " + err.Error())
	}
	q := db.New(conn)
	crawler, err := newDyCrawler(loadConfig.DouYingQueueName, q)
	if err != nil {
		fmt.Println(err)
	}

	err = crawler.Start("https://www.douyin.com/aweme/v1/web/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web")
	if err != nil {
		fmt.Println(err)
	}
}
