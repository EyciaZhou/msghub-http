package C
import (
	"gopkg.in/macaron.v1"
)

func RouterGroup(m *macaron.Macaron) {
	m.Get("/error", genMsgView("error"))
	m.Get("/info", genMsgView("info"))
}