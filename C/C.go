package C
import (
	"gopkg.in/macaron.v1"
	"net/url"
	"html"
)

type JSON struct {
	Err    int         `json:"err"`
	Data   interface{} `json:"data"`
	Reason string      `json:"reason"`
}

type canToMap interface {
	ToMap() map[string]interface{}
}

func Pack(v interface{}) *JSON {
	/*
	if cmp, ok := v.(canToMap); ok {
		return &JSON{
			Err:    0,
			Data:   cmp.ToMap(),
			Reason: "",
		}
	}
	*/
	return &JSON{
		Err:    0,
		Data:   v,
		Reason: "",
	}
}

func genMsgView(template_name string) func(ctx *macaron.Context) {
	return func(ctx *macaron.Context) {
		to, msg := ctx.Query("to"), ctx.Query("msg")

		_url, err := url.Parse(to)
		if err != nil {
			to = "/"
		} else {
			_url_final := ctx.Req.URL.ResolveReference(_url)
			if _url_final.Host != ctx.Req.URL.Host {
				to = "/"
			} else {
				to = _url_final.String()
			}
		}

		ctx.Data["error_redirect_to"] = to
		ctx.Data["error_msg"] = html.EscapeString(msg) //xss
		ctx.HTML(200, template_name)
	}
}

func RouterGroup(m *macaron.Macaron) {
	m.Get("/error", genMsgView("error"))
	m.Get("/info", genMsgView("info"))
}

func HtmlErrorView(ctx *macaron.Context, status int, to string, msg string) {
	v := url.Values{}
	_to, err := url.Parse(to)
	if err != nil {
		_to, _ = url.Parse("/")
	}
	v.Set("to", ctx.Req.URL.ResolveReference(_to).String())
	v.Set("msg", msg)
	ctx.Redirect("/error?" + v.Encode())
}

func HtmlInfoView(ctx *macaron.Context, status int, to string, msg string) {
	v := url.Values{}
	_to, err := url.Parse(to)
	if err != nil {
		_to, _ = url.Parse("/")
	}
	v.Set("to", ctx.Req.URL.ResolveReference(_to).String())
	v.Set("msg", msg)
	ctx.Redirect("/info?" + v.Encode())
}

func PackError(v interface{}, e error) *JSON {
	if e != nil {
		return &JSON{
			Err:    1,
			Data:   nil,
			Reason: e.Error(),
		}
	}
	return Pack(v)
}

func Error(e error) *JSON {
	return &JSON{
		Err:    1,
		Data:   nil,
		Reason: e.Error(),
	}
}