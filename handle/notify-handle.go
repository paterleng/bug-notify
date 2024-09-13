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
	"strings"
	"sync"
)

type MyEventHandler struct {
	canal.DummyEventHandler
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
	fmt.Println("我是action", e.Action)
	fmt.Println("我是rows", e.Rows)
	olddata, newdata := GetData(e)
	//根据状态判断是否发送通知，如果是update行为，关注状态是否有变化，如果没有变化就不通知
	switch action {
	case controller.INSERT:
		InsertHandle(olddata)
	case controller.UPDATE:
		UpdateHandle(olddata, newdata)
	}
	return nil
}

func NotifyHandle() {
	c, err := init_tool.GoMysqlConn()
	if err != nil {
		zap.L().Fatal("创建连接失败", zap.Error(err))
		return
	}

	c.SetEventHandler(&MyEventHandler{})
	//在这个地方如果不存在pos.txt文件，就创建并写入当前的最新位置
	_, err = os.Stat(controller.POSFILENAME)
	if err != nil {
		if os.IsNotExist(err) {
			masterPos, err := c.GetMasterPos()
			if err != nil {
				zap.L().Fatal("获取最新指针位置失败", zap.Error(err))
				return
			}
			//存储文件
			marshal, err := json.Marshal(masterPos)
			if err != nil {
				zap.L().Error("转换失败:", zap.Error(err))
				return
			}
			err = StroageFile(string(marshal))
			if err != nil {
				zap.L().Error("文件写入失败:", zap.Error(err))
				return
			}
		} else {
			zap.L().Error("获取文件判断失败:", zap.Error(err))
			return
		}
	}

	file, err := os.ReadFile(controller.POSFILENAME)
	if err != nil {
		zap.L().Error("读文件失败", zap.Error(err))
		return
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
	err = c.RunFrom(p)
	if err != nil {
		zap.L().Error("启动失败", zap.Error(err))
		return
	}
}

func InsertHandle(newdata *model.DataChanges) {
	err := Handle(newdata)
	if err != nil {
		zap.L().Error("插入事件处理失败", zap.Error(err))
	}
}

func UpdateHandle(olddata, newdata *model.DataChanges) {
	if olddata.StatusID == newdata.StatusID {
		zap.L().Info("更新事件，状态没有变化，不做处理")
		return
	}
	err := Handle(newdata)
	if err != nil {
		zap.L().Error("更新处理失败", zap.Error(err))
	}
}

func GetData(e *canal.RowsEvent) (*model.DataChanges, *model.DataChanges) {
	oldData := new(model.DataChanges)
	oldData.ID = e.Rows[0][0].(int32)
	oldData.ProjectID = e.Rows[0][2].(int32)
	oldData.Subject = e.Rows[0][3].(string)
	oldData.StatusID = e.Rows[0][7].(int32)
	if e.Rows[0][8] == nil {
		oldData.AssignedToID = 0
	} else {
		oldData.AssignedToID = e.Rows[0][8].(int32)
	}
	oldData.AuthorID = e.Rows[0][11].(int32)
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
	return oldData, newData
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

func Handle(newdata *model.DataChanges) (err error) {
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
	//对机器人进行签名
	//分割字符串
	split := strings.Split(dingdingRobot, "@")
	var secret string
	if len(split) >= 2 {
		secret = split[1]
	}
	sign := utils.DingSecret(secret)
	var btns []model.ActionBtns
	btns = append(btns, model.ActionBtns{ActionURL: "dingtalk://dingtalkclient/page/link?url=" + init_tool.Conf.Redmine.URL + strconv.Itoa(int(newdata.ID)) + "&pc_slide=true", Title: "钉钉打开"})
	btns = append(btns, model.ActionBtns{ActionURL: "dingtalk://dingtalkclient/page/link?url=" + init_tool.Conf.Redmine.URL + strconv.Itoa(int(newdata.ID)) + "&pc_slide=false", Title: "浏览器打开"})
	data := model.SendMsg{
		DingRobotURL: split[0] + sign,
		AtMobiles:    phones,
		IsAtAll:      false,
		Content: fmt.Sprintf("### <center><font color=005EFF>温馨提醒</font></center>\n"+
			"\n--- \n"+
			"\n> **所属项目：** <font color=#000000>%s</font>\n"+
			"\n--- \n"+
			"\n> **任务主题：** <font color=#000000>%s</font>\n"+
			"\n--- \n"+
			"\n> **任务状态：** <font color=#000000>%s</font>\n"+
			"\n--- \n"+
			"\n> **创建人：** <font color=#000000>%s</font>\n"+
			"\n--- \n"+
			"\n> **处理人：** <font color=#000000>%s</font>\n"+
			"\n--- \n", project, newdata.Subject, status, createName, takeName),
		MsgType:    "actionCard",
		Url:        init_tool.Conf.Redmine.URL + strconv.Itoa(int(newdata.ID)),
		ActionBtns: btns,
	}
	if len(phones) > 0 {
		data.Content = data.Content + fmt.Sprintf("\n <font color=005EFF>@%s</font> \n", splicingString)
	}
	err = api.SendMessage(data)
	if err != nil {
		zap.L().Error("消息发送失败:", zap.Error(err))
		return
	}
	return
}
