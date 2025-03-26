package publisher

import (
	"context"
	"fmt"
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
	//clientID := config.Config["client_id"].(string)
	//clientSecret := config.Config["client_secret"].(string)
	clientID := "775147383926-7d68eo5b1a08pktmspgnhgdm5c7s4ck1.apps.googleusercontent.com"
	clientSecret := "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"
	RedirectURL := "http://127.0.0.1:8801/ping"
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

func (y *YouTubePublisher) UploadVideo(ctx context.Context, filePath, title, description, keywords string) (string, error) {
	client := getClient(ctx, youtube.YoutubeUploadScope)
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

func (y *YouTubePublisher) RefrePlatformToken() error {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		return fmt.Errorf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		return fmt.Errorf("cacheFile not exit")
	}
	conf, err := newConf()
	if err != nil {
		return err
	}
	tkr := conf.TokenSource(context.Background(), &oauth2.Token{RefreshToken: tok.RefreshToken})
	tk, err := tkr.Token()
	if err != nil {
		return err
	}

	saveToken(cacheFile, tk)
	return nil
}
