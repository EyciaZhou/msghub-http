package msghub
import (
	"database/sql"
	"github.com/EyciaZhou/msghub-http/M/MUtils"
)

type MsgResult interface {
	ToMap() map[string]interface{}
	Bytes() []byte
}

type Msg struct {
	MsgInfo
	Body string
}

type MsgInfo struct {
	Id string
	SnapTime int64
	PubTime int64
	SourceURL string
	Title string
	SubTitle string
	CoverImgId sql.NullString
	ViewType int
	Frm string
	Tag string
	Topic sql.NullString
}

func (m *Msg) ToMap() map[string]interface{} {
	mp := m.MsgInfo.ToMap()
	mp["Body"] = m.Body
	return mp
}

func (m *Msg) Bytes() []byte {
	return MUtils.BytesPanic(m.ToMap())
}

func (m *MsgInfo) ToMap() map[string]interface{} {
	return map[string]interface{} {
		"Id":m.Id,
		"SnapTime":m.SnapTime,
		"PubTime":m.PubTime,
		"SourceURL":m.SourceURL,
		"Title":m.Title,
		"SubTitle":m.SubTitle,
		"CoverImgId":MUtils.CanNullToInterface(m.CoverImgId),
		"ViewType":m.ViewType,
		"Frm":m.Frm,
		"Tag":m.Tag,
		"Topic":MUtils.CanNullToInterface(m.Topic),
	}
}

func (m *MsgInfo) Bytes() []byte {
	return MUtils.BytesPanic(m.ToMap())
}
