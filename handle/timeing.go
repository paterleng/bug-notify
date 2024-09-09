package handle

import (
	"bug-notify/api"
	"bug-notify/dao"
	"bug-notify/model"
	"github.com/andeya/goutil/calendar/cron"
	"go.uber.org/zap"
	"strconv"
)

const (
	NOTPROCESSEDID = 2
	PROCESSINGID   = 3
)

func TimeingTasks() {
	c := cron.New()
	c.AddFunc("0 21 * * *", func() {
		notProceddedNums, err1 := dao.GetStatusNumByID(NOTPROCESSEDID)
		processingNums, err2 := dao.GetStatusNumByID(PROCESSINGID)
		if err1 != nil || err2 != nil {
			zap.L().Error("获取status_id为2的数量失败:", zap.Error(err1))
			zap.L().Error("获取status_id为3的数量失败:", zap.Error(err2))
			//fmt.Println("sssss \n")
			return
		}
		//var content strings.Builder
		//content.WriteString("## status_id nums \n")
		//content.WriteString("**未处理**：")

		content := "## 事务状态统计 \n"
		content = content + "**未处理**：" + strconv.Itoa(int(notProceddedNums))
		content = content + "\n \n **处理中**：" + strconv.Itoa(int(processingNums))
		content = content + "\n \n @所有人"
		data := model.SendMsg{
			Content: content,
			IsAtAll: true,
		}
		api.SendMessage(data)
	})
	c.Start()
}
