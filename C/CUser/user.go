package CUser
import (
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/M/MUser"
	"encoding/hex"
	"github.com/EyciaZhou/msghub-http/C"
	"github.com/go-macaron/session"
	"net/http"
	"github.com/go-macaron/captcha"
	"github.com/EyciaZhou/msghub-http/Utils"
	"github.com/go-macaron/csrf"
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
	pwd_hex := ctx.Query("password")

	pwd, err := hex.DecodeString(pwd_hex)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, C.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, C.PackError(MUser.DBUser.Add(username, email, pwd)))
}


func html_sign_get(ctx *macaron.Context, sess session.Store, sf *session.Flash, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		C.HtmlErrorView(ctx, 200, "/", "请先登出")
		return
	}
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "sign")
}

func html_sign_post(ctx *macaron.Context, cpt *captcha.Captcha) {
	if !cpt.VerifyReq(ctx.Req) {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "sign", "验证码错误")
		return
	}

	username, email := ctx.Query("username"), ctx.Query("email")
	pwd, retype := ctx.Query("password"), ctx.Query("retype")

	if retype != pwd {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "sign", "两次输入密码不匹配")
	}

	if pwd == "" {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "sign", "密码不能为空")
		return
	}

	pwd_sha256 := Utils.Sha256(([]byte)(pwd))

	_, e := MUser.DBUser.Add(username, email, pwd_sha256)
	if e != nil {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "sign", e.Error())
		return
	}
	C.HtmlInfoView(ctx, http.StatusOK, "login", "注册成功")
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

func html_login_get(ctx *macaron.Context, sess session.Store, sf *session.Flash, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		ctx.Redirect("/")
		return
	}
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "login")
}

func html_login_post(ctx *macaron.Context, f session.Store, cpt *captcha.Captcha) {
	if !cpt.VerifyReq(ctx.Req) {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "login", "验证码错误")
		return
	}

	uname, pwd := ctx.Query("uname"), ctx.Query("password")

	if pwd == "" {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "sign", "密码不能为空")
		return
	}

	pwd_sha256 := Utils.Sha256(([]byte)(pwd))



	info, err := MUser.DBUser.Pwd_verify(uname, pwd_sha256)
	if err != nil {
		C.HtmlErrorView(ctx, http.StatusUnauthorized, "login", err.Error())
		return
	}

	f.Set("uid", info.Id)
	C.HtmlInfoView(ctx, http.StatusOK, "/", "登陆成功")
}