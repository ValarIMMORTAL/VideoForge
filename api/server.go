package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
)

type Server struct {
	config config.Config //读取文件配置
	store  db.Store
	//tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(conf config.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: conf,
		store:  store,
		router: gin.Default(),
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := server.router
	//todo route定义

	server.router = router
}
