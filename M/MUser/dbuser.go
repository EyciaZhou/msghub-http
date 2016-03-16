package MUser
import (
	"regexp"
	"github.com/go-sql-driver/mysql"
	"github.com/Sirupsen/logrus"
	"github.com/EyciaZhou/msghub-http/Utils"
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

func (*ruler) Pwd_sha256(bs []byte) bool {
	return len(bs) == 16
}

func (*ruler) Nickname(name string) bool {
	return len(name) > 0 && len(name) <= 40
}

func (*ruler) Uname(uname string) bool {
	return Ruler.Username.MatchString(uname) || Ruler.Email.MatchString(uname) || Ruler.Uid.MatchString(uname)
}

type Dbuser struct{}
var DBUser = &Dbuser{}

func (dbuser *Dbuser) new_salt_pwd(pwd_sha256 []byte) (pwd_sha256_salted_sha256 []byte, salt []byte) {
	salt = Utils.GenSalt()
	pwd_sha256_salted_sha256 = dbuser.salt_pwd(pwd_sha256, salt)
	return
}

func (dbuser *Dbuser) salt_pwd(pwd_sha256 []byte, salt []byte) (pwd_sha256_salted_sha256 []byte) {
	pwd_sha256_salted := append(pwd_sha256, salt...)
	pwd_sha256_salted_sha256 = Utils.Sha256(pwd_sha256_salted)
	return
}

func (Dbuser *Dbuser) Pwd_verify(uname string, challenge []byte) (_user *User_base_info, _err error) {
	flag := Ruler.Uname(uname)
	if !flag {
		return nil, newUserError("验证用户时错误", "用户名格式错误")
	}

	uname = strings.ToLower(uname)

	row := db.QueryRow(`
		SELECT
				id, username, email, master, pwd, salt, nickname
			FROM _user
			WHERE (username=? OR email=? OR id=?)
			LIMIT 1
	`, uname, uname, uname)
	var (
		old_salted_pwd []byte
		salt []byte
	)

	_user = &User_base_info{}

	err := row.Scan(&_user.Id, &_user.Username, &_user.Email, &_user.Master, &old_salted_pwd, &salt, &_user.Nickname)

	if err == sql.ErrNoRows {
		return nil, newUserError("验证用户时错误", "不存在的用户")
	} else if err != nil {
		logrus.Error("验证用户时错误", err.Error())
		return nil, newUserError("验证用户时错误", "数据库错误")
	}

	challenge_salted := Dbuser.salt_pwd(challenge, salt)

	if bytes.Compare(challenge_salted, old_salted_pwd) != 0 {
		return nil, newUserError("验证用户时错误", "密码错误")
	}

	_err = nil
	return
}

func (Dbuser *Dbuser) Change_nickname(uname string, nickname string) (error) {
	if !Ruler.Uname(uname) {
		return newUserError("验证用户时错误", "用户名格式错误")
	}

	uname = strings.ToLower(uname)

	result, err := db.Exec(`
		UPDATE
			_user
		SET
			username=?
		WHERE (username=? OR email=? OR id=?)
		LIMIT 1
	`, uname, uname, uname)

	if err != nil {
		logrus.Error("修改昵称时错误", err.Error())
		return newUserErrorByError("修改昵称时错误", err)
	}

	cow_cnt, _ := result.RowsAffected()
	if cow_cnt != 1 {
		return newUserError("修改昵称时错误", "修改失败")
	}

	return nil
}

func (dbuser *Dbuser) Add(username string, email string, nickname string, pwd_sha256 []byte) (int64, error) {
	if !Ruler.Username.MatchString(username) {
		return 0, newUserError("创建用户时错误", "用户名格式错误")
	}
	if !Ruler.Email.MatchString(email) {
		return 0, newUserError("创建用户时错误", "邮箱格式错误")
	}
	if !Ruler.Pwd_sha256(pwd_sha256) {
		return 0, newUserError("创建用户时错误", "密码格式错误")
	}
	if !Ruler.Nickname(nickname) {
		return 0, newUserError("创建用户时错误", "昵称格式错误")
	}

	username = strings.ToLower(username)
	email = strings.ToLower(email)

	pwd_sha256_salted_sha256, salt := dbuser.new_salt_pwd(pwd_sha256)

	result, err := db.Exec(`
		INSERT INTO
				_user (username, email, pwd, salt, nickname)
			VALUE
				(?,?,?,?,?)
	`, username, email, pwd_sha256_salted_sha256, salt, nickname)
	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 2525 {
				return 0, newUserError("创建用户时错误", "检测到重复的用户信息")
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
		return newUserError("修改密码时错误", "服务器错误")
	}

	row_cnt, _ := result.RowsAffected()
	if row_cnt != 1 {
		return newUserError("修改密码时错误", "修改失败")
	}

	return nil
}

func (dbuser *Dbuser) GetId(uname string) (string, error) {
	if !Ruler.Uname(uname) {
		return "", newUserError("验证用户时错误", "用户名格式错误")
	}

	strings.ToLower(uname)

	row := db.QueryRow(`
		SELECT
				id
			FROM _user
			WHERE (username=? OR email=? OR id=?)
			LIMIT 1
	`, uname, uname, uname)

	var id string

	err := row.Scan(&id)

	if err != nil {
		logrus.Error("获取id时错误", err.Error())
		return "", newUserError("获取id时错误", "服务器错误")
	}

	return id, nil
}

func (dbuser *Dbuser) Master(from_uname string, from_pwd []byte, grantTo string, level int) error {
	base_user_info, err := dbuser.Pwd_verify(from_uname, from_pwd)
	if err != nil {
		return err
	}

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
	`, level, base_user_info.Id, level, grantTo)

	if err != nil {
		return newUserErrorByError("升级管理员时错误", err)
	}

	if cnt, _ := result.RowsAffected(); cnt != 1 {
		return newUserError("升级管理员时错误", "权限不足")
	}

	return nil
}