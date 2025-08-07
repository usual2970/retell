package channel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/usual2970/retell/pkg/http"
)

const baseUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/"

type Feishu struct {
	botToken string
}

func NewFeishu(botToken string) *Feishu {
	return &Feishu{
		botToken: botToken,
	}
}

// {
// 	"StatusCode": 0,
// 	"StatusMessage": "success",
// 	"code": 0,
// 	"data": {},
// 	"msg": "success"
// }

type FeishuResponse struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
	Code          int    `json:"code"`
	Data          any    `json:"data"`
	Msg           string `json:"msg"`
}

func (f *Feishu) Send(ctx context.Context, title, content string) error {
	var data any
	json.Unmarshal([]byte(content), &data)
	resp, err := http.Post(ctx, baseUrl+f.botToken, data, nil)
	if err != nil {
		return errors.New("failed to send message to Feishu: " + err.Error())
	}

	feishuResp := &FeishuResponse{}
	if err := json.Unmarshal([]byte(resp.Body), feishuResp); err != nil {
		return errors.New("failed to parse Feishu response: " + err.Error())
	}

	if feishuResp.Code != 0 {
		return errors.New("Feishu response error: " + feishuResp.Msg)
	}

	return nil
}
