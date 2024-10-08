package init_tool

import (
	"fmt"
	"go.uber.org/zap"
)

func Init() (err error) {
	err = ViperInit()
	if err != nil {
		fmt.Errorf("init viper failed,err: %v", err)
		return
	}
	err = LoggerInit()
	if err != nil {
		fmt.Errorf("init logger failed,err: %v", err)
		return
	}
	err = MysqlInit()
	if err != nil {
		lg.Error("init mysql failed,err: ", zap.Error(err))
		return
	}

	err = CreateTable()
	if err != nil {
		lg.Error("init table failed,err: ", zap.Error(err))
		return
	}
	return nil
}
