package CUser
import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/csrf"
	"github.com/EyciaZhou/msghub-http/M/MUser"
)

func RouterGroup(m *macaron.Macaron) {
	m.Group("/usr", func() {
		m.Get("/html/sign", html_sign_get)
		m.Get("/html/login", html_login_get)
		m.Post("/html/sign", csrf.Validate, html_sign_post)
		m.Post("/html/login", csrf.Validate, html_login_post)

		m.Post("/api/sign", api_sign)
		m.Post("/api/login", api_login)

		m.Get("/api/head/token", api_head_token)
		m.Post("/api/head/callback", MUser.HeadStore.Callback)
	})
}