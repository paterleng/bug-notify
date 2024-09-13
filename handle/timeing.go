package handle

import (
	"bug-notify/api"
	"bug-notify/dao"
	"bug-notify/model"
	"fmt"
	"github.com/andeya/goutil/calendar/cron"
	"go.uber.org/zap"
	"time"
)

const (
	// 未处理 处理中的事务id
	NOTPROCESSEDID = 2
	PROCESSINGID   = 3

	//// bug处理等级
	//LEVEL_HIGHT  = 1
	//LEVEL_MIDDLE = 2
	//LEVEL_TAIL   = 3
)

var P = []int{NOTPROCESSEDID, PROCESSINGID}

func TimeingTasks() {
	c := cron.New()
	c.AddFunc("0 21 * * *", func() {
		a, err := dao.GetStatusNumByID(P)
		if err != nil {
			zap.L().Error("获取status_id为2的数量失败:", zap.Error(err))
			return
		}
		maps := map[int]map[int]int{}
		for _, value := range a {
			childMap := make(map[int]int)
			childMap[value.StatusId] = int(value.Count)
			maps[value.PriorityId] = childMap
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 2; j++ {
				if _, ok := maps[i+1][j+1]; !ok {
					maps[i+1][j+1] = 0
				}
			}
		}
		//使用字符串builders方式 高效切割
		//var content strings.Builder
		//content.WriteString("## status_id nums \n")
		//content.WriteString("**未处理**：")

		// 使用markdown格式
		content := "# %s 任务状态统计 \n" +
			"\n| **级别** | **未处理** | **处理中** | " +
			"\n| :--: | :--: | :--: | " +
			"\n| **重要**     | %d       | %d       |  " +
			"\n| **中等**     | %d       | %d       |  " +
			"\n| **普通**     | %d       | %d       |" +
			"\n\n <font color=005EFF>[@所有人](#)\n</font>"
		nowTime := time.Now().Format("2006-01-02")

		content = fmt.Sprintf(content, nowTime, maps[1][2], maps[1][3], maps[2][2], maps[2][3], maps[3][2], maps[3][3])

		data := model.SendMsg{
			Content: content,
			IsAtAll: true,
			MsgType: "markdown",
		}
		api.SendMessage(data)
	})
	c.Start()
}

// // 返回查出的一个切片
//
//	func selectBugLevel() (a [3][2]int64, err error) {
//		// 简单便利 处理等级和遍历关联 没有关联也可自行转换关联
//		//for i, _ := range a {
//		//	for j, _ := range a[i] {
//		//		a[i][j], err = dao.GetStatusNumByID(j+2, i+1)
//		//	}
//		//}
//		// 笨方法
//		//a[0][0], err = dao.GetStatusNumByID(NOTPROCESSEDID, LEVEL_HIGHT)
//		//a[0][1], err = dao.GetStatusNumByID(PROCESSINGID, LEVEL_HIGHT)
//		//
//		//a[1][0], err = dao.GetStatusNumByID(NOTPROCESSEDID, LEVEL_MIDDLE)
//		//a[1][1], err = dao.GetStatusNumByID(PROCESSINGID, LEVEL_MIDDLE)
//		//
//		//a[2][0], err = dao.GetStatusNumByID(NOTPROCESSEDID, LEVEL_TAIL)
//		//a[2][1], err = dao.GetStatusNumByID(PROCESSINGID, LEVEL_TAIL)
//		return a, err
//	}
