package CUser
import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/captcha"
	"github.com/EyciaZhou/msghub-http/C"
	"net/http"
	"github.com/EyciaZhou/msghub-http/Utils"
	"github.com/EyciaZhou/msghub-http/M/MUser"
)

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

	_, e := MUser.DBUser.AddUser(username, email, pwd_sha256, username)
	if e != nil {
		C.HtmlErrorView(ctx, http.StatusBadRequest, "sign", e.Error())
		return
	}
	C.HtmlInfoView(ctx, http.StatusOK, "login", "注册成功")
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

	info, err := MUser.DBUser.VerifyPassword(uname, pwd_sha256)
	if err != nil {
		C.HtmlErrorView(ctx, http.StatusUnauthorized, "login", err.Error())
		return
	}

	f.Set("uid", info.Id)
	C.HtmlInfoView(ctx, http.StatusOK, "/", "登陆成功")
}