package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kataras/golog"

	"github.com/jinil-ha/blind-msg/utils/config"
)

var pushMessageURL string
var messageAccessToken string

func init() {
	pushMessageURL = "https://api.line.me/v2/bot/message/push"
	messageAccessToken = config.GetString("line.message.access_token")
}

type sendMessageReqType struct {
	ToUserID string        `json:"to"`
	Messages []messageType `json:"messages"`
}

type messageType struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SendMessage push LINE message to user
func SendMessage(userID string, msg string) error {
	v := &sendMessageReqType{
		ToUserID: userID,
		Messages: []messageType{{Type: "text", Text: msg}},
	}
	data, _ := json.Marshal(v)
	golog.Debugf("push line message : %s", data)

	req, err := http.NewRequest("POST", pushMessageURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", messageAccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		golog.Debugf(" push line message error : %d %s", resp.StatusCode, body)
	}

	return nil
}
