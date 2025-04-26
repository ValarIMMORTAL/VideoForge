package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/cloud"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/publisher"
	"github.com/pule1234/VideoForge/mq"
	"github.com/pule1234/VideoForge/pb"
	"github.com/pule1234/VideoForge/token"
	"github.com/pule1234/VideoForge/worker"
	"google.golang.org/grpc"
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
	mq               *mq.RabbitMQ
	grpcClient       pb.VideosForgeClient
	taskDistributor  *worker.TaskDistributor
	taskprocessor    *worker.TaskProcessor
}

func NewServer(conf config.Config, store db.Store, factory *publisher.PublisherFactory, taskDistributor *worker.TaskDistributor, taskprocessor *worker.TaskProcessor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	fmt.Println("ready grpc client success")
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	fmt.Println("connect grpc client success")
	server := &Server{
		config:           conf,
		store:            store,
		router:           gin.Default(),
		redis:            cache.RedisClient,
		qnManager:        cloud.QNManager,
		publisherFactory: factory,
		tokenMaker:       tokenMaker,
		mq:               mq.GlobalRabbitMQ,
		grpcClient:       pb.NewVideosForgeClient(conn),
		taskDistributor:  taskDistributor,
		taskprocessor:    taskprocessor,
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
	router.GET("/ping")
	//加入token

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker)) // 使用中间件进行认证
	authRoutes.GET("/getCopyWritings", server.getCopyWriting)
	authRoutes.GET("/getVideos", server.getVideos)
	authRoutes.POST("/generateVideo", server.generateVideo)
	authRoutes.POST("/upload-video", server.UploadVideo)
	authRoutes.POST("/upload-video-grpc", server.UploadVideoGrpc)
	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
