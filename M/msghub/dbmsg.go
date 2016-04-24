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
				msg.id, SnapTime, PubTime, SourceURL, Title, SubTitle, msg.CoverImg, ViewType,
                		author.id as AuthorId, author.coverImg as AuthorCoverImg, author.name as AuthorName,
                		Tag, Topic, Body
			FROM msg, author
			WHERE msg.id=? and msg.AuthorId = author.id
			LIMIT 1`, id)

	_res = &Msg{}
	_err = row.Scan(
		&_res.Id, &_res.SnapTime, &_res.PubTime, &_res.SourceURL, &_res.Title,
		&_res.SubTitle, &NullCoverImgId, &_res.ViewType, &_res.AuthorId, &_res.AuthorCoverImgId, &_res.AuthorName,
		&_res.Tag, &NullTopic, &_res.Body,
	)

	if _err != nil {
		_res = nil
		if _err != sql.ErrNoRows {
			_err = newMsghubError("Server Error : msg.GetById", _err)
		}
		return
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

func (*Dbmsg) GetRecentPageFlip(ChanId string, Limit int, lstti int64, lstid string, ignoreChan bool) (_res []*MsgInfo, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetRecentPageFlip", err.(error))
		}
	}()

	chansMutex.RLock()
	if _, ok := chans[ChanId]; !ignoreChan && !ok {
		chansMutex.RUnlock()
		return []*MsgInfo{}, nil
	}
	chansMutex.RUnlock()

	var (
		NullCoverImgId sql.NullString
		NullTopic      sql.NullString
	)

	var (
		rows *sql.Rows
	)

	if ignoreChan {
		rows, _err = db.Query(`
		SELECT
				msg.id, SnapTime, PubTime, SourceURL, Title, SubTitle, msg.CoverImg, ViewType,
				author.id as AuthorId, author.coverImg as AuthorCoverImg, author.name as AuthorName,
				Tag, Topic
			FROM msg, author
			WHERE ? >= SnapTime and AND msg.AuthorId = author.id
			ORDER BY PubTime DESC
			LIMIT ?`, lstti, Limit + 1)
	} else {
		rows, _err = db.Query(`
		SELECT
				msg.id, SnapTime, PubTime, SourceURL, Title, SubTitle, msg.CoverImg, ViewType,
				author.id as AuthorId, author.coverImg as AuthorCoverImg, author.name as AuthorName,
				Tag, Topic
			FROM msg, author
			WHERE ? >= SnapTime AND Topic=? AND msg.AuthorId = author.id
			ORDER BY PubTime DESC
			LIMIT ?`, lstti, ChanId, Limit + 1)
	}

	//plus one because of lstid included, it's design to avoid the error
	// when msgs have same lastti

	if _err != nil {
		return nil, _err
	}

	defer rows.Close()

	_res = make([]*MsgInfo, Limit+1)

	var i int
	for i = 0; rows.Next(); i++ {
		info := &MsgInfo{}
		_err = rows.Scan(
			&info.Id, &info.SnapTime, &info.PubTime, &info.SourceURL, &info.Title,
			&info.SubTitle, &NullCoverImgId, &info.ViewType, &info.AuthorId, &info.AuthorCoverImgId, &info.AuthorName,
			&info.Tag, &NullTopic,
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

func (*Dbmsg) GetRecentFirstPage(ChanId string, Limit int, ignoreChan bool) (_res []*MsgInfo, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetRecentFirstPage", err.(error))
		}
	}()

	chansMutex.RLock()

	if _, ok := chans[ChanId]; !ignoreChan && !ok {
		chansMutex.RUnlock()
		return []*MsgInfo{}, nil
	}
	chansMutex.RUnlock()

	var (
		NullCoverImgId sql.NullString
		NullTopic      sql.NullString
	)

	var (
		rows *sql.Rows
	)

	if ignoreChan {
		rows, _err = db.Query(`
		SELECT
				msg.id, SnapTime, PubTime, SourceURL, Title, SubTitle, msg.CoverImg, ViewType,
				author.id as AuthorId, author.coverImg as AuthorCoverImg, author.name as AuthorName,
				Tag, Top
			FROM msg, author
			WHERE msg.AuthorId = author.id
			ORDER BY PubTime DESC
			LIMIT ?`, Limit) //TODO: ? SnapTime or PubTime
	} else {
		rows, _err = db.Query(`
		SELECT
				msg.id, SnapTime, PubTime, SourceURL, Title, SubTitle, msg.CoverImg, ViewType,
				author.id as AuthorId, author.coverImg as AuthorCoverImg, author.name as AuthorName,
				Tag, Top
			FROM msg, author
			WHERE Topic=? AND msg.AuthorId = author.id
			ORDER BY PubTime DESC
			LIMIT ?`, ChanId, Limit) //TODO: ? SnapTime or PubTime
	}
	if _err != nil {
		return nil, _err
	}

	defer rows.Close()

	_res = make([]*MsgInfo, Limit)

	var i int
	for i = 0; rows.Next(); i++ {
		info := &MsgInfo{}
		_err = rows.Scan(
			&info.Id, &info.SnapTime, &info.PubTime, &info.SourceURL, &info.Title,
			&info.SubTitle, &NullCoverImgId, &info.ViewType, &info.AuthorId, &info.AuthorCoverImgId, &info.AuthorName,
			&info.Tag, &NullTopic,
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

func (*Dbmsg) GetChanInfos() []*ChanInfo {
	chansMutex.RLock()
	defer chansMutex.RUnlock()

	return chansArray
}