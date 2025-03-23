package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/token"
	"net/http"
)

type Server struct {
	config config.Config //读取文件配置
	store  db.Store
	//tokenMaker token.Maker
	router     *gin.Engine
	redis      *cache.Redis
	tokenMaker token.Maker
}

func NewServer(conf config.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: conf,
		store:  store,
		router: gin.Default(),
		redis:  cache.RedisClient,
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//todo route定义
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	//加入token

	router.POST("/generateVideo", server.generateVideo)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
