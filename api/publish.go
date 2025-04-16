package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/util"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

	userID, userName, err := util.GetUserByToken(c, authorizationPayloadKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 处理文件上传
	file, err := c.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败: " + err.Error()})
		return
	}
	var tempFilePath string
	if err == http.ErrMissingFile { // 不选择本地视频的情况
		fmt.Println("missing file")
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败: " + err.Error()})
			return
		}
	} else {
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败"})
			return
		}
		defer srcFile.Close()
		subscribe := time.Now().Unix()
		parts := strings.Split(file.Filename, ".")
		objectName := parts[0] + userName + fmt.Sprintf("%d", subscribe) + ".mp4"
		err = server.qnManager.UploadDataSource(c, srcFile, "videofore-videos", file.Filename, objectName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败" + err.Error()})
			return
		}
		tempFilePath, err = server.qnManager.DownloadFile(server.config, parts[0], userName, "videofore-videos", subscribe, "su15t494p.hn-bkt.clouddn.com")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败: " + err.Error()})
			return
		}
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

//func (server *Server) UploadVideoGrpc(c *gin.Context) {
//	//获取需要发送的平台名称
//	platformName := c.PostForm("platform")
//	title := c.PostForm("title")
//	description := c.PostForm("description")
//	if platformName == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "未指定平台"})
//		return
//	}
//
//	publisher, err := server.publisherFactory.CreatePublisher(platformName)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	userID, userName, err := util.GetUserByToken(c, authorizationPayloadKey)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	// 3. 处理文件上传
//	file, err := c.FormFile("video")
//	if err != nil && err != http.ErrMissingFile {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败: " + err.Error()})
//		return
//	}
//	var tempFilePath string
//	if err == http.ErrMissingFile { // 不选择本地视频的情况
//		videoId := c.PostForm("video_id")
//		if videoId == "" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "未指定视频"})
//			return
//		}
//		num, err := strconv.ParseInt(videoId, 10, 64)
//		video, err2 := server.store.GetVideosById(c, num)
//		if err2 != nil {
//			return
//		}
//		// 在grpc中执行文件下载操作
//		tempFilePath, err = server.qnManager.DownloadFile(server.config, video.Title, userName, "videofore-videos", video.Subscribe, "su15t494p.hn-bkt.clouddn.com")
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败: " + err.Error()})
//			return
//		}
//	} else { // 用户自己上传的文件直接保存在七牛云上
//		// 保存到临时文件
//		tempDir := "tmp"
//		if err := os.MkdirAll(tempDir, 0777); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时目录失败"})
//			return
//		}
//
//		tempFilePath = filepath.Join(tempDir, file.Filename)
//		if err = c.SaveUploadedFile(file, tempFilePath); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
//			return
//		}
//		fmt.Println("文件路径 : " + tempFilePath)
//		timestamp := time.Now().Unix()
//		parts := strings.Split(file.Filename, ".")
//		err = server.qnManager.UploadFile(c, tempFilePath, "videofore-videos", file.Filename, parts[0]+userName+fmt.Sprintf("%d", timestamp)+".mp4")
//
//	}
//	defer os.Remove(tempFilePath)
//
//	//调用grpc接口
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "视频上传成功",
//		//"video_id": videoID,
//		"platform": publisher.Platform(),
//	})
//}
