package handle

import (
	"bufio"
	"bug-notify/api"
	"bug-notify/controller"
	"bug-notify/dao"
	init_tool "bug-notify/init-tool"
	"bug-notify/model"
	"bug-notify/utils"
	"encoding/json"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"go.uber.org/zap"
	"os"
	"strconv"
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
		cfg := canal.NewDefaultConfig()
		cfg.Addr = init_tool.Conf.MySQLConfig.Host + ":" + strconv.Itoa(init_tool.Conf.MySQLConfig.Port)
		cfg.User = init_tool.Conf.MySQLConfig.User
		cfg.Password = init_tool.Conf.MySQLConfig.Password
		cfg.Dump.TableDB = init_tool.Conf.Table.TableDB
		cfg.Dump.Tables = init_tool.Conf.Table.TableName

		c, err := canal.NewCanal(cfg)

		if err != nil {
			zap.L().Fatal("shibai")
		}
		defer c.Close()

		c.SetEventHandler(&MyEventHandler{})

		masterPos, err := c.GetMasterPos()
		var pos uint32
		if e.Header != nil {
			pos = e.Header.LogPos
			fmt.Println("header", e.Header)
		}
		p := mysql.Position{
			Name: masterPos.Name,
			Pos:  pos,
		}

		/**
		发布的时候，转发给指定的人
		更新的时候，转发内容为更改了哪些内容
		获取到一些数据，然后进行入库，消息通知的操作，先通知，再入库
		*/

		action := e.Action
		olddata, newdata := GetData(e)
		//对比处理的数据差异
		switch action {
		//case "insert":
		//	InsertHandle(olddata, p)
		case "update":
			UpdateHandle(olddata, newdata, p)
		}
	}
	fmt.Println("表%s", e.Table)
	fmt.Println("数据", e.Rows)
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
	cfg.Addr = init_tool.Conf.MySQLConfig.Host + ":" + strconv.Itoa(init_tool.Conf.MySQLConfig.Port)
	cfg.User = init_tool.Conf.MySQLConfig.User
	cfg.Password = init_tool.Conf.MySQLConfig.Password
	cfg.Dump.TableDB = init_tool.Conf.Table.TableDB
	cfg.Dump.Tables = init_tool.Conf.Table.TableName

	c, err := canal.NewCanal(cfg)
	if err != nil {
		zap.L().Fatal("shibai")
	}

	c.SetEventHandler(&MyEventHandler{})

	file, err := os.ReadFile("pos.txt")
	if err != nil {
		zap.L().Error("读文件失败", zap.Error(err))

	}
	var pos model.Potion
	err = json.Unmarshal(file, &pos)
	if err != nil {
		zap.L().Error("反序列化失败:", zap.Error(err))
	}
	p := mysql.Position{
		Name: pos.Name,
		Pos:  pos.Pos,
	}
	c.RunFrom(p)
	//masterPos, err := c.GetMasterPos()
	//c.RunFrom(masterPos)

}

func Ttttt() {
	data := model.SendMsg{
		AtMobiles: []string{"17638641623", "15938479072"},
		IsAtAll:   false,
		//AtMobiles: []string{"17638641623"},
		//IsAtAll:   false,
		Content: "bug@17638641623@15938479072",
		MsgType: "actionCard",
		Url:     "http://192.168.10.6:3000/issues/16",
	}
	//panic("chucuol;")
	take, _ := dao.GetUserInfoByUserID(12)
	fmt.Println(take)
	fmt.Println(data.AtMobiles)

	err := api.SendMessage(data)
	if err != nil {
		zap.L().Error("消息发送失败:", zap.Error(err))
		return
	}
}

func InsertHandle(olddata *model.DataChanges, position mysql.Position) {
	//对比数据，看有什么变化
	//project, err := dao.GetProject(olddata.ProjectID)
	//if err != nil {
	//	zap.L().Error("获取项目失败:", zap.Error(err))
	//	return
	//}
	//phone, err := dao.GetPhoneByUserID(olddata.AssignedToID)
	//if err != nil {
	//	return
	//}
	//takeName, createName, err := GetUserName(olddata.AssignedToID, olddata.AuthorID)
	//if err != nil {
	//	return
	//}
	//data := model.SendMsg{
	//	AtMobiles: []string{phone},
	//	IsAtAll:   false,
	//	Content: fmt.Sprintf(
	//		"<center><font color=Blue size=5>温馨提醒</font></center>"+
	//			"---"+
	//			"> **所属项目：%s**"+
	//			"---"+
	//			"> **bug主题：%s**"+
	//			"---"+
	//			"> **创建人：%s**"+
	//			"---"+
	//			"> **处理人：%s**"+
	//			"---"+
	//			"\n @%s \n", project, olddata.Subject, createName, takeName, phone),
	//	MsgType: "actionCard",
	//	Url:     "http://192.168.10.6:3000/issues/" + strconv.Itoa(int(olddata.ID)),
	//}
	//file, err := os.ReadFile("pos.txt")
	//if err != nil {
	//	zap.L().Error("读文件失败", zap.Error(err))
	//}
	//var pos model.Potion
	//json.Unmarshal(file, pos)
	//err = api.SendMessage(data)
	//if err != nil {
	//	zap.L().Error("消息发送失败:", zap.Error(err))
	//	return
	//}
	//if pos.Pos != 0 {
	//	//存储文件
	//	marshal, err := json.Marshal(position)
	//	if err != nil {
	//		zap.L().Error("转换失败:", zap.Error(err))
	//		return
	//	}
	//	err = StroageFile(string(marshal))
	//	if err != nil {
	//		zap.L().Error("文件写入失败:", zap.Error(err))
	//		return
	//	}
	//}
}

