package msghub

import (
	"database/sql"
	"strings"
	"errors"
	"github.com/EyciaZhou/msghub.go/generant"
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
		NullCoverImgId  sql.NullString
		NullTopic       sql.NullString
		NullAuthorCover sql.NullString
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
		&_res.SubTitle, &NullCoverImgId, &_res.ViewType, &_res.AuthorId, &NullAuthorCover, &_res.AuthorName,
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
	_res.CoverImg = NullCoverImgId.String
	_res.AuthorCoverImg = NullAuthorCover.String

	_res.PicRefs, _err = d.GetReferredPictures(id)
	if _err != nil {
		return nil, _err
	}

	return
}

func (*Dbmsg) GetReferredPictures(id string) (_res []PicRef, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetReferredPictures", err.(error))
		}
	}()

	rows, err := db.Query(`
		SELECT
			Ref, Description, Pixes, pid, nodenum, nodetype
			FROM (
				SELECT
					*
				FROM picref
				WHERE mid=?
			) tb
			LEFT JOIN pic_task_queue ON tb.pid=pic_task_queue.id;`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	_res = make([]PicRef, 0, 10)

	var nodenumNull sql.NullString
	var nodetypeNull sql.NullString
	var pid string

	var i int
	for i = 0; rows.Next(); i++ {
		ref := PicRef{}

		_err = rows.Scan(&ref.Ref, &ref.Description, &ref.Pixes, &pid, &nodenumNull, &nodetypeNull)

		ref.Url = DBMsg.BuildPic(pid, nodenumNull.String, nodetypeNull.String)

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

const _RAW_SQL_PAGE_UP = `
SELECT
		tb.id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType,
		AuthorId, AuthorCoverImg, AuthorName, Tag, Topic,
		GROUP_CONCAT(DISTINCT CONCAT(pic_task_queue.id,',',pic_task_queue.nodenum,",",pic_task_queue.nodetype) SEPARATOR ' ; ')
	FROM (
		SELECT
			ta.id, SnapTime, PubTime, SourceURL, Title, SubTitle,
			CONCAT(pb.id,',',pb.nodenum,",",pb.nodetype) as CoverImg,
            		ViewType, Tag, Topic, author.id as AuthorId,
            		author.name as AuthorName, CONCAT(pc.id,',',pc.nodenum,",",pc.nodetype) as AuthorCoverImg
		FROM (
			SELECT *
			FROM msg
`
const _RAW_SQL_PAGE_DOWN = `
			ORDER BY SnapTime DESC, id DESC
			LIMIT ?
        	) as ta
		LEFT JOIN author ON ta.AuthorId = author.id
		LEFT JOIN pic_task_queue pb ON pb.id = ta.CoverImg
		LEFT JOIN pic_task_queue pc ON pc.id = author.coverImg
	) as tb
	LEFT JOIN picref ON (
		SELECT count(*)
		FROM picref as b
		WHERE b.mid = tb.id AND b.pid <= picref.pid
	) <= 9 AND picref.mid = tb.id
	LEFT JOIN pic_task_queue ON pic_task_queue.id = picref.pid
	GROUP BY tb.id;`


func (*Dbmsg) GetRecentPageFlip(ChanId string, Limit int, lstti int64, lstid string, ignoreChan bool) (_res []*Msg, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = nil, newMsghubError("Server Error : DBMsg.GetRecentPageFlip", err.(error))
		}
	}()

	chansMutex.RLock()
	if _, ok := chans[ChanId]; !ignoreChan && !ok {
		chansMutex.RUnlock()
		return []*Msg{}, nil
	}
	chansMutex.RUnlock()

	var (
		NullCoverImg  sql.NullString
		NullTopic       sql.NullString
		NullAuthorCover sql.NullString
		NullNinePics sql.NullString

		rows *sql.Rows
		WHERE string
	)

	if lstti >= 0 {
		if ignoreChan {
			WHERE = "WHERE ? > SnapTime OR (SnapTime=? AND id < ?)"
			rows, _err = db.Query(_RAW_SQL_PAGE_UP + WHERE + _RAW_SQL_PAGE_DOWN, lstti, lstti, lstid, Limit)
		} else {
			WHERE = "WHERE (? > SnapTime OR (SnapTime=? AND id < ?)) AND Topic=?"
			rows, _err = db.Query(_RAW_SQL_PAGE_UP + WHERE + _RAW_SQL_PAGE_DOWN, lstti, lstti, lstid, ChanId, Limit)
		}
	} else {
		if ignoreChan {
			rows, _err = db.Query(_RAW_SQL_PAGE_UP + _RAW_SQL_PAGE_DOWN, Limit)
		} else {
			WHERE = "WHERE Topic=?"
			rows, _err = db.Query(_RAW_SQL_PAGE_UP + WHERE + _RAW_SQL_PAGE_DOWN, ChanId, Limit)
		}
	}

	if _err != nil {
		return nil, _err
	}
	defer rows.Close()

	_res = make([]*Msg, Limit)

	var i int
	for i = 0; rows.Next(); i++ {
		info := &Msg{}
		_err = rows.Scan(
			&info.Id, &info.SnapTime, &info.PubTime, &info.SourceURL,
			&info.Title, &info.SubTitle, &NullCoverImg, &info.ViewType,
			&info.AuthorId, &NullAuthorCover, &info.AuthorName,
			&info.Tag, &NullTopic, &NullNinePics,
		)

		if _err != nil {
			return nil, _err
		}

		info.Topic = NullTopic.String
		info.CoverImg = DBMsg.genPicUrl(NullCoverImg.String)
		info.AuthorCoverImg = DBMsg.genPicUrl(NullAuthorCover.String)

		if info.ViewType == generant.VIEW_TYPE_PICTURES {
			info.PicRefs = DBMsg.genNinePics(NullNinePics.String)
		}

		_res[i] = info
	}
	if _err = rows.Err(); _err != nil {
		return nil, _err
	}

	return _res[:i], nil
}

func (*Dbmsg) genPicUrl(plain string) string {
	if plain == "" {
		return ""
	}

	fields := strings.Split(plain, ",")
	if len(fields) != 3 {
		panic(errors.New("Illegal plain, when generate nine pictures"))
	}

	return DBMsg.BuildPic(fields[0], fields[1], fields[2])
}

func (*Dbmsg) genNinePics(plain string) []PicRef {
	if plain == "" {
		return []PicRef{}
	}

	parted_plain := strings.Split(plain, " ; ")

	result := make([]PicRef, len(parted_plain))

	for i, p := range parted_plain {
		result[i] = PicRef{ Url: DBMsg.genPicUrl(p) }
	}

	return result

}

func (*Dbmsg) GetChanInfos() []*ChanInfo {
	chansMutex.RLock()
	defer chansMutex.RUnlock()

	return chansArray
}
