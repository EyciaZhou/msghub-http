package main

import (
	"github.com/EyciaZhou/msghub-http/C/CMsgs"
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/C/CPic"
)

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	CMsgs.RouterGroup(m)
	CPic.RouterGroup(m)

	m.Run()
}
