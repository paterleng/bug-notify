package handle

import (
	"bug-notify/api"
	"bug-notify/dao"
	"bug-notify/model"
	"bug-notify/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

const (
	// 未处理 处理中的事务id
	NOTPROCESSEDID = 2
	PROCESSINGID   = 3
)

var P = []int{NOTPROCESSEDID, PROCESSINGID}

func TimeingTasks() {
	c := cron.New()
	c.AddFunc("0 21  * * *", func() {
		//c.AddFunc("@every 3s", func() {
		ids, err := dao.GetAllProjectID()
		if err != nil {
			zap.L().Error("获取项目id失败:", zap.Error(err))
			return
		}
		urls, err := dao.GetURLByProjectId(ids)
		if err != nil {
			zap.L().Error("获取钉钉webhook失败:", zap.Error(err))
			return
		}
		urlMap := make(map[int64]string)
		for i := 0; i < len(urls); i++ {
			urlMap[urls[i].CustomizedId] = urls[i].Value
		}
		for _, id := range ids {
			a, err := dao.GetStatusNumByID(P, id)
			if len(a) == 0 {
				//没有任务要处理,直接跳过
				continue
			}

			if err != nil {
				zap.L().Error("查询数量失败:", zap.Error(err))
				return
			}
			maps := make(map[int]map[int]int)
			for i := 0; i < 3; i++ {
				maps[i+1] = make(map[int]int)
			}
			for _, value := range a {
				maps[value.PriorityId][value.StatusId] = int(value.Count)
			}
			content := "# %s 任务状态统计 \n" +
				"\n| **级别** | **未处理** | **处理中** | " +
				"\n| :--: | :--: | :--: | " +
				"\n| **重要**     | %d       | %d       |  " +
				"\n| **中等**     | %d       | %d       |  " +
				"\n| **普通**     | %d       | %d       |  " +
				"\n <font color=005EFF>[@所有人](#)\n</font>"
			nowTime := time.Now().Format("2006-01-02")
			content = fmt.Sprintf(content, nowTime, maps[1][2], maps[1][3], maps[2][2], maps[2][3], maps[3][2], maps[3][3])
			//分割字符串
			split := strings.Split(urlMap[id], "@")
			var secret string
			if len(split) >= 2 {
				secret = split[1]
			}
			sign := utils.DingSecret(secret)
			data := model.SendMsg{
				DingRobotURL: split[0] + sign,
				Content:      content,
				IsAtAll:      true,
				MsgType:      "markdown",
			}
			err = api.SendMessage(data)
			if err != nil {
				zap.L().Error(" 定时任务发送消息失败:", zap.Error(err))
				return
			}
		}
	})
	c.Start()
}
