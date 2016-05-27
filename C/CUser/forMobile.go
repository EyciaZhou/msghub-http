package CUser
import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"encoding/hex"
	"net/http"
	"github.com/EyciaZhou/msghub-http/C"
	"github.com/EyciaZhou/msghub-http/M/MUser"
	"qiniupkg.com/x/errors.v7"
)

/*
	three query field:
		username: including alphabet and digital, start with alphabet, length can from 5 to 16
		email: normal email
		nickname: nickname
		password: password not processed
	won't set session field
	return:
		uid
 */
func api_sign(ctx *macaron.Context, f session.Store) {
	username, email, nickname := ctx.Query("username"), ctx.Query("email"), ctx.Query("nickname")
	pwd_hex := ctx.Query("password")

	pwd, err := hex.DecodeString(pwd_hex)

	if err != nil {
		ctx.JSON(http.StatusOK, C.Error(errors.New("密码格式错误")))
		return
	}

	ctx.JSON(http.StatusOK, C.PackError(MUser.DBUser.AddUser(username, email, pwd, nickname)))
}

/*
api_login
	two query field:
		uname: can be username, uid or email
		pwd: in hex format, including password after sha256

	will set session: app_uid if login success
	return:
		user_base_info
 */
func api_login(ctx *macaron.Context, f session.Store) {
	uname, pwd_hex := ctx.Query("uname"), ctx.Query("pwd")
	pwd, err := hex.DecodeString(pwd_hex)

	if err != nil {
		ctx.JSON(http.StatusOK, C.Error(errors.New("密码格式错误")))
		return
	}

	info, err := MUser.DBUser.VerifyPassword(uname, pwd)
	if err != nil {
		ctx.JSON(http.StatusOK, C.Error(err))
		return
	}

	/*
	to protect from csrf, set uid name different from html side.
	because mobile side won't be csrf attack, so mobile side api not including
	csrf token.
	if hacker use mobile api to attack html, it will failure, because when user
	login on html side, won't set api_uid. so html side only can be attack by
	html api, but all html side api having csrf protecting.
	and user won't login using mobile api on html side.
	 */
	f.Set("api_uid", info.Id)
	ctx.JSON(http.StatusOK, C.Pack(info))
}

func api_head_token(ctx *macaron.Context, f session.Store) {
	username := f.Get("api_uid");
	if username == nil || username == "" {
		ctx.JSON(200, C.PackError(nil, errors.New("没有登陆或者登录过期")))
		return
	}
	token := MUser.HeadStore.MakeupUploadToken(username.(string))
	ctx.JSON(http.StatusOK, C.Pack(token))
}