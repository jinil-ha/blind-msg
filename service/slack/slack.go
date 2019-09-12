package slack

import (
	"log"
	"os"

	"github.com/nlopes/slack"

	"github.com/jinil-ha/blind-msg/utils/config"
)

var api *slack.Client
var rtm *slack.RTM
var channelID string

// SendChannel send message to channel
func SendChannel(msg string) {
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, channelID))
}

// Start starts slack bot service
func Start() {
	token := config.GetString("slack.token")
	channelID = config.GetString("slack.channel_id")
	level := config.GetString("log_level")

	if level == "debug" {
		api = slack.New(token,
			slack.OptionDebug(true), slack.OptionLog(log.New(os.Stdout, "bot: ", log.Lshortfile|log.LstdFlags)))
	} else {
		api = slack.New(token)
	}

	rtm = api.NewRTM()
	go rtm.ManageConnection()
}
