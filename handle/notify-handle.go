package handle

import (
	"bug-notify/api"
	"bug-notify/dao"
	init_tool "bug-notify/init-tool"
	"bug-notify/model"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"go.uber.org/zap"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	//根据表来判断是否处理这条数据
	//监控到变化后进行消息通知
	//通过header来判断是否为空，来判断是否是新数据

	/*
		创建个map存储要监控的表
	*/
	tableMap := make(map[string]int)
	for _, t := range init_tool.Conf.Table.TableName {
		tableMap[t]++
	}

	if _, ok := tableMap[e.Table.Name]; ok {

		/**
		发布的时候，转发给指定的人
		更新的时候，转发内容为更改了哪些内容
		获取到一些数据，然后进行入库，消息通知的操作，先通知，再入库
		*/
		fmt.Println("表%s", e.Table)
		fmt.Println("数据", e.Rows)
		//pos := e.Header.LogPos
		//action := e.Action
		//olddata, newdata := GetData(e)
		////对比处理的数据差异
		//switch action {
		//case "insert":
		//	InsertHandle(olddata, newdata)
		//case "update":
		//	UpdateHandle(olddata, newdata)
		//}
	}
	if e.Header != nil {
		fmt.Println("header", e.Header)
	}

	fmt.Println("我是action", e.Action)
	//fmt.Println("我是row", e.Rows[0][2])
	s := e.String()
	fmt.Println("我是s", s)

	//log.Infof("%s %v\n", e.Action, e.Rows)
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func NotifyHandle() {
	//在项目启动的时候记录指针的位置，用于下次启动时使用
	cfg := canal.NewDefaultConfig()
	cfg.Addr = init_tool.Conf.MySQLConfig.Host + ":" + init_tool.Conf.MySQLConfig.Port
	cfg.User = init_tool.Conf.MySQLConfig.User
	cfg.Password = init_tool.Conf.MySQLConfig.Password
	cfg.Dump.TableDB = init_tool.Conf.Table.TableDB
	cfg.Dump.Tables = init_tool.Conf.Table.TableName

	c, err := canal.NewCanal(cfg)
	if err != nil {
		zap.L().Fatal("shibai")
	}

	c.SetEventHandler(&MyEventHandler{})

	masterPos, err := c.GetMasterPos()

	c.RunFrom(masterPos)
	// Start canal
	//c.Run()
}

func Ttttt() {
	data := model.SendMsg{
		IsAtAll: true,
		Content: "test",
	}
	api.SendMessage(data)
}

func InsertHandle(olddata, newdata *model.DataChanges) {
	//对比数据，看有什么变化
	content := "bug \n"
	_, err := dao.GetProject(newdata.ProjectID)
	if err != nil {
		zap.L().Error("获取项目失败:", zap.Error(err))
		return
	}
	content = content + "所属项目：" + newdata.Subject + "\n"
	content = content + "bug主题：" + newdata.Subject + "\n"
	content = content + "bug描述：" + newdata.Description + "\n"
	if olddata.StatusID != newdata.StatusID {
		//	查看当前状态
		status, err := dao.GetStatusByID(newdata.StatusID)
		if err != nil {
			zap.L().Error("获取状态失败:", zap.Error(err))
			return
		}
		content = content + "状态：" + status + "\n"
	}

	phone, err := dao.GetPhoneByUserID(newdata.AssignedToID)
	if err != nil {
		return
	}
	userid, err := api.GetUserIDByPhone(phone)
	data := model.SendMsg{
		AtUserID: userid,
		IsAtAll:  true,
		Content:  "test",
	}
	api.SendMessage(data)
}

func UpdateHandle(olddata, newdata *model.DataChanges) {

}

func DataHandle() {

}

func GetData(e *canal.RowsEvent) (*model.DataChanges, *model.DataChanges) {
	oldData := new(model.DataChanges)
	oldData.ProjectID = e.Rows[0][2].(int)
	oldData.Subject = e.Rows[0][3].(string)
	oldData.Description = e.Rows[0][4].(string)
	oldData.StatusID = e.Rows[0][7].(int)
	oldData.AssignedToID = e.Rows[0][8].(int)
	oldData.AuthorID = e.Rows[0][11].(int)
	newData := new(model.DataChanges)
	newData.ProjectID = e.Rows[1][2].(int)
	newData.Subject = e.Rows[1][3].(string)
	newData.Description = e.Rows[1][4].(string)
	newData.StatusID = e.Rows[1][7].(int)
	newData.AssignedToID = e.Rows[1][8].(int)
	newData.AuthorID = e.Rows[1][11].(int)
	return oldData, newData
}
