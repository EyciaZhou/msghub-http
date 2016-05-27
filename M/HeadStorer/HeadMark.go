package HeadStorer
import (
	"database/sql"
	"fmt"
)

type HeadMark interface {
	Set(username string) error
	Get(username string) bool
}

type MysqlHeadMark struct {
	MysqlHeadMarkConfig

	queryGet string
	querySet string
}

type MysqlHeadMarkConfig struct {
	Db *sql.DB
	Table string
	UsernameField string
	HeadField string
	Value string
}

func (p *MysqlHeadMark) Set(username string) error {
	_, err := p.Db.Exec(p.querySet, p.Value, username)
	return err
}

func (p *MysqlHeadMark) Get(username string) bool {
	var result sql.NullString
	row := p.Db.QueryRow(p.queryGet, username)
	err := row.Scan(&result)
	if err != nil {
		return false
	}
	return result.Valid
}

func NewMysqlHeadMark(conf *MysqlHeadMarkConfig) *MysqlHeadMark {
	return &MysqlHeadMark{
		*conf,
		fmt.Sprintf(`SELECT %s
				FROM %s
				WHERE %s=?`,
			conf.HeadField, conf.Table, conf.UsernameField),
		fmt.Sprintf(`UPDATE %s
				SET %s=?
				WHERE %s=?`,
			conf.Table, conf.HeadField, conf.UsernameField),
	}
}