package gapi

import (
	"fmt"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/cloud"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/publisher"
	"github.com/pule1234/VideoForge/mq"
	"github.com/pule1234/VideoForge/pb"
	"github.com/pule1234/VideoForge/token"
)

type Server struct {
	pb.UnimplementedVideosForgeServer
	config config.Config //读取文件配置
	store  db.Store
	//tokenMaker token.Maker
	redis            *cache.Redis
	tokenMaker       token.Maker
	publisherFactory *publisher.PublisherFactory
	qnManager        *cloud.QiNiu
	mq               *mq.RabbitMQ
}

func NewServer(conf config.Config, store db.Store, factory *publisher.PublisherFactory) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:           conf,
		store:            store,
		redis:            cache.RedisClient,
		qnManager:        cloud.QNManager,
		publisherFactory: factory,
		tokenMaker:       tokenMaker,
		mq:               mq.GlobalRabbitMQ,
	}

	return server, nil
}
