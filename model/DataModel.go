package model

type DataChanges struct {
	ID           int32  `json:"id"`
	ProjectID    int32  `json:"project_id"`     //项目ID
	Subject      string `json:"subject"`        //标题
	Description  string `json:"description"`    //问题描述
	StatusID     int32  `json:"status_id"`      //状态ID
	AssignedToID int32  `json:"assigned_to_id"` //处理人ID
	AuthorID     int32  `json:"author_id"`      //创建者ID
}

type DingResponseCommon struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type SendMsg struct {
	AtMobiles []string `json:"at_mobiles"`
	MsgType   string   `json:"msgtype"`
	Url       string   `json:"url"`
	Content   string   `json:"content"`
	IsAtAll   bool     `json:"is_at_all"`
}
type Potion struct {
	Name string `json:"Name"`
	Pos  uint32 `json:"Pos"`
}

type UserName struct {
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
}
