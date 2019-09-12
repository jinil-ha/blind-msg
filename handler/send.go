package handler

import (
	"fmt"

	"github.com/kataras/golog"
	"github.com/kataras/iris"

	"github.com/jinil-ha/blind-msg/service"
	"github.com/jinil-ha/blind-msg/service/line"
	"github.com/jinil-ha/blind-msg/service/slack"
	"github.com/jinil-ha/blind-msg/utils/bac"
)

// GetSendHandler send form html for sending msg
func GetSendHandler(ctx iris.Context) {
	b := ctx.URLParam("bac")
	if b == "" {
		ctx.WriteString("Error")
		return
	}

	msg := "こんにちは。この車のドライバーです。伝えたいことがあったらメッセージを入力して送信ボタンを押してください。もし電話番号やLINE IDを残したらこちらから連絡させていただきます。"
	ctx.ViewData("msg", msg)
	ctx.ViewData("token", "test_token")
	ctx.View("send.html")

	dmsg := fmt.Sprintf("Open send page\n bac: %s\n IP: %s", b, ctx.RemoteAddr())
	slack.SendChannel(dmsg)
}

// PostSendHandler send msg to user
func PostSendHandler(ctx iris.Context) {
	b := ctx.URLParam("bac")
	contact := ctx.PostValue("contact")
	message := ctx.PostValue("msg")

	// TODO: check BAC
	s, uid := bac.GetUserInfo(b)
	if uid == "" {
		golog.Errorf("cannot get user info: bac(%s)", b)
		ctx.WriteString("Error")
		return
	}

	// TODO: check cracking code in Message

	msg := fmt.Sprintf("From: %s\nMessage: %s", contact, message)
	switch s {
	case service.LINE:
		err := line.SendMessage(uid, msg)
		if err != nil {
			golog.Errorf("line message error: %s", err)
		}
	}

	// for debug
	dmsg := fmt.Sprintf("Message Sent\n bac: %s\n service: %s\n id: %s\n from: %s\n message: %s",
		b, s, uid, contact, message)
	slack.SendChannel(dmsg)

	ctx.View("send_ok.html")
}
