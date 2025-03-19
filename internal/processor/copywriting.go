package processor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/models"
	"log"
	"time"
)

// 根据根据关键字及关键字来源生成对应的文案，并存储
func CreateCopyWriting(items []models.TrendingItem, dbStore *db.Queries) error {
	conf, err := config.LoadConfig("../../")
	if err != nil {
		return err
	}
	url := conf.AiUrl

	titleArgs := []string{}
	sourceArgs := []string{}
	contentArgs := []string{}
	dateArgs := []pgtype.Timestamp{}
	for _, item := range items {
		titleArgs = append(titleArgs, item.Title)
		sourceArgs = append(sourceArgs, item.Source)
		dateArgs = append(dateArgs, pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		})
		requestData := map[string]interface{}{
			"model": conf.AiModel,
			"messages": []map[string]string{
				{
					"role":    conf.Role,
					"content": conf.CopyWritingContent + item.Source + "，" + item.Title, //拼接提示词，按照app.env中预先写好的格式
				},
			},
		}

		// todo 调用api逻辑替换为SendRequest
		allResp, err := SendPostRequest(requestData, url, conf)
		if err != nil {
			//todo 日志
			return err
		}

		//jsonRequest, _ := json.Marshal(requestData)
		//payload := strings.NewReader(string(jsonRequest))
		//
		//client := &http.Client{}
		//req, err := http.NewRequest(method, url, payload)
		//if err != nil {
		//	return errors.New("failed to create copywriting request : " + err.Error())
		//}
		//
		//req.Header.Add("Authorization", conf.ApiKey)
		//req.Header.Add("Content-Type", "application/json")
		//res, err := client.Do(req)
		//if err != nil {
		//	return errors.New("failed to push copywriting request: " + err.Error())
		//}
		//
		//if res.StatusCode != 200 {
		//	return errors.New("failed to push copywriting request: HTTP code " + strconv.Itoa(res.StatusCode))
		//}
		//
		//allResp, err := ioutil.ReadAll(res.Body)
		//if err != nil {
		//	return err
		//}
		//将返回的数据写入到日志
		log.Println(item.Source + " copy is :" + string(allResp))

		var airesponse Response
		err = json.Unmarshal(allResp, &airesponse)

		if err != nil {
			return errors.New("failed to unmarshal AI response: " + err.Error())
		}

		contentArgs = append(contentArgs, airesponse.Choices[0].Message.Content)

	}

	arg := db.CreateMultipleCopyParams{
		Column1: titleArgs,
		Column2: sourceArgs,
		Column3: contentArgs,
		Column4: dateArgs,
	}
	err = dbStore.CreateMultipleCopy(context.Background(), arg)
	if err != nil {
		return errors.New("createMultipleCopy failed: " + err.Error())
	}

	return nil
}
