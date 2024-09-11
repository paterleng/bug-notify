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
	"github.com/go-mysql-org/go-mysql/replication"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRotate(header *replication.EventHeader, rotateEvent *replication.RotateEvent) error {
	fmt.Println(header)
	fmt.Println(rotateEvent)
	return nil

}

func (h *MyEventHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	//存储文件
	marshal, err := json.Marshal(pos)
	if err != nil {
		zap.L().Error("转换失败:", zap.Error(err))
		return err
	}
	err = StroageFile(string(marshal))
	if err != nil {
		zap.L().Error("文件写入失败:", zap.Error(err))
		return err
	}
	return nil
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	action := e.Action
	newdata := GetData(e)
	switch action {
	case controller.UPDATE:
		UpdateHandle(newdata)
	}
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func NotifyHandle() {
	c, err := init_tool.GoMysqlConn()
	if err != nil {
		zap.L().Fatal("创建连接失败")
		return
	}

	c.SetEventHandler(&MyEventHandler{})

	file, err := os.ReadFile(controller.POSFILENAME)
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
}

func UpdateHandle(newdata *model.DataChanges) {
	project, err := dao.GetProject(newdata.ProjectID)
	if err != nil {
		zap.L().Error("获取项目失败:", zap.Error(err))
		return
	}

	userids, err := dao.GetWatchUserID(newdata.ID, "Issue")
	if err != nil {
		zap.L().Error("获取关注用户失败:", zap.Error(err))
		return
	}

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
		zap.L().Error("名字获取失败", zap.Error(err))
		return
	}
	if takeName == "" {
		takeName = controller.NOSPECIFIED
	}
	dingdingRobot, err := dao.GetDingRobotByid(newdata.ProjectID)
	if err != nil {
		zap.L().Error("机器人获取失败", zap.Error(err))
		return
	}
	if dingdingRobot == "" {
		zap.L().Error("机器人为空", zap.Error(err))
		return
	}
	splicingString := utils.SplicingString(phones, "@")
	var btns []model.ActionBtns
	btns = append(btns, model.ActionBtns{ActionURL: "dingtalk://dingtalkclient/page/link?url=" + init_tool.Conf.Redmine.URL + strconv.Itoa(int(newdata.ID)) + "&pc_slide=true", Title: "钉钉打开"})
	btns = append(btns, model.ActionBtns{ActionURL: "dingtalk://dingtalkclient/page/link?url=" + init_tool.Conf.Redmine.URL + strconv.Itoa(int(newdata.ID)) + "&pc_slide=false", Title: "浏览器打开"})
	data := model.SendMsg{
		DingRobotURL: dingdingRobot,
		AtMobiles:    phones,
		IsAtAll:      false,
		Content: fmt.Sprintf("### <center><font color=005EFF>温馨提醒</font></center>\n"+
			"\n--- \n"+
			"\n> **所属项目：** <font color=#161823>%s</font>\n"+
			"\n--- \n"+
			"\n> **任务主题：** <font color=#161823>%s</font>\n"+
			"\n--- \n"+
			"\n> **任务状态：** <font color=#161823>%s</font>\n"+
			"\n--- \n"+
			"\n> **创建人：** <font color=#161823>%s</font>\n"+
			"\n--- \n"+
			"\n> **处理人：** <font color=#161823>%s</font>\n"+
			"\n--- \n"+
			"\n <font color=005EFF>@%s</font> \n", project, newdata.Subject, status, createName, takeName, splicingString),
		MsgType:    "actionCard",
		Url:        init_tool.Conf.Redmine.URL + strconv.Itoa(int(newdata.ID)),
		ActionBtns: btns,
	}
	err = api.SendMessage(data)
	if err != nil {
		zap.L().Error("消息发送失败:", zap.Error(err))
		return
	}
}

func GetData(e *canal.RowsEvent) *model.DataChanges {
	newData := new(model.DataChanges)
	if e.Action == controller.UPDATE {
		newData.ID = e.Rows[1][0].(int32)
		newData.ProjectID = e.Rows[1][2].(int32)
		newData.Subject = e.Rows[1][3].(string)
		newData.StatusID = e.Rows[1][7].(int32)
		if e.Rows[1][8] == nil {
			newData.AssignedToID = 0
		} else {
			newData.AssignedToID = e.Rows[1][8].(int32)
		}
		newData.AuthorID = e.Rows[1][11].(int32)
	}
	return newData
}

func StroageFile(data string) (err error) {
	file, err := os.OpenFile(controller.POSFILENAME, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	mu := sync.Mutex{}
	mu.Lock()
	_, err = writer.WriteString(data)
	if err != nil {
		return
	}
	mu.Unlock()
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
