package dao

import (
	init_tool "bug-notify/init-tool"
	"bug-notify/model"
)

func GetPhoneByUserID(id []int32) (phone []string, err error) {
	err = init_tool.DB.Table("custom_values").Where("customized_id in ? and customized_type = ?", id, "Principal").Select("value").Find(&phone).Error
	return
}

func GetStatusByID(id int32) (status string, err error) {
	err = init_tool.DB.Table("issue_statuses").Where("id = ?", id).Select("name").Find(&status).Error
	return
}

func GetProject(id int32) (project string, err error) {
	err = init_tool.DB.Table("projects").Where("id = ?", id).Select("name").Find(&project).Error
	return
}
func GetUserInfoByUserID(id int32) (name model.UserName, err error) {
	err = init_tool.DB.Table("users").Select("lastname", "firstname").Where("id = ?", id).Find(&name).Error
	return
}

func GetStatusNumByID(id int) (int64, error) {
	var a int64
	err := init_tool.DB.Table("issues").Where("status_id = ?", id).Count(&a).Error
	return a, err
}

func GetWatchUserID(watchid int32, watchtype string) (userid []int32, err error) {
	err = init_tool.DB.Table("watchers").Where("watchable_id = ? and watchable_type = ?", watchid, watchtype).Select("user_id").Find(&userid).Error
	return

}
