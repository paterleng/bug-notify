package dao

import (
	init_tool "bug-notify/init-tool"
	"bug-notify/model"
)

func GetPhoneByUserID(id int32) (phone string, err error) {
	err = init_tool.DB.Table("custom_values").Where("customized_id = ?", id).Select("value").Find(&phone).Error
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
	err = init_tool.DB.Table("users").Where("id = ?", id).Find(&name).Error
	return
}
