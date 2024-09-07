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
	id1 = 2
	id2 = 3
)

func TimeingTasks() {
	c := cron.New()
	c.AddFunc("0 21 * * *", func() {
		_, err1 := dao.GetStatusNumByID(id1)
		fmt.Println("55555555555")
		nums2, err2 := dao.GetStatusNumByID(id2)
		if err1 != nil || err2 != nil {
			zap.L().Error("获取status_id为2的数量失败:", zap.Error(err1))
			zap.L().Error("获取status_id为3的数量失败:", zap.Error(err2))
			fmt.Println("sssss \n")
			return
		}
		content := "status_id nums \n"
		//content = content + "未处理：" + strconv.Itoa(int(nums1)) + "\n"
		content = content + "处理中：" + strconv.Itoa(int(nums2)) + "\n"
		data := model.SendMsg{
			AtUserID: "",
			Content:  content,
			IsAtAll:  true,
		}
		fmt.Println(content)
		api.SendMessage(data)
	})
	c.Start()

	//a, err := dao.GetProject(2)
	//if err != nil {
	//	fmt.Println("111111")
	//} else {
	//	fmt.Println("2222222" + a)
	//}
}
