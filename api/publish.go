package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	userId, _ := c.Get(authorizationPayloadKey)
	videoID, err := publisher.UploadVideo(c.Request.Context(), tempFilePath, title, description, "keyword", userId)
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

//func ()  {
//
//}
