package cloud

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/global"
	"log"
	"testing"
)

func TestQiNiu(t *testing.T) {
	loadConfig, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := pgx.Connect(global.GlobalCtx, loadConfig.DBSource)

	if err != nil {
		log.Fatal("connect postgres err ", err)
	}
	q := db.New(conn)
	InitQiNiu(q)

	//localFile := "/Users/a0000/PycharmProjects/MoneyPrinterTurbo/storage/tasks/aacefe25-c9e3-418e-83b6-296f622277ba/combined-1.mp4"
	bucketName := "videofore-videos"
	//fileName := "运动.mp4"
	//ObjectName := "运动.mp4"
	//domain := "http://su15t494p.hn-bkt.clouddn.com/"
	//err = QNManager.UploadFile(context.Background(), localFile, bucketName, fileName, ObjectName)
	//if err != nil {
	//	fmt.Println("上传错误", err)
	//	return
	//}
	//publicURL := storage.MakePublicURL(domain, ObjectName)
	//fmt.Println("外链地址:", publicURL)
	_, err = QNManager.GetAllFileByUser(context.Background(), "ssaicyo", bucketName)
	if err != nil {
		fmt.Println("获取", err)
		return
	}
}
