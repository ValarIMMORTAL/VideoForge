package cloud

import (
	"context"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"log"
	"sync"
)

var (
	QNManager *QiNiu
	once      sync.Once // 确保只初始化一次
)

func InitQiNiu(store db.Store) {
	thirdKey, err := store.GetThirdKeyByName(context.Background(), "qiniu")
	if err != nil {
		log.Println("获取七牛云密钥失败")
		return
	}
	once.Do(func() {
		QNManager = NewQiNiu(thirdKey.Ak, thirdKey.Sk)
	})
}
