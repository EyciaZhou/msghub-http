package main

import (
	"github.com/EyciaZhou/msghub-http/C/CMsgs"
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/C/CPic"
	"github.com/EyciaZhou/msghub-http/C/CUser"
	"github.com/go-macaron/session"
	_ "github.com/go-macaron/session/redis"
	"github.com/EyciaZhou/msghub-http/C"
)

func main() {
	m := macaron.Classic()
	m.Use(session.Sessioner(session.Options{
		Provider:"redis",
		ProviderConfig:"addr=127.0.0.1:6379,password=fmttm233",
	}))
	m.Use(macaron.Renderer())

	C.RouterGroup(m)
	CMsgs.RouterGroup(m)
	CPic.RouterGroup(m)
	CUser.RouterGroup(m)

	m.Run()
}
