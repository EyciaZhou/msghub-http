package CUser
import (
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/M/MUser"
	"encoding/hex"
	"github.com/EyciaZhou/msghub-http/C"
	"github.com/go-macaron/session"
	"github.com/go-macaron/csrf"
	"net/http"
)

func RouterGroup(m *macaron.Macaron) {
	m.Group("/usr", func() {
		m.Get("/html/sign", html_sign_get)
		m.Get("/html/login", html_login_get)
		m.Post("/html/sign", csrf.Validate, html_sign_post)
		m.Post("/html/login", csrf.Validate, html_login_post)

		m.Post("/api/sign", api_sign)
		m.Post("/api/login", api_login)
	})
}

func api_sign(ctx *macaron.Context, f session.Store) {
	username, email := ctx.Query("username"), ctx.Query("email")
	nickname, pwd_hex := ctx.Query("nickname"), ctx.Query("pwd")

	pwd, err := hex.DecodeString(pwd_hex)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, C.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, C.PackError(MUser.DBUser.Add(username, email, nickname, pwd)))
}


func html_sign_get(ctx *macaron.Context, sess session.Store, sf session.Flash, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		C.HtmlErrorView(ctx, 200, "/", "请先登出")
		return
	}

	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "sign")
}

func html_sign_post(ctx *macaron.Context) {
	username, email := ctx.Query("username"), ctx.Query("email")
	nickname, pwd_hex := ctx.Query("nickname"), ctx.Query("pwd")

	pwd, err := hex.DecodeString(pwd_hex)

	if err != nil {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "/", "请求参数错误")
		return
	}

	_, e := MUser.DBUser.Add(username, email, nickname, pwd)
	if e != nil {
		C.HtmlErrorView(ctx, http.StatusOK, "/", e.Error())
		return
	}
	C.HtmlInfoView(ctx, http.StatusOK, "/usr/html/sign", "注册成功")
}

func api_login(ctx *macaron.Context, f session.Store) {
	uname, pwd_hex := ctx.Query("uname"), ctx.Query("pwd")
	pwd, err := hex.DecodeString(pwd_hex)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, C.Error(err))
		return
	}

	info, err := MUser.DBUser.Pwd_verify(uname, pwd)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, C.Error(err))
	}

	f.Set("api_uid", info.Id)
	ctx.JSON(http.StatusOK, C.Pack(info))
}

func html_login_get(ctx *macaron.Context, sess session.Store, sf session.Flash, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		ctx.Redirect("/")
		return
	}

	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "login")
}

func html_login_post(ctx *macaron.Context, f session.Store) {
	uname, pwd_hex := ctx.Query("uname"), ctx.Query("pwd")
	pwd, err := hex.DecodeString(pwd_hex)

	if err != nil {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "/login", "请求参数错误")
		return
	}

	info, err := MUser.DBUser.Pwd_verify(uname, pwd)
	if err != nil {
		C.HtmlErrorView(ctx, http.StatusUnauthorized, "/usr/html/login", err.Error())
		return
	}

	f.Set("uid", info.Id)
	C.HtmlInfoView(ctx, http.StatusOK, "/", "登陆成功")
}