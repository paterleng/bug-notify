package handle

import (
	"bug-notify/api"
	"bug-notify/model"
	"bug-notify/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const (
	// 未处理 处理中的事务id
	NOTPROCESSEDID = 2
	PROCESSINGID   = 3
)

var P = []int{NOTPROCESSEDID, PROCESSINGID}

func TimeingTasks() {
	c := cron.New()
	c.AddFunc("36 16  * * *", func() {
		//c.AddFunc("@every 3s", func() {
		//获取所有项目id
		//ids, err := dao.GetAllProjectID()
		//if err != nil {
		//	zap.L().Error("获取项目id失败:", zap.Error(err))
		//	return
		//}
		////根据项目id获取相对应的回调地址并加签
		//urls, err := dao.GetURLByProjectId(ids)
		//if err != nil {
		//	zap.L().Error("获取钉钉webhook失败:", zap.Error(err))
		//	return
		//}
		//urlMap := make(map[int64]string)
		//for i := 0; i < len(urls); i++ {
		//	urlMap[urls[i].CustomizedId] = urls[i].Value
		//}
		//for _, id := range ids {
		//	a, err := dao.GetStatusNumByID(P, id)
		//	if err != nil {
		//		zap.L().Error("查询数量失败:", zap.Error(err))
		//		return
		//	}
		//	maps := map[int]map[int]int{}
		//	for _, value := range a {
		//		childMap := make(map[int]int)
		//		childMap[value.StatusId] = int(value.Count)
		//		maps[value.PriorityId] = childMap
		//	}
		//	for i := 0; i < 3; i++ {
		//		for j := 0; j < 2; j++ {
		//			if _, ok := maps[i+1][j+1]; !ok {
		//				maps[i+1][j+1] = 0
		//			}
		//		}
		//	}

		// 使用markdown格式

		//content := "# %s 任务状态统计 \n" +
		//	"\n| **级别** | **未处理** | **处理中** | \n" +
		//	"\n| :--: | :--: | :--: | \n" +
		//	"\n| **重要**     | %d       | %d       |  \n" +
		//	"\n| **中等**     | %d       | %d       |  \n" +
		//	"\n| **普通**     | %d       | %d       |  \n" +
		//	"\n <font color=005EFF>[@所有人](#)\n</font>"
		//nowTime := time.Now().Format("2006-01-02")

		//content := fmt.Sprintf("# %s 任务状态统计\n"+
		//	"\n| **级别** | **未处理** | **处理中** |\n"+
		//	"\n| ---- | ---- | ---- |\n"+
		//	"\n| **重要**     | %d       | %d       |\n"+
		//	"\n| **中等**     | %d       | %d       |\n"+
		//	"\n| **普通**     | %d       | %d       |\n"+
		//	"\n <font color=005EFF>[@所有人](#)\n</font>", nowTime, maps[1][2], maps[1][3], maps[2][2], maps[2][3], maps[3][2], maps[3][3])
		//分割字符串
		//split := strings.Split(urlMap[id], "@")
		//var secret string
		//if len(split) >= 2 {
		//	secret = split[1]
		//}
		content := fmt.Sprintf(
			"\n--- \n" +
				"**所属项目：** <font color=#000000>123|</font>" +
				"**任务主题：** <font color=#000000>456</font>" +
				"**任务状态：** <font color=#000000>789</font>")
		sign := utils.DingSecret("SEC7c4636349cf3dde4ce3ba607a236e2adf2b74d32935832ed5d86d097726b790f")
		data := model.SendMsg{
			DingRobotURL: "https://oapi.dingtalk.com/robot/send?access_token=8ff6cde9a01910e897cb6461e75bd515ed9d683cc4924aad9439fda3e9689de1" + sign,
			Content:      content,
			IsAtAll:      true,
			MsgType:      "markdown",
		}
		err := api.SendMessage(data)
		if err != nil {
			zap.L().Error(" 定时任务发送消息失败:", zap.Error(err))
			return
		}
		//}
	})
	c.Start()
}
