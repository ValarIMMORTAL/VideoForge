package processor

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/models"
	"testing"
)

func TestCreateCopyWriting(t *testing.T) {
	loadConfig, err := config.LoadConfig("../../")
	if err != nil {
		return
	}

	conn, err := pgx.Connect(context.Background(), loadConfig.DBSource)
	if err != nil {
		t.Error("connect postgres err : " + err.Error())
	}
	var items = []models.TrendingItem{}
	items = append(items, models.TrendingItem{
		Source: "", // 填写参数
		Title:  "", // 填写参数
	})
	defer conn.Close(context.Background())

	q := db.New(conn)

	err = CreateCopyWriting(items, q)
	if err != nil {
		return
	}

}
