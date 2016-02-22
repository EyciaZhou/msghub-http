package msghub

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/Sirupsen/logrus"
)

type msghubError struct {
	error
	_time string
}

func (m *msghubError) Error() string {
	return m._time + m.error.Error()
}

func newMsghubError(_time string, _err error) *msghubError {
	return &msghubError{_err, _time}
}

type Config struct {
	QueueTableName  string `default:"pic_task_queue"`
	PicRefTableName string `default:"picref"`
	MsgTableName    string `default:"msg"`
	TopicTableName  string `default:"topic"`

	DBAddress  string `default:"127.0.0.1"`
	DBPort     string `default:"3306"`
	DBName     string `default:"msghub"`
	DBUsername string `default:"root"`
	DBPassword string `default:"fmttm233"`
}

var config Config
var (
	db	*sql.DB
)

func init() {
	var err error
	log.Info("M.msghub Start Connect mysql")
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBAddress, config.DBPort, config.DBName)
	db, err = sql.Open("mysql", url)
	if err != nil {
		log.Panic("M.msghub Can't Connect DB REASON : " + err.Error())
		return
	}
	err = db.Ping()
	if err != nil {
		log.Panic("M.msghub Can't Connect DB REASON : " + err.Error())
		return
	}
	log.Info("M.msghub connected")
}
