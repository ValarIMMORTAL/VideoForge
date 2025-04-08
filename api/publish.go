package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/token"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	userName := authPayload.Username

	// 3. 处理文件上传
	file, err := c.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败: " + err.Error()})
		return
	}
	var tempFilePath string
	if err == http.ErrMissingFile {
		videoId := c.PostForm("video_id")
		if videoId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "未指定视频"})
			return
		}
		num, err := strconv.ParseInt(videoId, 10, 64)
		video, err2 := server.store.GetVideosById(c, num)
		if err2 != nil {
			return
		}
		tempFilePath, err = server.qnManager.DownloadFile(server.config, video.Title, userName, "videofore-videos", video.Subscribe, "su15t494p.hn-bkt.clouddn.com")
		if err != nil {
			return
		}
	} else {
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
	}
	defer os.Remove(tempFilePath)

	// 4. 调用 Publisher 上传
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
