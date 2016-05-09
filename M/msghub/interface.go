package msghub

import ()

type MsgResult interface {
	ToMap() map[string]interface{}
	Bytes() []byte
}

type ChanInfo struct {
	Id         string
	Title      string
	LastModify int64
}

type Msg struct {
	MsgInfo
	Body    string
	PicRefs []PicRef `json:",omitempty"`
}

type MsgLine struct {
	MsgInfo
	Pics []string
}

type MsgInfo struct {
	Id             string
	SnapTime       int64
	PubTime        int64
	SourceURL      string
	Title          string
	SubTitle       string
	CoverImg       string `json:",omitempty"`
	ViewType       int
	AuthorId       string `json:",omitempty"`
	AuthorCoverImg string `json:",omitempty"`
	AuthorName     string `json:",omitempty"`
	Tag            string
	Topic          string `json:",omitempty"`
}

type PicRef struct {
	Url         string
	Ref         string `json:",omitempty"`
	Pixes       string `json:",omitempty"`
	Description string
}
