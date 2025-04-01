package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/cloud"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/publisher"
	"github.com/pule1234/VideoForge/token"
)

type Server struct {
	config config.Config //读取文件配置
	store  db.Store
	//tokenMaker token.Maker
	router           *gin.Engine
	redis            *cache.Redis
	tokenMaker       token.Maker
	publisherFactory *publisher.PublisherFactory
	qnManager        *cloud.QiNiu
}

func NewServer(conf config.Config, store db.Store, factory *publisher.PublisherFactory) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:           conf,
		store:            store,
		router:           gin.Default(),
		redis:            cache.RedisClient,
		qnManager:        cloud.QNManager,
		publisherFactory: factory,
		tokenMaker:       tokenMaker,
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	router.GET("/ping", server.callback)
	//加入token

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker)) // 使用中间件进行认证
	authRoutes.POST("/generateVideo", server.generateVideo)
	authRoutes.POST("/upload-video", server.UploadVideo)
	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
