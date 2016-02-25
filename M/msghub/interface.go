package msghub

import (
	"github.com/EyciaZhou/msghub-http/M/MUtils"
)

type MsgResult interface {
	ToMap() map[string]interface{}
	Bytes() []byte
}

type ChanInfo struct {
	Id string
	Title string
	LastModify int64
}

type Msg struct {
	MsgInfo
	Body string
	PicRefs []*PicRef `json:",omitempty`
}

type MsgInfo struct {
	Id         string
	SnapTime   int64
	PubTime    int64
	SourceURL  string
	Title      string
	SubTitle   string
	CoverImgId string `json:",omitempty"`
	ViewType   int
	Frm        string
	Tag        string
	Topic      string `json:",omitempty"`
}

type PicRef struct {
	Pid         string
	Ref         string `json:",omitempty"`
	Pixes	string `json:",omitempty"`
	Description string
	Node int
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
	return map[string]interface{}{
		"Id":         m.Id,
		"SnapTime":   m.SnapTime,
		"PubTime":    m.PubTime,
		"SourceURL":  m.SourceURL,
		"Title":      m.Title,
		"SubTitle":   m.SubTitle,
		"CoverImgId": m.CoverImgId,
		"ViewType":   m.ViewType,
		"Frm":        m.Frm,
		"Tag":        m.Tag,
		"Topic":      m.Topic,
	}
}

func (m *MsgInfo) Bytes() []byte {
	return MUtils.BytesPanic(m.ToMap())
}
