package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/token"
	"net/http"
	"os"
	"path/filepath"
)

func (server *Server) UploadVideo(c *gin.Context) {
	//获取需要发送的平台名称
	platformName := c.PostForm("platform")
	title := c.PostForm("title")
	description := c.PostForm("description")
	if platformName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未指定平台"})
		return
	}

	publisher, err := server.publisherFactory.CreatePublisher(platformName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 处理文件上传
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败: " + err.Error()})
		return
	}

	// 保存到临时文件
	tempDir := "tmp"
	if err := os.MkdirAll(tempDir, 0777); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时目录失败"})
		return
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}
	fmt.Println("文件路径 : " + tempFilePath)
	defer os.Remove(tempFilePath)

	// 4. 调用 Publisher 上传
	payload, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	// 类型断言，将 interface{} 转换为具体类型
	authPayload, ok := payload.(*token.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid authorization payload"})
		return
	}

	userID := authPayload.UserId
	//var userID int32 = 1
	videoID, err := publisher.UploadVideo(c.Request.Context(), tempFilePath, title, description, "keyword", userID, server.store)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "视频上传成功",
		"video_id": videoID,
		"platform": publisher.Platform(),
	})
}
