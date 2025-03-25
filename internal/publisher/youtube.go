package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pule1234/VideoForge/global"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
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
	client := y.getClient(ctx, youtube.YoutubeUploadScope)
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

func (y *YouTubePublisher) getClient(ctx context.Context, scope string) *http.Client {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	fmt.Println("cachefile + ", cacheFile)
	tok, err := tokenFromFile(cacheFile)
	//fmt.Println(tok.AccessToken)
	//fmt.Println(tok.RefreshToken)
	//err = fmt.Errorf("test")
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		tok, err = getTokenFromWeb(config, authURL)
		fmt.Println(tok.AccessToken)
		fmt.Println(tok.RefreshToken)
		if err == nil {
			saveToken(cacheFile, tok)
		}
	}
	return config.Client(ctx, tok)
}

// tokenCacheFile生成凭证文件路径/filename。
// 返回生成的凭证路径/文件名。
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, "credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go.json")), err
}

// tokenFromFile从给定的文件路径中获取Token。
// 返回检索到的Token和遇到的任何读取错误。
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config, authURL string) (*oauth2.Token, error) {
	err := openURL(authURL)
	if err != nil {
		log.Fatalf("Unable to open authorization URL in web server: %v", err)
	} else {
		fmt.Println("Your browser has been opened to an authorization URL.",
			" This program will resume once authorization has been provided.\n")
		fmt.Println(authURL)
	}

	// Wait for the web server to get the code.
	fmt.Println("阻塞")
	//code := <-codeCh
	code := <-global.OauthCodeChan
	fmt.Println("code :" + code)
	return exchangeToken(config, code)
}

func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:4001/").Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}

func exchangeToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	ctx := context.Background()
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}
	return tok, nil
}

func startWebServer() (codeCh chan string, err error) {
	listener, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		return nil, err
	}
	codeCh = make(chan string)

	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this browser window.", code)
	}))

	fmt.Println("临时服务器启动")
	return codeCh, nil
}
