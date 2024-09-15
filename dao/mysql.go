package dao

import (
	init_tool "bug-notify/init-tool"
	"bug-notify/model"
)

func GetPhoneByUserID(id []int32) (phone []string, err error) {
	err = init_tool.DB.Table("custom_values").Where("customized_id in ? and customized_type = ?", id, "Principal").Select("value").Find(&phone).Error
	return
}
func GetDingRobotByid(id int32) (dingdingRobot string, err error) {
	err = init_tool.DB.Table("custom_values").Where("customized_id = ? and customized_type = ?", id, "Project").Select("value").Find(&dingdingRobot).Error
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

func GetStatusNumByID(statusId []int, projectId int64) ([]model.TimeData, error) {
	var a []model.TimeData
	err := init_tool.DB.Table("issues").Select("status_id, priority_id, count(*) as count").Where("status_id in ? and project_id = ?", statusId, projectId).Group("priority_id").Group("status_id").Find(&a).Error
	return a, err
}

func GetWatchUserID(watchid int32, watchtype string) (userid []int32, err error) {
	err = init_tool.DB.Table("watchers").Where("watchable_id = ? and watchable_type = ?", watchid, watchtype).Select("user_id").Find(&userid).Error
	return
}

func GetAllProjectID() (ids []int64, err error) {
	err = init_tool.DB.Table("projects").Select("id").Find(&ids).Error
	return
}
func GetURLByProjectId(ids []int64) (urls []model.RobotUrl, err error) {
	err = init_tool.DB.Table("custom_values").Where("customized_id in ? and customized_type = ?", ids, "Project").Select("customized_id,value").Find(&urls).Error
	return

}
