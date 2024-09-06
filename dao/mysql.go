package dao

import init_tool "bug-notify/init-tool"

func GetPhoneByUserID(id int) (phone string, err error) {
	err = init_tool.DB.Table("custom_values").Where("customized_id = ?", id).Select("value").Find(&phone).Error
	return
}

func GetStatusByID(id int) (status string, err error) {
	err = init_tool.DB.Table("issue_statuses").Where("id = ?", id).Select("name").Find(&status).Error
	return
}

func GetProject(id int) (project string, err error) {
	err = init_tool.DB.Table("projects").Where("id = ?", id).Select("name").Find(&project).Error
	return
}

func GetStatusNumByID(id int) (int64, error) {
	var a *int64
	err := init_tool.DB.Table("issues").Where("status_id = ?", id).Count(a).Error
	return *a, err
}
