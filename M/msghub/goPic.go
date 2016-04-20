package msghub

import (
	"database/sql"
	"fmt"
	"errors"
)

func (d *Dbmsg) GetPic(id string) (_res string, _err error) {
	defer func() {
		err := recover()
		if err != nil {
			_res, _err = "", newMsghubError("Server Error : DBMsg.GetPic", err.(error))
		}
	}()

	row := db.QueryRow(`
		SELECT
				nodenum, ext
			FROM pic_task_queue
			WHERE id=?
			LIMIT 1`, id);

	var res sql.NullString
	var ext sql.NullString

	_err = row.Scan(&res, &ext)

	if _err != nil {
		if _err != sql.ErrNoRows {
			_err = newMsghubError("Server Error : DBMsg.GetPic", _err)
		} else {
			_err = errors.New("Not Found")
		}
		return
	}

	if !res.Valid {
		return "", errors.New("Not Found")
	}

	return fmt.Sprintf("https://pic%s.eycia.me:8080/%s.%s", res.String, id, ext.String), nil
}
