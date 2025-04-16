package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/pb"
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
		tempFilePath, err = server.qnManager.DownloadFile(server.config.TempDir, video.Title, userName, "videofore-videos", video.Subscribe, "su15t494p.hn-bkt.clouddn.com")
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
		tempFilePath, err = server.qnManager.DownloadFile(server.config.TempDir, parts[0], userName, "videofore-videos", subscribe, "su15t494p.hn-bkt.clouddn.com")
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

func (server *Server) UploadVideoGrpc(c *gin.Context) {
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

		arg := &pb.UploadVideoRequest{
			TempDir:      server.config.TempDir,
			Title:        title, //上传到平台的标题
			UserName:     userName,
			Bucket:       "videofore-videos",
			Subscribe:    video.Subscribe,
			Domain:       "su15t494p.hn-bkt.clouddn.com",
			FileName:     video.Title,
			PlatformName: platformName,
			UserId:       userID,
			Description:  description,
		}
		uploadRes, err := server.grpcClient.UploadVideo(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "视频上传成功",
			"video_id": uploadRes.VideoId,
			"platform": publisher.Platform(),
		})
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

		arg := &pb.UploadVideoRequest{
			TempDir:      server.config.TempDir,
			Title:        title, //上传到平台的标题
			UserName:     userName,
			Bucket:       "videofore-videos",
			Subscribe:    subscribe,
			Domain:       "su15t494p.hn-bkt.clouddn.com",
			FileName:     parts[0],
			PlatformName: platformName,
			UserId:       userID,
			Description:  description,
		}
		video, err := server.grpcClient.UploadVideo(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "视频上传成功",
			"video_id": video.VideoId,
			"platform": publisher.Platform(),
		})
	}
}
