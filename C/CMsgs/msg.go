package CMsgs

import (
	"github.com/EyciaZhou/msghub-http/C"
	"github.com/EyciaZhou/msghub-http/M/msghub"
	"gopkg.in/macaron.v1"
)

func RouterGroup(m *macaron.Macaron) {
	m.Group("/msgs", func() {
		m.Get("/pages/:limit/:lstid/:lstti", getMsgs)
		m.Get("/:id", getMsg)
	})
}

func getMsgs(ctx *macaron.Context) {
	limit, lstti, lstid := ctx.ParamsInt(":limit"), ctx.ParamsInt64(":lstti"), ctx.Params("lstid")
	if limit > 20 || limit <= 0 {
		limit = 20 //default
	}
	if lstti < 0 {
		ctx.JSON(200, C.PackError(msghub.DBMsg.GetRecentFirstPage(limit)))
		return
	}
	ctx.JSON(200, C.PackError(msghub.DBMsg.GetRecentPageFlip(limit, lstti, lstid)))
}

func getMsg(ctx *macaron.Context) {
	id := ctx.Params(":id")
	ctx.JSON(200, C.PackError(msghub.DBMsg.GetById(id)))
}
