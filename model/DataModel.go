package model

type DataChanges struct {
	ProjectID    int    `json:"project_id"`     //项目ID
	Subject      string `json:"subject"`        //标题
	Description  string `json:"description"`    //问题描述
	StatusID     int    `json:"status_id"`      //状态ID
	AssignedToID int    `json:"assigned_to_id"` //处理人ID
	AuthorID     int    `json:"author_id"`      //创建者ID
}

type DingResponseCommon struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type SendMsg struct {
	AtUserID string `json:"at_user_id"`
	Content  string `json:"content"`
	IsAtAll  bool   `json:"is_at_all"`
}
