package CMsgs

import (
	"github.com/EyciaZhou/msghub-http/C"
	"github.com/EyciaZhou/msghub-http/M/msghub"
	"gopkg.in/macaron.v1"
)

func ApiRouterGroup(m *macaron.Macaron) {
	m.Group("/msgs", func() {
		m.Get("/page/:limit/:lstid/:lstti", getMsgs)
		m.Get("/chan/:chan/page/:limit/:lstid/:lstti", getMsgs)
		m.Get("/:id", getMsg)
		m.Get("/chan", getChans)
	})
}
/*
func HtmlRouterGroup(m *macaron.Macaron) {
	m.Group("/msgs", func() {
		m.Get("/page/:limit/:lstid/:lstti", htmlGetMsgs)
		m.Get("/chan/:chan/page/:limit/:lstid/:lstti", htmlGetMsgs)
		m.Get("/:id", htmlGetMsg)
		m.Get("/chan", htmlGetChans)
	})
}
*/
func getMsgs(ctx *macaron.Context) {
	_chan, limit, lstti, lstid := ctx.Params("chan"), ctx.ParamsInt(":limit"), ctx.ParamsInt64(":lstti"), ctx.Params(":lstid")
	if limit > 20 || limit <= 0 {
		limit = 20 //default
	}
	if lstti < 0 {
		ctx.JSON(200, C.PackError(msghub.DBMsg.GetRecentFirstPage(_chan, limit, _chan=="")))
		return
	}
	ctx.JSON(200, C.PackError(msghub.DBMsg.GetRecentPageFlip(_chan, limit, lstti, lstid, _chan=="")))
}

func getMsg(ctx *macaron.Context) {
	id := ctx.Params(":id")
	ctx.JSON(200, C.PackError(msghub.DBMsg.GetById(id)))
}

func getChans(ctx *macaron.Context) {
	ctx.JSON(200, C.Pack(msghub.DBMsg.GetChanInfos()))
}
