package MUser
import (
	"regexp"
	"github.com/wendal/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/Sirupsen/logrus"
	"github.com/EyciaZhou/msghub-http/M/MUtils"
	"strings"
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

func (dbuser *Dbuser) salt_pwd(pwd_sha256 []byte) (pwd_sha256_salted_sha256 []byte, salt []byte) {
	salt = MUtils.GenSalt()

	pwd_sha256_salted := append(pwd_sha256, salt...)
	pwd_sha256_salted_sha256 = MUtils.Sha256(pwd_sha256_salted)

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

	pwd_sha256_salted_sha256, salt := dbuser.salt_pwd(pwd_sha256)

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

func (dbuser *Dbuser) ChangePwd(uname string, old_pwd []byte, new_pwd []byte) error {
	flag := false
	flag |= Ruler.Username.MatchString(uname) | Ruler.Email.MatchString(uname) | Ruler.Uid.MatchString(uname)

	if !flag {
		return newUserError("修改密码", "用户名格式错误")
	}

	uname = strings.ToLower(uname)

	row := db.QueryRow(`
		SELECT
				id, pwd
			FROM _user
			WHERE (username=? OR email=? OR id=?)
			LIMIT 1
	`, uname, uname, uname)
	var (
		id string
		old_salted_pwd []byte
	)

	row.Scan(&id, &old_salted_pwd)
}