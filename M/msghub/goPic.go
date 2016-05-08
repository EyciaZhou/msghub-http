package msghub

import (
	"database/sql"
	"fmt"
	"errors"
	"strings"
)

func (d *Dbmsg) BuildPic(id string, bu string, node string) string {
	id = strings.TrimLeft(id, "0")

	switch node {
	case "FS":
		return fmt.Sprintf("https://pic%s.eycia.me/%s.%s", bu, id)
	case "QINIU":
		return fmt.Sprintf("http://7xtaud.com2.z0.glb.qiniucdn.com/"+id)
	}
	return ""
}

func (d *Dbmsg) GetPic(id string) (_res string, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = "", newMsghubError("Server Error : DBMsg.GetPic", err.(error))
		}
	}()

	row := db.QueryRow(`
		SELECT
				nodetype, nodenum, ext
			FROM pic_task_queue
			WHERE id=?
			LIMIT 1`, id);


	var nodetype sql.NullString
	var nodenum sql.NullString
	var mime sql.NullString

	_err = row.Scan(&nodetype, &nodenum, &mime)

	if _err != nil {
		if _err != sql.ErrNoRows {
			_err = newMsghubError("Server Error : DBMsg.GetPic", _err)
		} else {
			_err = errors.New("Not Found")
		}
		return
	}

	if !nodetype.Valid {
		_err = errors.New("not valid pic")
		return
	}

	_res = d.BuildPic(id, nodenum.String, nodetype.String)
	if _res == "" {
		_err = errors.New("not valid pic")
	}
	return
}
