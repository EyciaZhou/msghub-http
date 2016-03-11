package MUser
import (
	"regexp"
	"github.com/wendal/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/Sirupsen/logrus"
	"github.com/EyciaZhou/msghub-http/M/MUtils"
	"strings"
	"database/sql"
	"bytes"
)

type ruler struct{
	Username *regexp.Regexp
	Email *regexp.Regexp
	Uid *regexp.Regexp
}
var Ruler = newRuler()

func newRuler() *ruler {
	return &ruler {
		Username:regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{5,16}$"),				//start with
		Email:regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`),
		Uid:regexp.MustCompile(`^[0-9]+$`),
	}
}

func (*ruler) Pwd_sha256(bs []byte) {
	return len(bs) == 16
}

type Dbuser struct{}
var DBUser = &Dbuser{}

func (dbuser *Dbuser) new_salt_pwd(pwd_sha256 []byte) (pwd_sha256_salted_sha256 []byte, salt []byte) {
	salt = MUtils.GenSalt()
	pwd_sha256_salted_sha256 = dbuser.salt_pwd(pwd_sha256, salt)
	return
}

func (dbuser *Dbuser) salt_pwd(pwd_sha256 []byte, salt []byte) (pwd_sha256_salted_sha256 []byte) {
	pwd_sha256_salted := append(pwd_sha256, salt...)
	pwd_sha256_salted_sha256 = MUtils.Sha256(pwd_sha256_salted)
	return
}

func (Dbuser *Dbuser) Pwd_verify(uname string, challenge []byte) (_user *User_base_info, _err error) {
	flag := false
	flag |= Ruler.Username.MatchString(uname) | Ruler.Email.MatchString(uname) | Ruler.Uid.MatchString(uname)

	if !flag {
		return "", newUserError("验证用户时错误", "用户名格式错误")
	}

	uname = strings.ToLower(uname)

	row := db.QueryRow(`
		SELECT
				id, username, email, master, pwd, salt
			FROM _user
			WHERE (username=? OR email=? OR id=?)
			LIMIT 1
	`, uname, uname, uname)
	var (
		old_salted_pwd []byte
		salt []byte
	)

	_user = &User_base_info{}

	err := row.Scan(&_user.Id, &_user.Username, &_user.Email, &_user.Master, &old_salted_pwd, &salt)

	if err == sql.ErrNoRows {
		return "", newUserError("验证用户时错误", "不存在的用户")
	} else if err != nil {
		logrus.Error("验证用户时错误", err.Error())
		return "", newUserError("验证用户时错误", "数据库错误")
	}

	challenge_salted := Dbuser.salt_pwd(challenge, salt)

	if bytes.Compare(challenge_salted, old_salted_pwd) != 0 {
		return "", newUserError("验证用户时错误", "密码错误")
	}

	_err = nil
	return
}

func (dbuser *Dbuser) Add(username string, email string, pwd_sha256 []byte) (int64, error) {
	if !Ruler.Username.MatchString(username) {
		return 0, newUserError("创建用户时错误", errors.New("用户名格式错误"))
	}
	if !Ruler.Email.MatchString(email) {
		return 0, newUserError("创建用户时错误", errors.New("邮箱格式错误"))
	}
	if !Ruler.Pwd_sha256(pwd_sha256) {
		return 0, newUserError("创建用户时错误", errors.New("密码格式错误"))
	}

	username = strings.ToLower(username)	//tolower
	email = strings.ToLower(email)

	pwd_sha256_salted_sha256, salt := dbuser.new_salt_pwd(pwd_sha256)

	result, err := db.Exec(`
		INSERT INTO
				_user (username, email, pwd, salt)
			VALUE
				(?,?,?,?)
	`, username, email, pwd_sha256_salted_sha256, salt)
	if err != nil {
		if e, ok := err.(mysql.MySQLError); ok {
			if e.Number == 2525 {
				return 0, newUserError("创建用户时错误", errors.New("检测到重复的用户信息"))
			}
		}
		logrus.Error("创建用户时错误", err.Error())
		return 0, newUserError("创建用户时错误", "数据库错误")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (dbuser *Dbuser) Change_pwd(uname string, old_pwd []byte, new_pwd []byte) error {
	userInfo, err := dbuser.Pwd_verify(uname, old_pwd)
	if err != nil {
		return err
	}

	new_salted, salt := dbuser.new_salt_pwd(new_pwd)

	result, err := db.Exec(`
		UPDATE
				_user
			SET
				pwd=?, salt=?
			WHERE
				id=?
	`, new_salted, salt, userInfo.Id)

	if err != nil {
		logrus.Error("修改密码时错误", err.Error())
		return newUserError("修改密码时错误", errors.New("服务器错误"))
	}

	row_cnt, _ := result.RowsAffected()
	if row_cnt != 1 {
		return newUserError("修改密码时错误", errors.New("修改失败"))
	}

	return nil
}

func (dbuser *Dbuser) Master(fromId string, grantTo string, level int) error {
	result, err := db.Exec(`
		UPDATE
			_user
		SET
			master=?
		WHERE EXISTS (
			SELECT * FROM
				_usr
			WHERE
				id=? AND master > ?
		) AND id=?
	`, level, fromId, level, grantTo)

	if err != nil {
		return newUserError("升级管理员时错误", err)
	}

	if cnt, _ := result.RowsAffected(); cnt != 1 {
		return newUserError("升级管理员时错误", "权限不足")
	}

	return nil
}