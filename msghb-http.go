package main

import (
	"github.com/EyciaZhou/msghub-http/C/CMsgs"
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/C/CPic"
	"github.com/EyciaZhou/msghub-http/C/CUser"
	"github.com/go-macaron/session"
	_ "github.com/go-macaron/session/redis"
	"github.com/EyciaZhou/msghub-http/C"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/cache"
)

func main() {
	m := macaron.Classic()
	m.Use(session.Sessioner(session.Options{
		Provider:"redis",
		ProviderConfig:"addr=127.0.0.1:6379",
	}))
	m.Use(macaron.Renderer())
	m.Use(csrf.Csrfer())
	m.Use(cache.Cacher())
	m.Use(captcha.Captchaer())

	C.RouterGroup(m)
	CMsgs.ApiRouterGroup(m)
	CPic.RouterGroup(m)
	CUser.RouterGroup(m)

	m.Run()
}
