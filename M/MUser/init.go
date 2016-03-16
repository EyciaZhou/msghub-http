package MUser

import (
	"database/sql"
	"fmt"
	"git.eycia.me/eycia/configparser"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

type userError struct {
	_time string
	_err string
}

func (m *userError) Error() string {
	return m._time + " : " + m._err
}

func newUserError(_time string, _err string) *userError {
	return &userError{_time, _err}
}

func newUserErrorByError(_time string, _err error) *userError {
	return &userError{_time, _err.Error()}
}

type Config struct {
	DBAddress  string `default:"127.0.0.1"`
	DBPort     string `default:"3306"`
	DBName     string `default:"usr"`
	DBUsername string `default:"root"`
	DBPassword string `default:"fmttm233"`
}

var config Config
var (
	db *sql.DB
)

func init() {
	configparser.AutoLoadConfig("M.user", &config)

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
