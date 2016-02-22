package msghub

import (
	"github.com/EyciaZhou/msghub-http/M/MUtils"
	"database/sql"
)

type Dbmsg struct{}
var DBMsg = &Dbmsg{}

func (*Dbmsg) GetById(id interface{}) (_res *Msg, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_err = newMsghubError("Server Error : msg.GetById", err.(error))
		}
	}()

	idS := MUtils.IdPanic(id)
	row := db.QueryRow(`
		SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic, Body
			FROM msg
			WHERE id="001"
			LIMIT 1`, idS)

	_res = &Msg {}
	_err = row.Scan(
		_res.Id, _res.SnapTime, _res.PubTime, _res.SourceURL, _res.Title,
		_res.SubTitle, _res.CoverImgId, _res.ViewType, _res.Frm, _res.Tag, _res.Topic, _res.Body
	)

	if _err != nil {
		_res = nil
		if _err != sql.ErrNoRows {
			_err = newMsghubError("Server Error : msg.GetById", _err)
		}
	}
	return
}

func (*Dbmsg) GetRecent(Limit int) (_res []*MsgInfo, _err error) {
	
}