package handle

import (
	"bug-notify/api"
	"bug-notify/dao"
	"bug-notify/model"
	"fmt"
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
	c.AddFunc("49 14 * * *", func() {
		notProceddedNums, err1 := dao.GetStatusNumByID(NOTPROCESSEDID)
		//fmt.Println("55555555555")
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

		content := "## status_id nums \n"
		content = content + "**未处理**：" + strconv.Itoa(int(notProceddedNums)) + "\n"
		content = content + "**处理中**：" + strconv.Itoa(int(processingNums)) + "\n"
		data := model.SendMsg{
			Content: content,
			IsAtAll: true,
		}
		fmt.Println(content)
		api.SendMessage(data)
	})
	c.Start()
}
