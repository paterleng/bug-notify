package main

import (
	"bug-notify/handle"
	init_tool "bug-notify/init-tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	err := init_tool.Init()
	if err != nil {
		return
	}
	engine := gin.Default()
	//启动一个协程用于执行binlog
	//go safelyRun(handle.Ttttt)
	go safelyRun(handle.NotifyHandle)
	//go safelyRun(handle.TimeingTasks)
	engine.Run(init_tool.Conf.ProjectConfig.Address + ":" + init_tool.Conf.ProjectConfig.Port)
}

func safelyRun(f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			zap.L().Error("数据中心异常，请联系管理员处理")
		}
	}()
	f()
}
