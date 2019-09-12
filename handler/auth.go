package handler

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"

	"github.com/jinil-ha/blind-msg/service"
	"github.com/jinil-ha/blind-msg/service/line"
)

var ckService string
var ckToken string

func init() {
	ckService = "_service"
	ckToken = "_token"
}

// GoAuthHandler redirect to LINE Login page
func GoAuthHandler(ctx iris.Context) {
	s := ctx.URLParam("service")

	var url string
	switch s {
	case service.LINE:
		url = line.AuthorizeURL()

	default:
		golog.Errorf("wrong service: %s", s)
		ctx.Redirect("/")
		return
	}

	ctx.SetCookieKV(ckService, s)
	golog.Warnf("Auth start: service(%s), url(%s)", s, url)
	ctx.Redirect(url)
}

// AuthHandler do auth process
func AuthHandler(ctx iris.Context) {
	var token string

	s := ctx.GetCookie(ckService)
	switch s {
	case service.LINE:
		err := line.Auth(ctx.URLParams(), &token)
		if err != nil {
			golog.Error(err)
		} else {
			ctx.SetCookieKV(ckToken, token)
		}

	}
	ctx.Redirect("/")
}

// LogoutHandler do logout process
func LogoutHandler(ctx iris.Context) {
	// TODO: logout

	ctx.RemoveCookie(ckService)
	ctx.RemoveCookie(ckToken)
	ctx.Redirect("/")
}
