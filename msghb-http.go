package main

import (
	"github.com/EyciaZhou/msghub-http/C/CMsgs"
	"gopkg.in/macaron.v1"
)

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	CMsgs.RouterGroup(m)

	m.Run()
}
