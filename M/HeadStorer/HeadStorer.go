package HeadStorer
import (
	"time"
	"encoding/json"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha1"
	"net/url"
	"errors"
	"github.com/Sirupsen/logrus"
	"net/http"
	"strings"
	"io/ioutil"
	"github.com/EyciaZhou/msghub-http/C"
	"gopkg.in/macaron.v1"
)

type HeadStorer interface {
	MakeupUploadToken(username string) (string);
	Callback(ctx *macaron.Context);

	/*GetHead
		return "", if this user never upload head
	 */
	GetHead(username string) (string);
}


type QiniuHeadStorer struct {
	QiniuHeadStorerConfig
	headMark HeadMark
}

type QiniuHeadStorerConfig struct {
	AccessKey string
	SecretKey string
	Bucket string
	DownloadUrl string
	CallbackUrl string
}

type putPolicy struct {
	Scope string `json:"scope"`
	Deadline int64 `json:"deadline"`
	CallbackUrl string `json:"callbackUrl"`
	CallbackBody string `json:"callbackBody"`

	EndUser string `json:"endUser"`
	FsizeLimit int `json:"fsizeLimit"`
	DetectMime int `json:"detectMime"`
	MimeLimit string `json:"mimeLimit"`
}

func (p *QiniuHeadStorer)makeupPutPolicy(username string) *putPolicy {
	return &putPolicy{
		Scope:p.Bucket + ":" + username,
		Deadline:time.Now().Unix() + 3600,
		CallbackUrl:p.CallbackUrl,
		CallbackBody:`endUser=$(endUser)`,
		EndUser:username,
		FsizeLimit:1*1024*1024, //1m
		DetectMime:1,
		MimeLimit:`image/jpeg;image/png`,
	}
}

func hmac_sha1(bs []byte, key string) []byte {
	_hmac := hmac.New(sha1.New, ([]byte)(key))
	_hmac.Write(bs)
	return _hmac.Sum(nil)
}

func (p *QiniuHeadStorer) MakeupUploadToken(username string) (string){
	putPolicyStuct := p.makeupPutPolicy(username);
	bs, _ := json.Marshal(putPolicyStuct)
	encodedPutPolicy := base64.URLEncoding.EncodeToString(bs)

	encodedSign := base64.URLEncoding.EncodeToString(hmac_sha1(([]byte)(encodedPutPolicy), p.SecretKey))

	return p.AccessKey + ":" + encodedSign + ":" + encodedPutPolicy
}

func (p *QiniuHeadStorer) callbackHeaderAuthorization(Authorization string, Path string, Body string) bool {
	if (strings.Index(Authorization, "QBox ") != 0) {
		return false
	}
	auth := strings.Split(Authorization[5:], ":")
	if (len(auth) != 2 || auth[0] != p.AccessKey) {
		return false
	}
	return base64.URLEncoding.EncodeToString(hmac_sha1(([]byte)(Path + "\n" + Body), p.SecretKey)) == auth[1]
}

func (p *QiniuHeadStorer) Callback(ctx *macaron.Context) {
	reason := ""
	e := (error)(nil)

	defer func() {
		if e != nil {
			logrus.Error(reason, e.Error())
			ctx.JSON(http.StatusOK, C.Error(errors.New(reason)))
		}
	}()

	defer ctx.Req.Body().ReadCloser().Close()
	body_bs, err := ioutil.ReadAll(ctx.Req.Body().ReadCloser())
	if err != nil {
		reason, e = "callback:读取Body失败", err
		return
	}
	body := (string)(body_bs)
	if (!p.callbackHeaderAuthorization(ctx.Req.Header.Get("Authorization"), ctx.Req.URL.Path, body)) {
		reason, e = "callback:验证Authorization失败", errors.New("验证Authorization失败")
		return
	}
	vals, err := url.ParseQuery((string)(body))
	if err != nil {
		reason, e = "callback:非法Query", err
		return
	}
	usernme := vals.Get("endUser")
	if usernme == "" {
		reason, e = "callback, 上传成功,未知错误", errors.New((string)(body))
		return
	}

	reason, e = "服务端错误", p.headMark.Set(usernme)

	if e == nil {
		ctx.JSON(http.StatusOK, C.Pack(p.GetHead(usernme)))
	}

	return
}

func (p *QiniuHeadStorer) GetHead(username string) (string) {
	if p.headMark.Get(username) {
		return p.DownloadUrl + username
	}
	return ""
}

func NewQiniuHeadStorer(config *QiniuHeadStorerConfig, headMark HeadMark) *QiniuHeadStorer {
	return &QiniuHeadStorer{
		*config,
		headMark,
	}
}