func UpdateHandle(olddata, newdata *model.DataChanges, position mysql.Position) {
	//对比数据，看有什么变化
	project, err := dao.GetProject(newdata.ProjectID)
	if err != nil {
		zap.L().Error("获取项目失败:", zap.Error(err))
		return
	}
	//根据bugid获取到bug的关注者，然后根据bug不同的状态通知给不同的人
	userids, err := dao.GetWatchUserID(newdata.ID, "Issue")
	if err != nil {
		zap.L().Error("获取关注用户失败:", zap.Error(err))
		return
	}
	//发消息的时候根据bug状态通知到作者或处理者
	status, err := dao.GetStatusByID(newdata.StatusID)
	if err != nil {
		zap.L().Error("获取bug状态失败:", zap.Error(err))
		return
	}

	//发消息的时候根据bug状态通知到作者或处理者
	//3、4、5、6通知创建者
	//1、2、7 通知处理者
	if controller.CreateMap[newdata.StatusID] {
		userids = append(userids, newdata.AuthorID)
	} else if controller.ProcessorMap[newdata.StatusID] {
		userids = append(userids, newdata.AssignedToID)
	}

	phones, err := dao.GetPhoneByUserID(userids)
	if err != nil {
		zap.L().Error("获取手机号失败:", zap.Error(err))
		return
	}
	takeName, createName, err := GetUserName(newdata.AssignedToID, newdata.AuthorID)
	if err != nil {
		return
	}
	splicingString := utils.SplicingString(phones, "@")
	data := model.SendMsg{
		AtMobiles: phones,
		IsAtAll:   false,
		Content: fmt.Sprintf("<center><font color=Blue size=6>温馨提醒</font></center>\n"+
			"\n--- \n"+
			"\n> **所属项目：%s**\n"+
			"\n--- \n"+
			"\n> **bug主题：%s**\n"+
			"\n--- \n"+
			"\n> **bug状态：%s**\n"+
			"\n--- \n"+
			"\n> **创建人：%s** \n"+
			"\n--- \n"+
			"\n> **处理人：%s** \n"+
			"\n--- \n"+
			"\n @%s \n", project, newdata.Subject, status, createName, takeName, splicingString),
		MsgType: "actionCard",
		Url:     "http://192.168.10.6:3000/issues/" + strconv.Itoa(int(newdata.ID)),
	}
	file, err := os.ReadFile("pos.txt")
	if err != nil {
		zap.L().Error("读文件失败", zap.Error(err))
	}
	var pos model.Potion
	json.Unmarshal(file, &pos)
	err = api.SendMessage(data)
	if err != nil {
		zap.L().Error("消息发送失败:", zap.Error(err))
		return
	}
	if pos.Pos != 0 {
		//存储文件
		marshal, err := json.Marshal(position)
		if err != nil {
			zap.L().Error("转换失败:", zap.Error(err))
			return
		}
		err = StroageFile(string(marshal))
		if err != nil {
			zap.L().Error("文件写入失败:", zap.Error(err))
			return
		}

	}
}

func GetData(e *canal.RowsEvent) (*model.DataChanges, *model.DataChanges) {
	oldData := new(model.DataChanges)
	oldData.ID = e.Rows[0][0].(int32)
	oldData.ProjectID = e.Rows[0][2].(int32)
	oldData.Subject = e.Rows[0][3].(string)
	oldData.StatusID = e.Rows[0][7].(int32)
	if e.Rows[0][8] != nil {
		oldData.AssignedToID = e.Rows[0][8].(int32)
	} else {
		oldData.AssignedToID = 0
	}
	oldData.AuthorID = e.Rows[0][11].(int32)
	newData := new(model.DataChanges)
	if e.Action == "update" {
		newData.ID = e.Rows[1][0].(int32)
		newData.ProjectID = e.Rows[1][2].(int32)
		newData.Subject = e.Rows[1][3].(string)
		newData.StatusID = e.Rows[1][7].(int32)
		if e.Rows[1][8] == nil {
			newData.AssignedToID = 0
		} else {
			fmt.Println("走着路", e.Rows)
			newData.AssignedToID = e.Rows[1][8].(int32)
		}
		newData.AuthorID = e.Rows[1][11].(int32)
	}
	return oldData, newData
}

func StroageFile(data string) (err error) {
	file, err := os.OpenFile("pos.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(data)
	if err != nil {
		return
	}

	err = writer.Flush()
	if err != nil {
		return
	}

	return nil
}

// 拿到发布者和接收者用户名
func GetUserName(takeUserID, createUserID int32) (takeName, createName string, err error) {
	take, err := dao.GetUserInfoByUserID(takeUserID)
	if err != nil {
		return "", "", err
	}
	takeName = take.Lastname + take.Firstname
	create, err := dao.GetUserInfoByUserID(createUserID)
	if err != nil {
		return takeName, "", err
	}
	createName = create.Lastname + create.Firstname
	return
}
