package handle

import (
	"bug-notify/api"
	"bug-notify/dao"
	"bug-notify/model"
	"fmt"
	"github.com/andeya/goutil/calendar/cron"
	"go.uber.org/zap"
)

const (
	// 未处理 处理中的事务id
	NOTPROCESSEDID = 2
	PROCESSINGID   = 3

	// bug处理等级
	LEVEL_HIGHT  = 1
	LEVEL_MIDDLE = 2
	LEVEL_TAIL   = 3
)

// 返回查出的一个切片
func selectBugLevel() (a [3][2]int64, err error) {
	// 简单便利 处理等级和遍历关联 没有关联也可自行转换关联
	for i, _ := range a {
		for j, _ := range a[i] {
			a[i][j], err = dao.GetStatusNumByID(j+2, i+1)
		}
	}
	// 笨方法
	//a[0][0], err = dao.GetStatusNumByID(NOTPROCESSEDID, LEVEL_HIGHT)
	//a[0][1], err = dao.GetStatusNumByID(PROCESSINGID, LEVEL_HIGHT)
	//
	//a[1][0], err = dao.GetStatusNumByID(NOTPROCESSEDID, LEVEL_MIDDLE)
	//a[1][1], err = dao.GetStatusNumByID(PROCESSINGID, LEVEL_MIDDLE)
	//
	//a[2][0], err = dao.GetStatusNumByID(NOTPROCESSEDID, LEVEL_TAIL)
	//a[2][1], err = dao.GetStatusNumByID(PROCESSINGID, LEVEL_TAIL)
	return a, err
}
func TimeingTasks() {
	c := cron.New()
	c.AddFunc("0 21 * * *", func() {
		//c.AddFunc("0 21 * * *", func() {
		a, err := selectBugLevel()
		if err != nil {
			zap.L().Error("获取status_id为2的数量失败:", zap.Error(err))
			//fmt.Println("sssss \n")
			return
		}
		//使用字符串builders方式 高效切割
		//var content strings.Builder
		//content.WriteString("## status_id nums \n")
		//content.WriteString("**未处理**：")

		// 使用markdown格式
		content := "# 学生信息表  \n" +
			"\n| **bug级别** | **未处理** | **处理中** | " +
			"\n| :--: | :--: | :--: | " +
			"\n| **重要**     | %d       | %d       |  " +
			"\n| **中等**     | %d       | %d       |  " +
			"\n| **普通**     | %d       | %d       |" +
			"\n\n <font color=005EFF>[@所有人](#)\n</font>"
		content = fmt.Sprintf(content, a[0][0], a[0][1], a[1][0], a[1][1], a[2][0], a[2][1])
		//content := "## 事务状态统计 \n"
		//content = content + "**重要未处理**：" + strconv.Itoa(int(a[0][0]))
		//content = content + "**重要处理中**：" + strconv.Itoa(int(a[0][1]))
		//content = content + "\n \n**中等未处理**：" + strconv.Itoa(int(a[1][0]))
		//content = content + "**中等处理中**：" + strconv.Itoa(int(a[1][1]))
		//content = content + "\n \n**普通未处理**：" + strconv.Itoa(int(a[2][0]))
		//content = content + "**普通处理中**：" + strconv.Itoa(int(a[2][1]))
		//content = content + "\n \n @所有人"

		data := model.SendMsg{
			Content: content,
			IsAtAll: true,
			MsgType: "markdown",
		}
		//fmt.Println(data.MsgType)
		api.SendMessage(data)
	})
	c.Start()
}
