package publisher

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"strings"
)

type YouTubePublisher struct {
	oauthConfig *oauth2.Config
}

func NewYouTubePublisher(config PlatformConfig) (Publisher, error) {
	clientID := config.Config["client_id"].(string)
	clientSecret := config.Config["client_secret"].(string)
	uris, ok := config.Config["redirect_uris"].([]interface{})
	if !ok || len(uris) == 0 {
		log.Fatal("redirect_uris is not a slice or is empty")
	}
	RedirectURL, ok := uris[0].(string)
	return &YouTubePublisher{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  RedirectURL,
			Scopes:       []string{youtube.YoutubeUploadScope},
			Endpoint:     google.Endpoint,
		},
	}, nil
}

func (y *YouTubePublisher) UploadVideo(ctx context.Context, filePath, title, description, keywords string, userId int) (string, error) {
	client, err := getClient(ctx, youtube.YoutubeUploadScope, userId)
	if err != nil {
		return "", err
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  "22",
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "private"},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}
	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", filePath, err)
	}

	response, err := call.Media(file).Do()
	if err != nil {
		return "", err
	}

	return response.Id, nil
}

func (y *YouTubePublisher) Platform() string {
	return "youtube"
}

// 获取回掉地址
func (y *YouTubePublisher) GetAuthURL() string {
	return y.oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}
