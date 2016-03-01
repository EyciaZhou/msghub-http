package CPic

import (
	"gopkg.in/macaron.v1"
	"github.com/EyciaZhou/msghub-http/M/msghub"
)

func RouterGroup(m *macaron.Macaron) {
	m.Group("/pic", func() {
		m.Get("/:id/", getId)
	})
}

func getId(ctx *macaron.Context) (int, string) {
	id := ctx.Params("id")
	url, err := msghub.DBMsg.GetPic(id)

	if err != nil {
		return 404, err.Error()
	}

	ctx.Header().Set("Location", url)
	return 302, ""
}
