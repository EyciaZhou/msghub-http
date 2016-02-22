package msghub

import (
	"database/sql"
)

type Dbmsg struct{}
var DBMsg = &Dbmsg{}

func (*Dbmsg) GetById(id string) (_res *Msg, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetById", err.(error))
		}
	}()

	//idS := MUtils.IdPanic(id)
	row := db.QueryRow(`
		SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic, Body
			FROM msg
			WHERE id=?
			LIMIT 1`, id)

	_res = &Msg {}
	_err = row.Scan(
		_res.Id, _res.SnapTime, _res.PubTime, _res.SourceURL, _res.Title,
		_res.SubTitle, _res.CoverImgId, _res.ViewType, _res.Frm, _res.Tag, _res.Topic, _res.Body,
	)

	if _err != nil {
		_res = nil
		if _err != sql.ErrNoRows {
			_err = newMsghubError("Server Error : msg.GetById", _err)
		}
	}
	return
}

func (*Dbmsg) GetRecentPageFlip(Limit int, lstti int64, lstid string) (_res []*MsgInfo, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetRecentPageFlip", err.(error))
		}
	}()

	rows, err := db.Query(`
		SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic
			FROM msg
			WHERE ? <= SnapTime
			ORDER BY SnapTime DESC
			LIMIT ?`, lstti, Limit+1) 	//plus one because of lstid included, it's design to avoid the error
							// when msgs have same lastti

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	_res = make([]*MsgInfo, Limit+1)

	var i int
	for i=0; rows.Next(); i++ {
		info := &MsgInfo{}
		_err = rows.Scan(
			info.Id, info.SnapTime, info.PubTime, info.SourceURL, info.Title,
			info.SubTitle, info.CoverImgId, info.ViewType, info.Frm, info.Tag, info.Topic,
		)

		if _err != nil {
			return nil, _err
		}

		if info.Id == lstid {
			i--
			continue
		} //remove lstti

		_res[i] = info
	}
	if _err = rows.Err(); _err != nil {
		return nil, _err
	}

	return _res[:i], nil
}

func (*Dbmsg) GetRecentFirstPage(Limit int) (_res []*MsgInfo, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetRecentFirstPage", err.(error))
		}
	}()

	rows, err := db.Query(`
		SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic
			FROM msg
			ORDER BY SnapTime DESC
			LIMIT ?`, Limit)	//? SnapTime or PubTime

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	_res = make([]*MsgInfo, Limit)

	var i int
	for i=0; rows.Next(); i++ {
		info := &MsgInfo{}
		_err = rows.Scan(
			info.Id, info.SnapTime, info.PubTime, info.SourceURL, info.Title,
			info.SubTitle, info.CoverImgId, info.ViewType, info.Frm, info.Tag, info.Topic,
		)

		if _err != nil {
			return nil, _err
		}

		_res[i] = info
	}
	if _err = rows.Err(); _err != nil {
		return nil, _err
	}

	return _res[:i], nil
}