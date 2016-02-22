package main

import (
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/C/CMsgs"
)

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	CMsgs.RouterGroup(m)

	m.Run()
}
