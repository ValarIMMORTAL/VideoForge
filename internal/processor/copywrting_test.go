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
	var item = models.TrendingItem{
		Source: "douying",
	}
	defer conn.Close(context.Background())

	q := db.New(conn)

	processor := NewProcessor(q)

	err = processor.CreateCopyWriting(item)
	if err != nil {
		return
	}
}
