package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/pb"
	"net/http"
)

func (server *Server) Dycrawler(c *gin.Context) {
	var req DycrawlerRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(req.Url)
	arg := &pb.CrawlerRequest{
		Url: req.Url,
	}

	crawler, err := server.grpcClient.Crawler(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"crawler success": crawler})
}
