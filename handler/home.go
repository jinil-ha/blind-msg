package handler

import (
	"fmt"

	"github.com/kataras/golog"
	"github.com/kataras/iris"

	"github.com/jinil-ha/blind-msg/service"
	svc "github.com/jinil-ha/blind-msg/service"
	"github.com/jinil-ha/blind-msg/service/line"
	"github.com/jinil-ha/blind-msg/service/slack"
	"github.com/jinil-ha/blind-msg/utils/bac"
)

// TopHandler services top page
func TopHandler(ctx iris.Context) {
	var profile service.ProfileType

	service := ctx.GetCookie(ckService)
	token := ctx.GetCookie(ckToken)
	logined := false

	var err error
	if token != "" {
		switch service {
		case svc.LINE:
			err = line.GetProfile(token, &profile)
			if err != nil {
				ctx.RemoveCookie(ckService)
				ctx.RemoveCookie(ckToken)
			} else {
				logined = true
			}
		}

		if !logined {
			golog.Infof("get profile error: %s", err)
			golog.Debugf("invalid token: %s %s", service, token)
		}
	}

	if logined {
		golog.Infof("user entered : %s", profile.UserID)
		dmsg := fmt.Sprintf("User logined\n id: %s\n name: %s\n pic: %s",
			profile.UserID, profile.DisplayName, profile.PictureURL)
		slack.SendChannel(dmsg)

		b, err := bac.GetBAC(service, profile.UserID)
		if err != nil {
			golog.Errorf("get BAC error: %s", err)
		}
		err = bac.CreateQR(b)
		if err != nil {
			golog.Errorf("create QR code error : %s", err)

		} else {
			ctx.ViewData("qr_url", bac.GetQRURL(b))
			ctx.ViewData("send_url", bac.GetSendURL(b))
		}

		ctx.ViewData("name", profile.DisplayName)
		ctx.ViewData("pic_url", profile.PictureURL)

		ctx.View("top.html")
	} else {
		ctx.View("index.html")

		dmsg := fmt.Sprintf("Open homepage\n IP: %s", ctx.RemoteAddr())
		slack.SendChannel(dmsg)
	}
}
