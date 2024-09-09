package model

import "gorm.io/gorm"

type DataChanges struct {
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
	AtMobiles []string `json:"atMobiles"`
	//AtUserID  string     `json:"at_user_id"`
	Content string `json:"content"`
	IsAtAll bool   `json:"is_at_all"`
}
type AtMobile struct {
	gorm.Model
	AtMobile string `json:"atMobile"`
	Name     string `json:"name"`
	AtID     uint   //AtMobile属于At，打上标签
}
type Potion struct {
	Name string `json:"Name"`
	Pos  uint32 `json:"Pos"`
}

type UserName struct {
	LastName  string `json:"lastname"`
	FirstName string `json:"firstname"`
}
