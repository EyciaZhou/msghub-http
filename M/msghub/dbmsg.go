package msghub

import (
	"database/sql"
)

type Dbmsg struct{}

var DBMsg = &Dbmsg{}

func (d *Dbmsg) GetById(id string) (_res *Msg, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetById", err.(error))
		}
	}()

	var (
		NullCoverImgId sql.NullString
		NullTopic      sql.NullString
	)

	//idS := MUtils.IdPanic(id)
	row := db.QueryRow(`
		SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic, Body
			FROM msg
			WHERE id=?
			LIMIT 1`, id)

	_res = &Msg{}
	_err = row.Scan(
		&_res.Id, &_res.SnapTime, &_res.PubTime, &_res.SourceURL, &_res.Title,
		&_res.SubTitle, &NullCoverImgId, &_res.ViewType, &_res.Frm, &_res.Tag, &NullTopic, &_res.Body,
	)

	if _err != nil {
		_res = nil
		if _err != sql.ErrNoRows {
			_err = newMsghubError("Server Error : msg.GetById", _err)
		}
	}

	_res.Topic = NullTopic.String
	_res.CoverImgId = NullCoverImgId.String

	_res.PicRefs, _err = d.GetReferredPictures(id)
	if _err != nil {
		return nil, _err
	}

	return
}

func (*Dbmsg) GetReferredPictures(id string) (_res []*PicRef, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetReferredPictures", err.(error))
		}
	}()

	rows, err := db.Query(`
		SELECT
				Ref, Description, Pixes, pid, nodenum
			FROM picref, pic_task_queue
			WHERE mid=? AND picref.pid=pic_task_queue.id;`, id)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	_res = [](*PicRef){}

	var nodeNull sql.NullInt64

	var i int
	for i = 0; rows.Next(); i++ {
		ref := &PicRef{}
		_err = rows.Scan(&ref.Ref, &ref.Description, &ref.Pixes, &ref.Pid, &nodeNull)
		ref.Node = (int)(nodeNull.Int64)

		if _err != nil {
			return nil, _err
		}
		_res = append(_res, ref)
	}
	if _err = rows.Err(); _err != nil {
		return nil, _err
	}

	return _res[:i], nil
}

func (*Dbmsg) GetRecentPageFlip(Limit int, lstti int64, lstid string) (_res []*MsgInfo, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetRecentPageFlip", err.(error))
		}
	}()

	var (
		NullCoverImgId sql.NullString
		NullTopic      sql.NullString
	)

	rows, err := db.Query(`
		SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic
			FROM msg
			WHERE ? >= SnapTime
			ORDER BY SnapTime DESC
			LIMIT ?`, lstti, Limit+1) //plus one because of lstid included, it's design to avoid the error
	// when msgs have same lastti

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	_res = make([]*MsgInfo, Limit+1)

	var i int
	for i = 0; rows.Next(); i++ {
		info := &MsgInfo{}
		_err = rows.Scan(
			&info.Id, &info.SnapTime, &info.PubTime, &info.SourceURL, &info.Title,
			&info.SubTitle, &NullCoverImgId, &info.ViewType, &info.Frm, &info.Tag, &NullTopic,
		)

		if _err != nil {
			return nil, _err
		}

		if info.Id == lstid {
			i--
			continue
		} //remove lstti

		info.Topic = NullTopic.String
		info.CoverImgId = NullCoverImgId.String

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

	var (
		NullCoverImgId sql.NullString
		NullTopic      sql.NullString
	)

	rows, err := db.Query(`
		SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic
			FROM msg
			ORDER BY SnapTime DESC
			LIMIT ?`, Limit) //? SnapTime or PubTime

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	_res = make([]*MsgInfo, Limit)

	var i int
	for i = 0; rows.Next(); i++ {
		info := &MsgInfo{}
		_err = rows.Scan(
			&info.Id, &info.SnapTime, &info.PubTime, &info.SourceURL, &info.Title,
			&info.SubTitle, &NullCoverImgId, &info.ViewType, &info.Frm, &info.Tag, &NullTopic,
		)

		if _err != nil {
			return nil, _err
		}

		info.Topic = NullTopic.String
		info.CoverImgId = NullCoverImgId.String

		_res[i] = info
	}
	if _err = rows.Err(); _err != nil {
		return nil, _err
	}

	return _res[:i], nil
}
