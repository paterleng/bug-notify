package main

import (
	"bug-notify/handle"
	init_tool "bug-notify/init-tool"
	"github.com/gin-gonic/gin"
)

func main() {
	err := init_tool.Init()
	if err != nil {
		return
	}
	engine := gin.Default()
	//启动一个协程用于执行binlog

	//handle.Ttttt()

	go handle.NotifyHandle()

	engine.Run(init_tool.Conf.ProjectConfig.Address + ":" + init_tool.Conf.ProjectConfig.Port)

}
