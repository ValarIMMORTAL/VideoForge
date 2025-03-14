package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 根据根据关键字及关键字来源生成对应的文案，并存储
func (p *Processor) CreateCopyWriting(item models.TrendingItem) error {
	conf, err := config.LoadConfig("../../")
	if err != nil {
		return err
	}
	url := conf.AiUrl
	method := "POST"

	requestData := map[string]interface{}{
		"model": conf.AiModel,
		"messages": []map[string]string{
			{
				"role":    conf.Role,
				"content": conf.CopyWritingContent + item.Source + "，" + item.Title, //拼接提示词，按照app.env中预先写好的格式
			},
		},
	}

	jsonRequest, _ := json.Marshal(requestData)
	payload := strings.NewReader(string(jsonRequest))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return errors.New("failed to create copywriting request : " + err.Error())
	}

	req.Header.Add("Authorization", conf.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	//if err != nil {
	//	return errors.New("failed to push copywriting request: " + err.Error())
	//}
		return errors.New("failed to push copywriting request: " + err.Error())
	}

	if res.StatusCode != 200 {
		return errors.New("failed to push copywriting request: HTTP code " + strconv.Itoa(res.StatusCode))
	}

	allResp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//将返回的数据写入到日志
	log.Println(item.Source + " copy is :" + string(allResp))

	var airesponse Response
	err = json.Unmarshal(allResp, &airesponse)

	if err != nil {
		return errors.New("failed to unmarshal AI response: " + err.Error())
	}

	//将生成的数据存储到数据库
	content := airesponse.Choices[0].Message.Content

	// 创建时间戳
	var date pgtype.Timestamp
	date.Time = time.Now()
	date.Valid = true

	// 构造数据库参数
	arg := db.CreateCopyParams{
		Title:   item.Title,
		Source:  item.Source,
		Content: content,
		Date:    date,
	}
	str, _ := json.Marshal(arg)
	fmt.Println(string(str))
	// 存储到数据库
	createRes, err := p.store.CreateCopy(context.Background(), arg)
	fmt.Println("---------------")
	if err != nil {
		fmt.Println("错误 : " + err.Error())
		return errors.New("failed to store copywriting in database: " + err.Error())
	}
	fmt.Println("------------")
	fmt.Println(createRes.ID)
	return nil
}
