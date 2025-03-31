package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/global"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	"time"
)

func getClient(ctx context.Context, scope string, userId int32, store db.Store) (*http.Client, error) {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	//todo 切换成通过scope、userId从数据库中获取token 若当前数据库中不存在对应的token则需要进行oauth2认证
	token, err := tokenFromDb(ctx, userId, scope, store)
	var tok *oauth2.Token
	if err != nil { //获取数据失败
		log.Println("获取预存储的token失败 :" + err.Error())
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		tok, err = getTokenFromWeb(config, authURL)
		if err == nil {
			//todo 将token存储在数据库中
			err = saveToken(ctx, tok, userId, scope, store)
			if err != nil {
				return nil, err
			}
		}
	} else {
		tok = &oauth2.Token{
			AccessToken:  token.AccessToken,
			TokenType:    token.TokenType,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry,
		}
	}
	return config.Client(ctx, tok), nil
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

func tokenFromDb(ctx context.Context, userId int32, scope string, store db.Store) (db.GetOauth2TokenRow, error) {
	arg := db.GetOauth2TokenParams{
		UserID:   userId,
		Provider: "Google",
		Api:      scope,
	}
	token, err := store.GetOauth2Token(ctx, arg)
	if err != nil {
		return db.GetOauth2TokenRow{}, err
	}
	// 验证token是否过期
	if token.Expiry.Add(-30 * time.Second).Before(time.Now()) {
		// 使用现有refreshToken函数刷新令牌
		if err := newAccessToken(ctx, token.RefreshToken, userId, scope, store); err != nil {
			return db.GetOauth2TokenRow{}, fmt.Errorf("刷新令牌失败: %w", err)
		}

		// 重新读取令牌
		token, err = store.GetOauth2Token(ctx, arg)
		if err != nil {
			return db.GetOauth2TokenRow{}, err
		}

		// 检查刷新后的令牌
		if token.Expiry.Before(time.Now()) {
			return db.GetOauth2TokenRow{}, errors.New("刷新后的令牌仍然过期")
		}
	}
	return token, nil
}

// 获取存储的token信息
// 同时判断accessoken是否过期，当accesstoken过期则使用
//func tokenFromFile(file string) (*oauth2.Token, error) {
//	token, err := readAndDecodeToken(file)
//	if err != nil {
//		return nil, err
//	}
//
//	// 检查令牌是否过期（带30秒缓冲）
//	if token.Expiry.Add(-30 * time.Second).Before(time.Now()) {
//		// 使用现有refreshToken函数刷新令牌
//		if err := newAccessToken(token.RefreshToken, file); err != nil {
//			return nil, fmt.Errorf("刷新令牌失败: %w", err)
//		}
//
//		// 重新读取令牌文件
//		token, err = readAndDecodeToken(file)
//		if err != nil {
//			return nil, err
//		}
//
//		// 检查刷新后的令牌
//		if token.Expiry.Before(time.Now()) {
//			return nil, errors.New("刷新后的令牌仍然过期")
//		}
//	}
//	return token, nil
//}

// 读取并解码令牌文件
func readAndDecodeToken(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("无法打开令牌文件: %w", err)
	}
	defer f.Close()

	var token oauth2.Token
	if err := json.NewDecoder(f).Decode(&token); err != nil {
		return nil, fmt.Errorf("解码令牌失败: %w", err)
	}

	// 基本验证
	if token.AccessToken == "" {
		return nil, errors.New("无效的令牌: 缺少access_token")
	}

	return &token, nil
}

func saveToken(ctx context.Context, tok *oauth2.Token, userId int32, scope string, store db.Store) error {
	insertArg := db.InsertOauth2TokenParams{
		UserID:       userId,
		Provider:     "Google",
		Api:          scope,
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	}
	_, err := store.InsertOauth2Token(ctx, insertArg)
	if err != nil {
		return err
	}
	return nil
}

func getTokenFromWeb(config *oauth2.Config, authURL string) (*oauth2.Token, error) {
	err := openURL(authURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to open authorization URL in web server: %v", err)
	} else {
		fmt.Println("Your browser has been opened to an authorization URL.",
			" This program will resume once authorization has been provided.")
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
		return nil, fmt.Errorf("Unable to retrieve token %v", err)
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

func newConf() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}
	type cred struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
		AuthURI      string   `json:"auth_uri"`
		TokenURI     string   `json:"token_uri"`
	}

	var j struct {
		Web       *cred `json:"web"`
		Installed *cred `json:"installed"`
	}
	if err = json.Unmarshal(b, &j); err != nil {
		return nil, fmt.Errorf("json Unmarshal fail: %v", err)
	}

	return &oauth2.Config{
		ClientID:     j.Web.ClientID,
		ClientSecret: j.Web.ClientSecret,
		RedirectURL:  j.Web.RedirectURIs[0],
		Scopes:       []string{"scope1", "scope2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  j.Web.AuthURI,
			TokenURL: j.Web.TokenURI,
		},
	}, nil
}

func newAccessToken(ctx context.Context, refreToken string, userId int32, scope string, store db.Store) error {
	conf, _ := newConf()
	tkr := conf.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreToken})
	tk, err := tkr.Token()
	if err != nil {
		return nil
	}
	//更新token
	arg := db.UpdateAccessTokenParams{
		AccessToken:  tk.AccessToken,
		TokenType:    tk.TokenType,
		Expiry:       tk.Expiry,
		UserID:       userId,
		Provider:     "Google",
		RefreshToken: refreToken,
		TokenType_2:  tk.TokenType,
	}
	_, err = store.UpdateAccessToken(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}
